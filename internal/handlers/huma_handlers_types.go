// Package handlers - part of handlers but only for huma query paramaters

package handlers

import "github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

type TeachersInput struct {
	Body struct {
		Teachers []TeacherBody `json:"teachers" doc:"Teachers"`
	}
}
type TeacherBody struct {
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

type TeachersOutput struct {
	Body struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}
}

type TeacherIDResponse struct {
	Body struct {
		Data models.Teacher `json:"data"`
	}
}

type TeachersUpdateInput struct {
	Body struct {
		Teacher TeacherUpdateBody `json:"teacher" doc:"Teacher"`
	}
}
type TeachersUpdateOutput struct {
	Body struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}
}

type TeacherUpdateBody struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" doc:"First name of the teacher"`
	LastName  string `json:"last_name"  doc:"Last name of the techer"`
	Class     string `json:"class"      doc:"The class of the teacher"`
	Subject   string `json:"subject"    doc:"Subject to teach"`
	Email     string `json:"email"      doc:"Email"`
}

type TeachersPatchInput struct {
	Body struct {
		Teacher TeacherPatchBody `json:"teacher" doc:"Teacher"`
	}
}

type TeachersPatchOutput struct {
	Body struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}
}

type TeacherPatchBody struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name,omitempty" doc:"First name of the teacher"`
	LastName  string `json:"last_name,omitempty"  doc:"Last name of the techer"`
	Class     string `json:"class,omitempty"      doc:"The class of the teacher"`
	Subject   string `json:"subject,omitempty"    doc:"Subject to teach"`
	Email     string `json:"email,omitempty"      doc:"Email"`
}
