// Package utils - utility functions for helpers or some common perpose
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

func GenereateInsertQuery(model any, name string) string {
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

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", name, columns, placeholders)
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

func PasswordHash(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}

func SighnToken(userID, username, role string, config config.Config) (string, error) {
	jwtSecret := config.JWTSecret
	jwtExpiresIn := config.JWTExpiresIn

	claims := jwt.MapClaims{
		"uid":  userID,
		"user": username,
		"role": role,
	}

	claims["exp"] = jwt.NewNumericDate(time.Now().Add(jwtExpiresIn))

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := tkn.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("internal error %v", err)
	}

	return signedToken, nil
}

func SendTestResetEmail(
	to, message string,
	host string,
	port string,
) error {
	from := "noreply@school-rest-api.example.com"
	subject := "Your password reset link"

	body := message

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, subject, body)

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", host, port),
		nil,
		from,
		[]string{to},
		[]byte(msg),
	)
}
