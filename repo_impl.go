package repoImpl

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

// Implement migration SQL queries for all models
// Parameters:
//   - model... : implementation of models with table name and fields
//
// Example:
//   - NewDBTable(models.User{})
//
// Returns:
//   - error: error if something went wrong
func NewDBTable(datas ...any) error {

	err := os.Mkdir("./migrations", 0755)
	if err != nil {
		if err == os.ErrExist {
			err := os.RemoveAll("./migrations")
			if err != nil {
				return err
			}
		}
	}
	log.Println("Starting migration implement...")
	for _, model := range datas {

		nameOfModel := fieldToDefault(reflect.TypeOf(model).Name())

		valueOfModel := reflect.ValueOf(model)
		typeOfModel := valueOfModel.Type()
		numberOfFields := valueOfModel.NumField()

		var fieldsOfModel []*field

		for i := 0; i < numberOfFields; i++ {

			fieldType := typeOfModel.Field(i)
			dbTag, ok := fieldType.Tag.Lookup("db")
			if !ok {
				dbTag = fieldToDefault(fieldType.Name)
			}
			field := field{
				Name:  fieldType.Name,
				Type:  fieldType.Type.Name(),
				DBTag: dbTag,
			}
			fieldsOfModel = append(fieldsOfModel, &field)
		}
		dbCreateQuery := createDbQuery(nameOfModel, fieldsOfModel)
		file, err := os.Create("./migrations/migrations_" + nameOfModel + ".sql")
		if err != nil {
			os.RemoveAll("migrations")
			return err
		}
		_, err = file.Write([]byte(dbCreateQuery))
		if err != nil {
			os.RemoveAll("migrations")
			return err
		}
	}
	log.Println("Succesfull implemented...")
	return nil
}

func fieldToDefault(fieldName string) string {
	var dbColumnName string
	for i, v := range fieldName {
		if isUpperLetter(v) {
			if i != 0 {
				dbColumnName += "_" + string(upperLetterToLowerLetter(v))
			} else {
				dbColumnName += string(upperLetterToLowerLetter(v))
			}
		} else {
			dbColumnName += string(v)
		}
	}
	return dbColumnName
}

func isUpperLetter(letter rune) bool {
	if 65 <= letter && letter <= 90 {
		return true
	}
	return false
}

func upperLetterToLowerLetter(letter rune) rune {
	return letter + 32
}

func createDbQuery(nameOfModel string, fields []*field) string {

	str := strings.Builder{}
	str.Write([]byte(fmt.Sprintf("CREATE TABLE IF NOT EXITS \"%ss\" (\n", nameOfModel)))
	for i, v := range fields {
		if i != len(fields)-1 {
			str.Write([]byte(fmt.Sprintf("\t\"%s\" %s,\n", v.DBTag, goTypeToDBType(v.Type))))
		} else {
			str.Write([]byte(fmt.Sprintf("\t\"%s\" %s\n", v.DBTag, goTypeToDBType(v.Type))))
		}
	}
	str.Write([]byte(");"))
	return str.String()
}

func goTypeToDBType(goType string) string {
	var result string
	switch goType {
	case "int":
		result = "NUMERIC"
	case "string":
		result = "TEXT"
	case "bool":
		result = "BOOLEN"
	case "Time":
		result = "TIMESTAMP"
	default:
		result = "ANY"
	}
	return result
}

type field struct {
	Name    string
	Type    string
	DBTag   string
	JSONTag string
}
