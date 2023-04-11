package repoimpl

import (
	"fmt"
	"os"
	"strings"
)

func createRepository(upperNameOfModel, lowerNameOfModel string, fields []*field) string {
	fmt.Println(upperNameOfModel, lowerNameOfModel)
	str := strings.Builder{}

	str.WriteString(headerOfRepository(upperNameOfModel, lowerNameOfModel, fields))
	str.WriteString(createMethod(upperNameOfModel, lowerNameOfModel, fields))
	str.WriteString(getByIDMethod(upperNameOfModel, lowerNameOfModel, fields))
	str.WriteString(getAllMethod(upperNameOfModel, lowerNameOfModel, fields))
	str.WriteString(updateMethod(upperNameOfModel, lowerNameOfModel, fields))
	str.WriteString(deleteMethod(upperNameOfModel, lowerNameOfModel, fields))

	return str.String()
}
func repositoryFiles(repos, nameOfModel string) error {

	file, err := os.Create(repositoryPath + "/" + postgresPath + "/" + nameOfModel + ".go")
	if err != nil {
		return clear(err)
	}
	_, err = file.Write([]byte(repos))
	if err != nil {
		return clear(err)
	}
	return nil
}

func headerOfRepository(upperNameOfModel, lowerNameOfModel string, fields []*field) string {

	str := strings.Builder{}
	str.WriteString("package postgres\n\n")
	str.WriteString("import (\n")
	str.WriteString("\t\"context\"\n")
	str.WriteString("\t\"database/sql\"\n")

	str.WriteString(")\n\n")
	str.WriteString("const (\n")
	str.WriteString(fmt.Sprintf("\t%sTable = \"%ss\"\n", lowerNameOfModel, lowerNameOfModel))
	str.WriteString(fmt.Sprintf("\tfieldsOf%s = \"%s\"\n", upperNameOfModel, strings.Join(getFields(fields), ",")))
	str.WriteString(")\n\n")
	str.WriteString(fmt.Sprintf("type %sRepository struct {\n", lowerNameOfModel))
	str.WriteString("\tdb *sql.DB\n")
	str.WriteString("}\n\n")
	str.WriteString(fmt.Sprintf("func New%sRepository(db *sql.DB) *%sRepository {\n", upperNameOfModel, lowerNameOfModel))
	str.WriteString(fmt.Sprintf("\treturn &%sRepository{db: db}\n", lowerNameOfModel))
	str.WriteString("}\n\n")
	return str.String()
}
func createMethod(upperNameOfModel, lowerNameOfModel string, fields []*field) string {

	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("// Create%s implements storage.%s\n", upperNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("func (r *%sRepository) Create%s(ctx context.Context, %s *models.%s) (*models.%s, error) {\n\n", lowerNameOfModel, upperNameOfModel, lowerNameOfModel, upperNameOfModel, upperNameOfModel))
	str.WriteString("\t// response result\n")
	str.WriteString(fmt.Sprintf("\tvar result models.%s\n\n", upperNameOfModel))
	str.WriteString("\t// query\n")
	str.WriteString(fmt.Sprintf("\tquery := `INSERT INTO %ss (` + fieldsOf%s + `", lowerNameOfModel, upperNameOfModel))
	str.WriteString(") VALUES (")
	for i, _ := range fields {
		if i != len(fields)-1 {
			str.WriteString(fmt.Sprintf("$%d,", i+1))
		} else {
			str.WriteString(fmt.Sprintf("$%d", i+1))
		}
	}
	str.WriteString(fmt.Sprintf(") RETURNING ` + fieldsOf%s", upperNameOfModel))

	str.WriteString("\n\n\t// execute query and scan result\n")
	str.WriteString("\terr := r.db.QueryRowContext(ctx, query,\n")
	for i, v := range fields {
		if i != len(fields)-1 {
			str.WriteString(fmt.Sprintf("\t\t%s.%s,\n", lowerNameOfModel, v.Name))
		} else {
			str.WriteString(fmt.Sprintf("\t\t%s.%s,\n", lowerNameOfModel, v.Name))
		}
	}
	str.WriteString("\t).Scan(\n")
	for i, v := range fields {
		if i != len(fields)-1 {
			str.WriteString(fmt.Sprintf("\t\t&result.%s,\n", v.Name))
		} else {
			str.WriteString(fmt.Sprintf("\t\t&result.%s,\n", v.Name))
		}
	}
	str.WriteString("\t)\n\n")
	str.WriteString("\t// checking error\n")
	str.WriteString("\tif err != nil {\n")
	str.WriteString("\t\treturn nil, err\n")
	str.WriteString("\t}\n\n")
	str.WriteString("\t// return result\n")
	str.WriteString("\treturn &result, nil\n")
	str.WriteString("}\n\n")
	return str.String()
}

func getByIDMethod(upperNameOfModel, lowerNameOfModel string, fields []*field) string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("// GetBy%s%s implements storage.%s\n", fields[0].Name, upperNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("func (r *%sRepository) Get%sByID(ctx context.Context, id int64) (*models.%s, error) {\n\n", lowerNameOfModel, upperNameOfModel, upperNameOfModel))
	str.WriteString("\t// response result\n")
	str.WriteString(fmt.Sprintf("\tvar result models.%s\n\n", upperNameOfModel))
	str.WriteString("\t// query\n")
	str.WriteString(fmt.Sprintf("\tquery := `SELECT ` + fieldsOf%s + ` FROM %ss WHERE %s = $1`", upperNameOfModel, lowerNameOfModel, fields[0].DBTag))
	str.WriteString("\n\n\t// execute query and scan result\n")
	str.WriteString("\terr := r.db.QueryRowContext(ctx, query, id).Scan(\n")
	for i, v := range fields {
		if i != len(fields)-1 {
			str.WriteString(fmt.Sprintf("\t\t&result.%s,\n", v.Name))
		} else {
			str.WriteString(fmt.Sprintf("\t\t&result.%s,\n", v.Name))
		}
	}
	str.WriteString("\t)\n\n")
	str.WriteString("\t// checking error\n")
	str.WriteString("\tif err != nil {\n")
	str.WriteString("\t\treturn nil, err\n")
	str.WriteString("\t}\n\n")
	str.WriteString("\t// return result\n")
	str.WriteString("\treturn &result, nil\n")
	str.WriteString("}\n\n")
	return str.String()
}

func getAllMethod(upperNameOfModel, lowerNameOfModel string, fields []*field) string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("// GetAll%s implements storage.%s\n", upperNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("func (r *%sRepository) GetAll%s(ctx context.Context) ([]*models.%s, error) {\n\n", lowerNameOfModel, upperNameOfModel, upperNameOfModel))
	str.WriteString("\t// response result\n")
	str.WriteString(fmt.Sprintf("\tvar result []*models.%s\n\n", upperNameOfModel))
	str.WriteString("\t// query\n")
	str.WriteString(fmt.Sprintf("\tquery := `SELECT ` + fieldsOf%s + `FROM %ss`", upperNameOfModel, lowerNameOfModel))
	str.WriteString("\n\n\t// execute query and scan result\n")
	str.WriteString("\trows, err := r.db.QueryContext(ctx, query)\n")
	str.WriteString("\tif err != nil {\n")
	str.WriteString("\t\treturn nil, err\n")
	str.WriteString("\t}\n\n")
	str.WriteString("\t// scan result\n")
	str.WriteString("\tfor rows.Next() {\n")
	str.WriteString(fmt.Sprintf("\t\tvar %s models.%s \n", lowerNameOfModel, upperNameOfModel))
	str.WriteString("\t\terr := rows.Scan(\n")
	for i, v := range fields {
		if i != len(fields)-1 {
			str.WriteString(fmt.Sprintf("\t\t\t&%s.%s,\n", lowerNameOfModel, v.Name))
		} else {
			str.WriteString(fmt.Sprintf("\t\t\t&%s.%s,\n", lowerNameOfModel, v.Name))
		}
	}

	str.WriteString("\t\t)\n")
	str.WriteString("\t\tif err != nil {\n")
	str.WriteString("\t\t\treturn nil, err\n")
	str.WriteString("\t\t}\n")
	str.WriteString(fmt.Sprintf("\t\tresult = append(result, &%s)\n", lowerNameOfModel))
	str.WriteString("\t}\n\n")
	str.WriteString("\t// checking error\n")
	str.WriteString("\tif err := rows.Err(); err != nil {\n")
	str.WriteString("\t\treturn nil, err\n")
	str.WriteString("\t}\n\n")
	str.WriteString("\t// return result\n")
	str.WriteString("\treturn result, nil\n")
	str.WriteString("}\n\n")
	return str.String()
}

func updateMethod(upperNameOfModel, lowerNameOfModel string, fields []*field) string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("// Update%s implements storage.%s\n", upperNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("func (r *%sRepository) Update%s(ctx context.Context, %s *models.%s) (*models.%s, error){\n\n", lowerNameOfModel, upperNameOfModel, lowerNameOfModel, upperNameOfModel, upperNameOfModel))
	str.WriteString("\t// response result\n")
	str.WriteString(fmt.Sprintf("\tvar result models.%s\n\n", upperNameOfModel))
	str.WriteString("\t// query\n")
	str.WriteString(fmt.Sprintf("\tquery := `UPDATE %ss SET ", lowerNameOfModel))
	for i, v := range fields {
		if i != len(fields)-1 {
			str.WriteString(fmt.Sprintf("%s = $%d, ", v.DBTag, i+2))
		} else {
			str.WriteString(fmt.Sprintf("%s = $%d ", v.DBTag, i+2))
		}
	}
	str.WriteString(fmt.Sprintf("WHERE %s = $1`", fields[0].DBTag))
	str.WriteString("\n\n\t// execute query and scan result\n")
	str.WriteString(fmt.Sprintf("\terr := r.db.QueryRowContext(ctx, query, %s.%s).Scan(\n", lowerNameOfModel, fields[0].Name))
	for i, v := range fields {
		if i != len(fields)-1 {
			str.WriteString(fmt.Sprintf("\t\t&result.%s,\n", v.Name))
		} else {
			str.WriteString(fmt.Sprintf("\t\t&result.%s,\n", v.Name))
		}
	}
	str.WriteString("\t)\n\n")
	str.WriteString("\t// checking error\n")
	str.WriteString("\tif err != nil {\n")
	str.WriteString("\t\treturn err\n")
	str.WriteString("\t}\n\n")
	str.WriteString("\t// return result\n")
	str.WriteString("\treturn nil\n")
	str.WriteString("}\n\n")
	return str.String()
}

func deleteMethod(upperNameOfModel, lowerNameOfModel string, fields []*field) string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("// Delete%s implements storage.%s\n", upperNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("func (r *%sRepository) Delete%s(ctx context.Context, %s %s) (*models.%s, error) {\n\n", lowerNameOfModel, upperNameOfModel, fields[0].DBTag, fields[0].Type, upperNameOfModel))
	str.WriteString("\t// response result\n")
	str.WriteString(fmt.Sprintf("\tvar result models.%s\n\n", upperNameOfModel))
	str.WriteString("\t// query\n")
	str.WriteString(fmt.Sprintf("\tquery := `DELETE FROM %ss WHERE %s = $1`", lowerNameOfModel, fields[0].DBTag))
	str.WriteString("\n\n\t// execute query and scan result\n")
	str.WriteString(fmt.Sprintf("\terr := r.db.QueryRowContext(ctx, query, %s.%s).Scan(\n", lowerNameOfModel, fields[0].Name))
	for i, v := range fields {
		if i != len(fields)-1 {
			str.WriteString(fmt.Sprintf("\t\t&result.%s,\n", v.Name))
		} else {
			str.WriteString(fmt.Sprintf("\t\t&result.%s,\n", v.Name))
		}
	}
	str.WriteString("\t)\n\n")
	str.WriteString("\t// checking error\n")
	str.WriteString("\tif err != nil {\n")
	str.WriteString("\t\treturn err\n")
	str.WriteString("\t}\n\n")
	str.WriteString("\t// return result\n")
	str.WriteString("\treturn nil\n")
	str.WriteString("}\n\n")
	return str.String()
}

func storageInterface(upperNameOfModel, lowerNameOfModel string, fields []*field) (string, string) {

	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("// %sI interface\n", upperNameOfModel))
	str.WriteString(fmt.Sprintf("type %sI interface {\n", upperNameOfModel))
	str.WriteString(fmt.Sprintf("\tCreate%s(ctx context.Context, %s *models.%s) (*models.%s, error)\n", upperNameOfModel, lowerNameOfModel, upperNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("\tGet%sBy%s(ctx context.Context, %s %s) (*models.%s, error)\n", upperNameOfModel, fields[0].Name, fields[0].Name, fields[0].Type, upperNameOfModel))
	str.WriteString(fmt.Sprintf("\tGetAll%s(ctx context.Context) ([]*models.%s, error)\n", upperNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("\tUpdate%s(ctx context.Context, %s *models.%s) (*models.%s, error)\n", upperNameOfModel, lowerNameOfModel, upperNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("\tDelete%s(ctx context.Context, %s %s) (*models.%s, error)\n", upperNameOfModel, fields[0].Name, fields[0].Type, upperNameOfModel))
	str.WriteString("}\n\n")
	return fmt.Sprintf("\t%s() %sI\n", upperNameOfModel, upperNameOfModel), str.String()
}
func storageHeader() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("package storage\n\n"))
	str.WriteString("import (\n")
	str.WriteString("\t\"context\"\n")
	str.WriteString(")\n\n")
	str.WriteString("type Storage interface {\n")
	str.WriteString("\tCloseDb() error\n")
	return str.String()
}

func storageFile(str string) error {
	f, err := os.Create(repositoryPath + "/storage.go")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(str)
	if err != nil {
		return err
	}
	return nil
}

func storageFooter() string {
	str := strings.Builder{}
	str.WriteString("}\n\n")
	return str.String()
}
