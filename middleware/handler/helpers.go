package handler

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

func getFieldName(namespace string, reqs map[string]interface{}) string {
	nsParts := strings.Split(namespace, ".")
	if len(nsParts) < 2 {
		return namespace
	}
	obj, ok := reqs[nsParts[0]]
	if !ok {
		return namespace[1:]
	}
	return buildFieldPath("", nsParts[1:], obj)
}

func buildFieldPath(resultField string, nsParts []string, obj interface{}) string {
	if len(nsParts) == 0 {
		return resultField
	}
	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		obj = reflect.New(reflect.TypeOf(obj).Elem()).Elem().Interface()
		return buildFieldPath(resultField, nsParts, obj)
	}
	nsField := nsParts[0]
	fieldPostfix := ""
	if strings.Contains(nsField, "[") {
		i := strings.Index(nsField, "[")
		fieldPostfix = nsField[i:]
		nsField = nsField[:i]
	}
	var structField reflect.StructField
	var nextObj interface{}
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Struct:
		sf, ok := reflect.TypeOf(obj).FieldByName(nsField)
		if !ok {
			return resultField
		}
		structField = sf
		nextObj = reflect.ValueOf(obj).FieldByName(structField.Name).Interface()
	default:
		return resultField
	}
	if len(resultField) > 0 {
		resultField = fmt.Sprintf("%s.", resultField)
	}
	nextPart := structField.Name
	if tag, ok := getStructJSONTag(structField); ok {
		nextPart = tag
	}
	resultField = resultField + nextPart + fieldPostfix
	return buildFieldPath(resultField, nsParts[1:], nextObj)
}

func getStructJSONTag(f reflect.StructField) (string, bool) {
	str := f.Tag.Get("json")
	if len(str) == 0 {
		return "", false
	}
	return str, true
}

// --------------------

func ucFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func lcFirst(str string) string {
	return strings.ToLower(str)
}

func split(src string) string {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	var entries []string
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}

	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}

	for index, word := range entries {
		if index == 0 {
			entries[index] = ucFirst(word)
		} else {
			entries[index] = lcFirst(word)
		}
	}
	justString := strings.Join(entries, " ")
	return justString
}
