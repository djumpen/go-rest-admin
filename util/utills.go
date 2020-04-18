package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// TODO: Function can't work with time.Time values. Unmarshal panics due to empty value
// Copy values from one struct to another
// Panic if `to` does not have every field defined in `from`
func CopyFields(from, to interface{}, skipFields ...string) interface{} {
	val := reflect.ValueOf(from)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	typeOfT := val.Type()
	for i := 0; i < typeOfT.NumField(); i++ {
		fromFieldName := typeOfT.Field(i).Name
		if StringInSlice(fromFieldName, skipFields) {
			continue
		}
		if err := StructFieldExists(to, fromFieldName); err != nil {
			panic(err)
		}
	}
	formBytes, err := json.Marshal(from)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(formBytes, to)
	if err != nil {
		panic(err)
	}
	return to
}

// Reflect if an interface is either a struct or a pointer to a struct
// and has the defined member field, if error is nil, the given
// FieldName exists and is accessible with reflect.
func StructFieldExists(Iface interface{}, FieldName string) error {
	ValueIface := reflect.ValueOf(Iface)

	// Check if the passed interface is a pointer
	if ValueIface.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueIface = reflect.New(reflect.TypeOf(Iface))
	}

	// 'dereference' with Elem() and get the field by name
	Field := ValueIface.Elem().FieldByName(FieldName)
	if !Field.IsValid() {
		return fmt.Errorf("Interface `%s` does not have the field `%s`", ValueIface.Type(), FieldName)
	}
	return nil
}

// Print json representation
func Println(vars ...interface{}) {
	for _, v := range vars {
		b, _ := json.Marshal(v)
		fmt.Println(string(b))
	}
}

// Format timestamp to "2006-01-02" format
func FmtDateTimeToDate(value string) string {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// Format date to date and time format
func FmtDateToDateTime(value string) string {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func StringDateToTime(str string) (time.Time, error) {
	// layout := "2006-01-02T15:04:05.000Z"
	layout := time.RFC3339
	t, err := time.Parse(layout, str)
	if err != nil {
		return t, err
	}
	return t, nil
}
