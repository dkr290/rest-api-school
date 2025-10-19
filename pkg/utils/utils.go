// Package utils - utility functions for helpers or some common perpose
package utils

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

func GenereateInsertQuery(model any) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, "omitempty")
		dbTag = strings.TrimSuffix(dbTag, ",")
		// skip id field if it is auto incrment
		if dbTag != "" && dbTag != "id" {

			columns += dbTag + ","
			placeholders += "?,"
		}
	}
	columns = removeLastComma(columns)
	placeholders = removeLastComma(placeholders)

	return fmt.Sprintf("INSERT INTO teachers (%s) VALUES (%s)", columns, placeholders)
}

func GetStructValues(model any) []any {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}
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

func removeLastComma(s string) string {
	if idx := strings.LastIndex(s, ","); idx != -1 {
		s = s[:idx]
	}
	return s
}

func EmailCheck(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	em := strings.TrimSpace(email)
	if !emailRegex.MatchString(em) {
		return fmt.Errorf("invalid email: %s", em)
	}
	return nil
}
