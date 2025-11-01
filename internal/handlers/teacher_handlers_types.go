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
		Teachers []models.Teacher `json:"teachers" doc:"Teachers"`
	}
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
		Teacher models.TeacherUpdateBody `json:"teacher" doc:"Teacher"`
	}
}
type TeachersUpdateOutput struct {
	Body struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}
}

type TeacherPatchInput struct {
	Body struct {
		Teacher models.TeacherPatchBody `json:"teacher" doc:"Teacher"`
	}
}

type TeacherPatchOutput struct {
	Body struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}
}

type TeachersPatrchInput struct {
	Body struct {
		Teachers []models.TeacherPatchBody `json:"teachers" doc:"Teachers patch"`
	}
}
type TeachersPatchOutput struct {
	Body struct {
		Status string           `json:"status"`
		Data   []models.Teacher `json:"data"`
	}
}

type DeleteTeachersInput struct {
	IDn []int `query:"idn" example:"[104,106,103]" doc:"Teachers IDn to delete"`
}
type DeleteTeachersOutput struct {
	Body struct {
		Status string `json:"status"`
		ID     []int  `json:"id"`
	}
}
