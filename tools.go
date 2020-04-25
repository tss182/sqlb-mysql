package sqlb

import (
	"bytes"
	"html"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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

func removeSpecialChar(char interface{}) string {
	val := valueInterface(char)[0]
	val = string(bytes.Trim([]byte(val), "\xef\xbb\xbf"))
	reg, err := regexp.Compile("[^ -~]+")
	if err != nil {
		return ""
	}
	str := reg.ReplaceAllString(val, "")
	str = addSlash(str)
	str = html.EscapeString(str)
	return str
}

func addSlash(char string) string {
	var str = strings.Replace(char, "'", "\\'", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	return str
}
