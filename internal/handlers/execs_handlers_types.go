// Package handlers - part of handlers but only for huma query paramaters
package handlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

type ExecsInput struct {
	Body struct {
		Execs []models.ExecInput `json:"execs" doc:"Execs"`
	}
}

type ExecsOutput struct {
	Body struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}
}

type ExecIDResponse struct {
	Body struct {
		Data models.Exec `json:"data"`
	}
}

type ExecPatchInput struct {
	Body struct {
		Exec models.ExecPatchBody `json:"exec" doc:"Exec patch body output"`
	}
}

type ExecPatchOutput struct {
	Body struct {
		Status string      `json:"status"`
		Data   models.Exec `json:"data"`
	}
}

type ExecsPatchInput struct {
	Body struct {
		Students []models.ExecPatchBody `json:"exec" doc:"Execs patch"`
	}
}
type ExecsPatchOutput struct {
	Body struct {
		Status string        `json:"status"`
		Data   []models.Exec `json:"data"`
	}
}

type DeleteExecsInput struct {
	IDn []int `query:"idn" example:"[104,106,103]" doc:"Execs IDn to delete"`
}
type DeleteExecsOutput struct {
	Body struct {
		Status string `json:"status"`
		ID     []int  `json:"id"`
	}
}

type ExecsLoginInput struct {
	Body struct {
		Exec models.ExecLoginInput `json:"execs" doc:"Execs"`
	}
}

type ExecsLoginOutput struct {
	Body struct {
		SetCookie http.Cookie `header:"Set-Cookie" json:"-"`
		Token     string      `                    json:"token"`
	}
}
