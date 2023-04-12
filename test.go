package repoimpl

import (
	"os"
)

func createTest(nameOfModel string, fields []*field) string {
	return ""
}
func testFiles(tests, nameOfModel string) error {

	file, err := os.Create(testPath + "/" + postgresPath + "/" + nameOfModel + ".go")
	if err != nil {
		return clear("test", err)
	}
	_, err = file.Write([]byte(tests))
	if err != nil {
		return clear("test", err)
	}
	return nil
}
