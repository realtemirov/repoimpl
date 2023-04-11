package repoimpl

import (
	"fmt"
	"os"
	"strings"
)

func createDbQuery(nameOfModel string, fields []*field) (string, string) {

	crt := strings.Builder{}
	crt.Write([]byte(fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%ss\" (\n", nameOfModel)))
	for i, v := range fields {
		if i != len(fields)-1 {
			crt.Write([]byte(fmt.Sprintf("\t\"%s\" %s,\n", v.DBTag, goTypeToDBType(v.Type))))
		} else {
			crt.Write([]byte(fmt.Sprintf("\t\"%s\" %s\n", v.DBTag, goTypeToDBType(v.Type))))
		}
	}
	crt.Write([]byte(");"))

	drp := fmt.Sprintf("DROP TABLE IF EXISTS \"%ss\";", nameOfModel)
	return crt.String(), drp
}

func migrationFiles(nameOfModel, create, drop string) error {

	file, err := os.Create(fmt.Sprintf("%s/%s/01_create_%s.up.sql", migrationPath, postgresPath, nameOfModel))
	if err != nil {
		return clear(err)
	}
	_, err = file.Write([]byte(create))
	if err != nil {
		return clear(err)
	}

	file, err = os.Create(fmt.Sprintf("%s/%s/01_create_%s.down.sql", migrationPath, postgresPath, nameOfModel))
	if err != nil {
		return clear(err)
	}
	_, err = file.Write([]byte(drop))
	if err != nil {
		return clear(err)
	}
	return nil
}
