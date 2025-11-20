package models

import "database/sql"

type Exec struct {
	ID                   int            `json:"id,omitempty"              db:"id,omitempty"`
	FirstName            string         `json:"first_name,omitempty"      db:"first_name,omitempty"`
	LastName             string         `json:"last_name,omitempty"       db:"last_name,omitempty"`
	Email                string         `json:"email,omitempty"           db:"email,omitempty"`
	Username             string         `json:"username,omitempty"        db:"username,omitempty"`
	Password             string         `json:"password,omitempty"        db:"password,omitempty"`
	PasswordChangedAt    sql.NullString `json:"password_changed_at"       db:"password_changed_at"`
	UserCreatedAt        sql.NullString `json:"user_created_at"           db:"user_created_at"`
	PasswordResetToken   sql.NullString `json:"password_reset_token"      db:"password_reset_token"`
	PasswordTokenExpires sql.NullString `json:"password_token_expires"    db:"password_token_expires"`
	InactiveStatus       bool           `json:"inactive_status,omitempty" db:"inactive_status,omitempty"`
	Role                 string         `json:"role,omitempty"            db:"role,omitempty"`
}
type ExecLoginInput struct {
	Username string `json:"username" required:"true" minLength:"2" maxLength:"255" doc:"username" examle:"username"`
	Password string `json:"password" required:"true" minLength:"2" maxLength:"255" doc:"password"                   example:"password"`
}

type ExecInput struct {
	FirstName string `json:"first_name" required:"true" minLength:"2" maxLength:"255" example:"Tom"                 doc:"First name of the exec"`
	LastName  string `json:"last_name"  required:"true" minLength:"2" maxLength:"255" example:"Last"                doc:"Last name of the exec"`
	Email     string `json:"email"      required:"true"               maxLength:"255" example:"teacher@example.com" doc:"Email"`
	Username  string `json:"username"   required:"true" minLength:"2" maxLength:"255"                               doc:"username"               examle:"username"`
	Password  string `json:"password"   required:"true" minLength:"2" maxLength:"255" example:"password"            doc:"password"`
	Role      string `json:"role"       required:"true"                               example:"admin"               doc:"role to use like admin"`
}

type ExecPatchBody struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name,omitempty" example:"Alice"           doc:"First name of the teacher"`
	LastName  string `json:"last_name,omitempty"  example:"Brown"           doc:"Last name of the techer"`
	Email     string `json:"email,omitempty"      example:"ac@example.com " doc:"Email"`
	Username  string `json:"username"                                       doc:"username"                  required:"true" minLength:"2" maxLength:"255" examle:"username"`
}

type ExecsQueryInput struct {
	FirstName string   `query:"first_name"`
	LastName  string   `query:"last_name"`
	Email     string   `query:"email"`
	Username  string   `query:"username"`
	Role      string   `query:"role"`
	SortBy    []string `query:"sort_by"    example:"first_name:asc" doc:"Order by asc or desc of the records"`
}
