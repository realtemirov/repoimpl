package repoimpl

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

const (
	migrationPath  = "migration"
	repositoryPath = "storage"
	testPath       = "test"
	postgresPath   = "postgres"
	modelPath      = "model"
)

// Implement repository, migration and test
// Parameters:
//   - model... : implementation of models with table name and fields
//
// Example:
//   - NewProject(models.User{})
//
// Returns:
//   - error: error if something went wrong
func NewProject(datas ...any) error {

	log.Println("Starting implement...")

	err := NewMigration(datas...)
	if err != nil {
		return err
	}

	err = NewRepository(datas...)
	if err != nil {
		return err
	}

	err = NewTest(datas...)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(datas ...any) error {

	var (
		postgresUp   strings.Builder
		postgresDown strings.Builder
		storage      strings.Builder
		interfaces   strings.Builder
	)
	err := os.Mkdir(repositoryPath, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			fmt.Println("repository folder already exist, If you want to create new repository click something else `y` or `Y`")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) != "y" {
				return nil
			}
			err := os.RemoveAll(repositoryPath)
			os.Mkdir(repositoryPath, 0755)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
	}

	err = os.Mkdir(repositoryPath+"/"+postgresPath, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			fmt.Println("repository/postgres folder already exist, If you want to create new repository/postgres click something else `y` or `Y`")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) != "y" {
				return nil
			}
			err := os.RemoveAll(repositoryPath + "/" + postgresPath)
			if err != nil {
				return err
			}
		}
	}

	storage.WriteString(storageHeader())
	postgresUp.WriteString(postgresHeader())
	for _, model := range datas {

		lowerNameOfModel := fieldToDefault(reflect.TypeOf(model).Name())
		upperNameOfModel := reflect.TypeOf(model).Name()
		valueOfModel := reflect.ValueOf(model)
		typeOfModel := valueOfModel.Type()
		numberOfFields := valueOfModel.NumField()

		fieldsOfModel := fields(typeOfModel, numberOfFields)

		log.Println("Starting repository...")
		err = repositoryFiles(createRepository(upperNameOfModel, lowerNameOfModel, fieldsOfModel), lowerNameOfModel)
		if err != nil {
			return clear("repository", err)
		}
		interfaceName, interfaceMethods := storageInterface(upperNameOfModel, lowerNameOfModel, fieldsOfModel)
		storage.WriteString(interfaceName)
		interfaces.WriteString(interfaceMethods)

		postgresUpString, postgresDownString := postgresInterface(upperNameOfModel, lowerNameOfModel)

		postgresUp.WriteString(postgresUpString)
		postgresDown.WriteString(postgresDownString)
	}
	postgresUp.WriteString(postgresFooter())
	postgresDown.WriteString(postgresNew())
	storage.WriteString(storageFooter())
	storage.WriteString(interfaces.String())
	storageFile(storage.String())
	postgresFile(postgresUp.String(), postgresDown.String())
	log.Println("Successful repository implemented")
	return nil
}

func NewMigration(datas ...any) error {
	err := os.Mkdir(migrationPath, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			fmt.Println("migration folder already exist, If you want to create new migration, click something else `y` or `Y`")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) != "y" {
				return nil
			}
			err := os.RemoveAll(migrationPath)
			if err != nil {
				return err
			}
		}
	} else {
		err := os.Mkdir(migrationPath+"/"+postgresPath, 0755)
		if err != nil {
			if strings.Contains(err.Error(), "exists") {
				fmt.Println("migration/postgres folder already exist, If you want to create new migration/postgres, click something else `y` or `Y`")
				var answer string
				fmt.Scanln(&answer)
				if strings.ToLower(answer) != "y" {
					return nil
				}
				err := os.RemoveAll(migrationPath + "/" + postgresPath)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, model := range datas {

		lowerNameOfModel := fieldToDefault(reflect.TypeOf(model).Name())
		// upperNameOfModel := reflect.TypeOf(model).Name()
		valueOfModel := reflect.ValueOf(model)
		typeOfModel := valueOfModel.Type()
		numberOfFields := valueOfModel.NumField()

		fieldsOfModel := fields(typeOfModel, numberOfFields)

		log.Println("Starting migration...")
		init, drop := createDbQuery(lowerNameOfModel, fieldsOfModel)
		err := migrationFiles(lowerNameOfModel, init, drop)
		if err != nil {
			return err
		}
		log.Println("Successful migration implemented")

	}

	return nil
}
func NewTest(datas ...any) error {
	err := os.Mkdir(testPath, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			fmt.Println("test folder already exist, If you want to create new test, click something else `y` or `Y`")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) != "y" {
				return nil
			}
			err := os.RemoveAll(testPath)
			if err != nil {
				return err
			}
		}
	} else {
		err := os.Mkdir(testPath+"/"+postgresPath, 0755)
		if err != nil {
			if strings.Contains(err.Error(), "exists") {
				fmt.Println("test/postgres folder already exist, If you want to create new test/postgres, click something else `y` or `Y`")
				var answer string
				fmt.Scanln(&answer)
				if strings.ToLower(answer) != "y" {
					return nil
				}
				err := os.RemoveAll(testPath + "/" + postgresPath)
				if err != nil {
					return err
				}
			}
		}
	}
	for _, model := range datas {

		lowerNameOfModel := fieldToDefault(reflect.TypeOf(model).Name())
		//		upperNameOfModel := reflect.TypeOf(model).Name()
		valueOfModel := reflect.ValueOf(model)
		typeOfModel := valueOfModel.Type()
		numberOfFields := valueOfModel.NumField()

		fieldsOfModel := fields(typeOfModel, numberOfFields)

		log.Println("Starting test...")
		err = testFiles(createTest(lowerNameOfModel, fieldsOfModel), lowerNameOfModel)
		if err != nil {
			return clear("test", err)
		}

		log.Println("Successful test implemented")
	}

	return nil
}

func clear(data string, response error) error {

	switch data {
	case "repository":
		{
			if _, err := os.Stat(repositoryPath + "/" + postgresPath); os.IsExist(err) {
				err = os.RemoveAll(repositoryPath + "/" + postgresPath)
				if err != nil {
					log.Printf("Error: %v", err)
				}
			} else {
				log.Printf("Error: %v", err)
			}
		}
	case "migration":
		{
			if _, err := os.Stat(migrationPath + "/" + postgresPath); os.IsExist(err) {
				err = os.RemoveAll(migrationPath + "/" + postgresPath)
				if err != nil {
					log.Printf("Error: %v", err)
				}
			} else {
				log.Printf("Error: %v", err)
			}
		}
	case "test":
		{
			if _, err := os.Stat(testPath + "/" + postgresPath); os.IsExist(err) {
				err = os.RemoveAll(testPath + "/" + postgresPath)
				if err != nil {
					log.Printf("Error: %v", err)
				}
			} else {
				log.Printf("Error: %v", err)
			}
		}
	}

	return response
}

func NewModels(datas ...any) error {
	err := os.Mkdir(modelPath, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			fmt.Println("models folder already exist, If you want to create new migration, click something else `y` or `Y`")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) != "y" {
				return nil
			}
			err := os.RemoveAll(migrationPath)
			if err != nil {
				return err
			}
		}
	} else {
		err := os.Mkdir(migrationPath+"/"+postgresPath, 0755)
		if err != nil {
			if strings.Contains(err.Error(), "exists") {
				fmt.Println("migration/postgres folder already exist, If you want to create new migration/postgres, click something else `y` or `Y`")
				var answer string
				fmt.Scanln(&answer)
				if strings.ToLower(answer) != "y" {
					return nil
				}
				err := os.RemoveAll(migrationPath + "/" + postgresPath)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, model := range datas {

		lowerNameOfModel := fieldToDefault(reflect.TypeOf(model).Name())
		// upperNameOfModel := reflect.TypeOf(model).Name()
		valueOfModel := reflect.ValueOf(model)
		typeOfModel := valueOfModel.Type()
		numberOfFields := valueOfModel.NumField()

		fieldsOfModel := fields(typeOfModel, numberOfFields)

		log.Println("Starting migration...")
		init, drop := createDbQuery(lowerNameOfModel, fieldsOfModel)
		err := migrationFiles(lowerNameOfModel, init, drop)
		if err != nil {
			return err
		}
		log.Println("Successful migration implemented")

	}

	return nil
}
