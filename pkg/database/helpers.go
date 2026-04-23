package database

import (
	"reflect"
)

// ColumnNames extracts column names from model if db tag is set.
func ColumnNames(v interface{}) []string {
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	cols := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}

		cols = append(cols, tag)
	}

	return cols
}

// Values extracts values from a model interface.
func Values(v interface{}) []interface{} {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	values := make([]interface{}, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		values = append(values, val.Field(i).Interface())
	}

	return values
}

// ColumnValues extracts column names and values and creates a key-value pair map used for set.
func ColumnValues(v interface{}) map[string]any {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}

	result := make(map[string]any, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")

		if tag == "" || tag == "-" {
			continue
		}

		result[tag] = val.Field(i).Interface()
	}

	return result
}
