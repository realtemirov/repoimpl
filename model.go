package repoimpl

import "reflect"

type field struct {
	Name    string
	Type    string
	DBTag   string
	JSONTag string
}

func getFields(fields []*field) []string {
	var fieldsOfModel []string
	for _, v := range fields {
		fieldsOfModel = append(fieldsOfModel, v.DBTag)
	}
	return fieldsOfModel
}

func fields(typeOfModel reflect.Type, numberOfFields int) []*field {
	models := make(map[string]*field)
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
		models[fieldType.Name] = &field
		fieldsOfModel = append(fieldsOfModel, models[fieldType.Name])
	}
	return fieldsOfModel
}
func fieldToDefault(fieldName string) string {
	var dbColumnName string
	for i, v := range fieldName {
		if isUpperLetter(v) {
			if i != 0 {
				if isUpperLetter(rune(fieldName[i-1])) {
					dbColumnName += string(upperLetterToLowerLetter(v))
				} else {
					dbColumnName += "_" + string(upperLetterToLowerLetter(v))
				}
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
func goTypeToDBType(goType string) string {
	var result string
	switch goType {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		result = "INTEGER"
	case "float32", "float64":
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
