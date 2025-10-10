// Package utils - utility functions for helpers or some common perpose
package utils

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func GenereateInsertQuery(model any) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println("Db tag:", dbTag)
		dbTag = strings.TrimSuffix(dbTag, "omitempty")
		// skip id field if it is auto incrment
		if dbTag != "" && dbTag != "id" {
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			placeholders += "?"
		}
	}
	return fmt.Sprintf("INSERT INTO teachers (%s) VALUES (%s)", columns, placeholders)
}

func GetStructValues(model any) []any {
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()

	var values []any

	for i := 0; i < modelType.NumField(); i++ {

		dbTag := modelType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" {
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	log.Println("Values", values)
	return values
}
