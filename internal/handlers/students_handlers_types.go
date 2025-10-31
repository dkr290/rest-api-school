// Package handlers - part of handlers but only for huma query paramaters
package handlers

import "github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"

type StudentsInput struct {
	Body struct {
		Students []models.StudentInput `json:"students" doc:"Students"`
	}
}

type StudentsOutput struct {
	Body struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}
}

type StudentIDResponse struct {
	Body struct {
		Data models.Student `json:"data"`
	}
}

type StudentsUpdateInput struct {
	Body struct {
		Student models.StudentUpdateBody `json:"student" doc:"Student update body output"`
	}
}
type StudentsUpdateOutput struct {
	Body struct {
		Status string         `json:"status"`
		Data   models.Student `json:"data"`
	}
}

type StudentPatchInput struct {
	Body struct {
		Student models.StudentPatchBody `json:"student" doc:"Student patch body output"`
	}
}

type StudentPatchOutput struct {
	Body struct {
		Status string         `json:"status"`
		Data   models.Student `json:"data"`
	}
}

type StudentsPatchInput struct {
	Body struct {
		Students []models.StudentPatchBody `json:"students" doc:"Students patch"`
	}
}
type StudentsPatchOutput struct {
	Body struct {
		Status string           `json:"status"`
		Data   []models.Student `json:"data"`
	}
}

type DeleteStudentsInput struct {
	IDn []int `query:"idn" example:"[104,106,103]" doc:"Students IDn to delete"`
}
type DeleteStudentsOutput struct {
	Body struct {
		Status string `json:"status"`
		ID     []int  `json:"id"`
	}
}
