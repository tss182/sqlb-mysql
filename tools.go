package sqlb

import (
	"reflect"
	"strconv"
)

func valueInterface(value interface{}) [2]string {
	var reflectValue = reflect.ValueOf(value)

	var result [2]string
	switch reflectValue.Kind() {
	case reflect.String:
		result[0] = value.(string)
		result[1] = "string"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := int(reflectValue.Uint())
		result[0] = strconv.Itoa(i)
		result[1] = "int"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := int(reflectValue.Int())
		result[0] = strconv.Itoa(i)
		result[1] = "int"
	case reflect.Bool:
		logic := reflectValue.Bool()
		result[0] = strconv.FormatBool(logic)
		result[1] = "bool"
	}

	return result
}
