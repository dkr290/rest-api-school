package models

import "database/sql"

type Exec struct {
	ID                int
	FirstName         string
	LastName          string
	Email             string
	Username          string
	Password          string
	PasswordChangedAt sql.NullString
	UserCreatedAt     sql.NullString
	PasswordResetCode string
	InactiveStatus    bool
	Role              string
}
