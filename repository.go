package repoimpl

import (
	"fmt"
	"os"
	"strings"
)

func createRepository(upperNameOfModel, lowerNameOfModel string, fields []*field) string {
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
		return clear("repository", err)
	}
	_, err = file.Write([]byte(repos))
	if err != nil {
		return clear("repository", err)
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
	for i := range fields {
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
	str.WriteString(fmt.Sprintf("func (r *%sRepository) Get%sByID(ctx context.Context,%s %s) (*models.%s, error) {\n\n", lowerNameOfModel, upperNameOfModel, fields[0].DBTag, fields[0].Type, upperNameOfModel))
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
	str.WriteString("\t\treturn nil, err\n")
	str.WriteString("\t}\n\n")
	str.WriteString("\t// return result\n")
	str.WriteString("\treturn &result, nil\n")
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
	str.WriteString(fmt.Sprintf("\terr := r.db.QueryRowContext(ctx, query, %s).Scan(\n", fields[0].DBTag))
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
	str.WriteString("type StorageI interface {\n")
	str.WriteString("\tCloseDB() error\n")
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

func postgresFile(up, down string) error {
	f, err := os.Create(repositoryPath + "/postgres/postgres.go")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(up)
	if err != nil {
		return err
	}
	_, err = f.WriteString(down)
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

func postgresHeader() string {

	str := strings.Builder{}

	str.WriteString(fmt.Sprintf("package postgres"))
	str.WriteString("\n\n")
	str.WriteString("import (\n")
	str.WriteString("\t\"context\"\n")
	str.WriteString("\t\"database/sql\"\n")
	str.WriteString("\t\"fmt\"\n")
	str.WriteString("\n")
	str.WriteString("\t_ \"github.com/lib/pq\"\n")
	str.WriteString("\t\"github.com/rs/zerolog\"\n")
	str.WriteString(")\n\n")
	str.WriteString(fmt.Sprintf("type Storage struct {\n"))
	str.WriteString(fmt.Sprintf("\tdb  *sql.DB        // database connection\n"))
	str.WriteString(fmt.Sprintf("\tlog zerolog.Logger // logger\n"))

	return str.String()
}

func postgresNew() string {
	str := strings.Builder{}
	str.WriteString("// CloseDB implements storage.StorageI\n")
	str.WriteString("func (s *Storage) CloseDB() error {\n")
	str.WriteString("\t// Close the database connection.\n")
	str.WriteString("\treturn s.db.Close()\n")
	str.WriteString("}\n\n")
	str.WriteString("// NewPostgres creates a new postgres storage\n")
	str.WriteString("func NewPostgres(ctx context.Context, cfg *config.Config, log zerolog.Logger) (storage.StorageI, error) {\n\n")
	str.WriteString("\t// Create a connection string to connect to the database.\n")
	str.WriteString("\tpostgresConnString := fmt.Sprintf(\n")
	str.WriteString("\t\t\"host=%%s port=%%s user=%%s dbname=%%s password=%%s sslmode=%%s\",\n")
	str.WriteString("\t\tcfg.PostgresHost,\n")
	str.WriteString("\t\tcfg.PostgresPort,\n")
	str.WriteString("\t\tcfg.PostgresUser,\n")
	str.WriteString("\t\tcfg.PostgresDB,\n")
	str.WriteString("\t\tcfg.PostgresPassword,\n")
	str.WriteString("\t\tcfg.PostgresSSLMode,\n")
	str.WriteString("\t)\n")
	str.WriteString("\tlog.Info().Msg(postgresConnString)\n")
	str.WriteString("\t// log")
	str.WriteString("\tlog.Info().Msg(\"Starting connection to the database...\")\n\n")
	str.WriteString("\t// Connect to the database.\n")
	str.WriteString("\tdb, err := sql.Open(\"postgres\", postgresConnString)\n\n")
	str.WriteString("\t// checking error\n")
	str.WriteString("\tif err != nil {\n")
	str.WriteString("\t\tlog.Info().AnErr(\"Method: NewPostgres Comment: Connect to the database Error: %%v\", err)\n")
	str.WriteString("\t\treturn nil, err\n")
	str.WriteString("\t}\n\n")
	str.WriteString("\t// Ping test\n")
	str.WriteString("\terr = db.Ping()\n\n")
	str.WriteString("\t// checking error\n")
	str.WriteString("\tif err != nil {\n")
	str.WriteString("\t\tlog.Info().AnErr(\"Method: NewPostgres Comment: Ping test Error: %%v\", err)\n")
	str.WriteString("\t\treturn nil, err\n")
	str.WriteString("\t}\n\n")
	str.WriteString("\t// log\n")
	str.WriteString("\tlog.Info().Msg(\"Connection to the database is successful.\")\n\n")
	str.WriteString("\t// Create a new instance of the Postgres storage and return it.\n")
	str.WriteString("\treturn &Storage{\n")
	str.WriteString("\t\tdb:  db,\n")
	str.WriteString("\t\tlog: log,\n")
	str.WriteString("\t}, nil\n")
	str.WriteString("}\n\n")

	return str.String()
}

func postgresInterface(upperNameOfModel, lowerNameOfModel string) (string, string) {

	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("// %s implements storage.StorageI\n", upperNameOfModel))
	str.WriteString(fmt.Sprintf("func (s *Storage) %s() storage.%sI {\n\n", upperNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("\t// if %sRepository is not nil, return it\n", lowerNameOfModel))
	str.WriteString(fmt.Sprintf("\tif s.%sRepository != nil {\n\n", lowerNameOfModel))
	str.WriteString(fmt.Sprintf("\t\t// return %sRepository\n", lowerNameOfModel))
	str.WriteString(fmt.Sprintf("\t\treturn s.%sRepository\n", lowerNameOfModel))
	str.WriteString(fmt.Sprintf("\t}\n\n"))
	str.WriteString(fmt.Sprintf("\t// if %sRepository is nil, create a new one\n", lowerNameOfModel))
	str.WriteString(fmt.Sprintf("\ts.%sRepository = new%sRepo(s.db, s.log)\n\n", lowerNameOfModel, upperNameOfModel))
	str.WriteString(fmt.Sprintf("\t// return %sRepository\n", lowerNameOfModel))
	str.WriteString(fmt.Sprintf("\treturn s.%sRepository\n", lowerNameOfModel))
	str.WriteString(fmt.Sprintf("}\n\n"))
	return fmt.Sprintf("\t%sRepository storage.%sI // %sRepository storage.%sI\n", lowerNameOfModel, upperNameOfModel, lowerNameOfModel, upperNameOfModel), str.String()

}

func postgresFooter() string {
	str := strings.Builder{}
	str.WriteString("}\n\n")
	return str.String()
}
