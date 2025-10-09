// Package models
package models

type Teacher struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" maxLength:"255" example:"Tom"     doc:"First name of the teacher" db:"id,omitempty"`
	LastName  string `json:"last_name"  maxLength:"255" example:"Last"    doc:"Last name of the techer"   db:"first_name,omitempty"`
	Class     string `json:"class"                      example:"10B"     doc:"The class of the teacher"  db:"last_name,omitempty"`
	Subject   string `json:"subject"    maxLength:"255" example:"History" doc:"Subject to teach"          db:"subject,omitempty"`
	Email     string `json:"email"      maxLength:"50"  example:"Email"   doc:"Email"                     db:"email,omitempty"`
}

type TeacherInput struct {
	FirstName string `json:"first_name" maxLength:"255" example:"Tom"     doc:"First name of the teacher"`
	LastName  string `json:"last_name"  maxLength:"255" example:"Last"    doc:"Last name of the techer"`
	Class     string `json:"class"                      example:"10B"     doc:"The class of the teacher"`
	Subject   string `json:"subject"    maxLength:"255" example:"History" doc:"Subject to teach"`
	Email     string `json:"email"      maxLength:"50"  example:"Email"   doc:"Email"`
}

type TeachersQueryInput struct {
	FirstName string   `query:"first_name"`
	LastName  string   `query:"last_name"`
	Class     string   `query:"class"`
	Subject   string   `query:"subject"`
	Email     string   `query:"email"`
	SortBy    []string `query:"sort_by"    example:"first_name:asc" doc:"Order by asc or desc of the records"`
}
type TeacherUpdateBody struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" example:"Alice"          doc:"First name of the teacher"`
	LastName  string `json:"last_name"  example:"Brown"          doc:"Last name of the techer"`
	Class     string `json:"class"      example:"9C"             doc:"The class of the teacher"`
	Subject   string `json:"subject"    example:"History"        doc:"Subject to teach"`
	Email     string `json:"email"      example:"ac@example.net" doc:"Email"`
}

type TeacherPatchBody struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name,omitempty" example:"Alice"           doc:"First name of the teacher"`
	LastName  string `json:"last_name,omitempty"  example:"Brown"           doc:"Last name of the techer"`
	Class     string `json:"class,omitempty"      example:"11C"             doc:"The class of the teacher"`
	Subject   string `json:"subject,omitempty"    example:"History"         doc:"Subject to teach"`
	Email     string `json:"email,omitempty"      example:"ac@example.com " doc:"Email"`
}

type Exec struct{}
