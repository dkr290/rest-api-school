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
		Token string `                    json:"token"`
	}
	SetCookie http.Cookie `header:"Set-Cookie"`
}

type ExecLogoutOutput struct {
	Body struct {
		Status string `json:"status"`
	}
	SetCookie http.Cookie `header:"Set-Cookie"`
}

type ExecUpdatePasswordInput struct {
	Body struct {
		ID              int    `json:"id"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
}
type ExecUpdatePasswordOutput struct {
	Body struct {
		PasswordUpdated string `json:"password_updated"`
	}
	SetCookie http.Cookie `header:"Set-Cookie"`
}

type ExecsForgotPasswordInput struct {
	Body struct {
		Email string `json:"email"`
	}
}

type ExecsPasswordResetInput struct {
	ResetCode string `path:"resetcode" doc:"Password reset code"`
	Body      struct {
		NewPassword     string `json:"new_password" required:"true" minLength:"2" maxLength:"255" doc:"New password"`
		ConfirmPassword string `json:"confirm_password" required:"true" minLength:"2" maxLength:"255" doc:"Confirm password"`
	}
}

type PasswordresetOutput struct {
	Body struct {
		Data string `json:"data"`
	}
}
