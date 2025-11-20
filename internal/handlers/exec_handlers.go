package handlers

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/config"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/dataops"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/logging"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/utils"
	"golang.org/x/crypto/argon2"
)

type ExecsHandlers struct {
	mutex   sync.Mutex
	execsDB dataops.ExecsInf
	logger  *logging.Logger
	conf    config.Config
}

func NewExecsHandler(
	tdb dataops.ExecsInf,
	logger *logging.Logger,
	conf config.Config,
) *ExecsHandlers {
	return &ExecsHandlers{
		execsDB: tdb,
		logger:  logger,
		conf:    conf,
	}
}

func (h *ExecsHandlers) ExecGetHandler(ctx context.Context, input *struct {
	ID int `path:"id"`
},
) (*ExecIDResponse, error) {
	resp := ExecIDResponse{}

	exec, err := h.execsDB.GetExecsByID(input.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Error quering database", err)
	}

	resp.Body.Data = exec
	return &resp, nil
}

func (h *ExecsHandlers) ExecAddHandler(
	ctx context.Context,
	input *ExecsInput,
) (*ExecsOutput, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	addedExecs := make([]models.Exec, len(input.Body.Execs))

	for i, newExec := range input.Body.Execs {

		err := utils.EmailCheck(newExec.Email)
		if err != nil {
			return nil, huma.Error400BadRequest(
				"Invalid mail format",
				fmt.Errorf("invalid email: %s", newExec.Email),
			)
		}
		encodedPass, err := utils.PasswordHash(newExec.Password)
		if err != nil {
			h.logger.Logging.Errorf("failed to generate salt %v", err)
			return nil, huma.Error400BadRequest("error adding data")
		}

		exec := models.Exec{
			FirstName: newExec.FirstName,
			LastName:  newExec.LastName,
			Email:     newExec.Email,
			Username:  newExec.Username,
			Password:  encodedPass,
			Role:      newExec.Role,
		}
		id, err := h.execsDB.InsertExecs(&exec)
		if err != nil {
			return nil, huma.Error500InternalServerError(
				"Error adding to the database",
				err,
			)
		}
		exec.ID = int(id)
		addedExecs[i] = exec
	}

	resp := &ExecsOutput{}
	resp.Body.Status = "Success"
	resp.Body.Count = len(addedExecs)
	resp.Body.Data = addedExecs
	return resp, nil
}

func (e *ExecsHandlers) ExecsGetHandler(
	ctx context.Context,
	input *models.ExecsQueryInput,
) (*ExecsOutput, error) {
	response := ExecsOutput{}

	params := map[string]string{
		"first_name": input.FirstName,
		"last_name":  input.LastName,
		"email":      input.Email,
		"username":   input.Username,
		"role":       input.Role,
	}

	sortBy := input.SortBy
	// filtering by params basically with query parameters anf filtering
	rows, err := e.execsDB.GetAllExecs(params, sortBy)
	if err != nil {
		return nil, huma.Error500InternalServerError("Error quering database", err)
	}

	execsList := make([]models.Exec, 0)

	for rows.Next() {
		var exec models.Exec
		err = rows.Scan(
			&exec.ID,
			&exec.FirstName,
			&exec.LastName,
			&exec.Email,
			&exec.Username,
			&exec.UserCreatedAt,
			&exec.InactiveStatus,
			&exec.Role,
		)
		if err != nil {
			return nil, huma.Error500InternalServerError("Error scanning database results", err)
		}
		execsList = append(execsList, exec)
	}
	defer rows.Close()

	response.Body.Status = "Sucess"
	response.Body.Count = len(execsList)
	response.Body.Data = execsList
	return &response, nil
}

func (h *ExecsHandlers) PatchExecsHandler(
	ctx context.Context,
	input *ExecPatchInput,
) (*ExecPatchOutput, error) {
	id := input.Body.Exec.ID
	if id <= 0 {
		return nil, huma.NewError(http.StatusBadRequest, "invalid exec id", nil)
	}

	email := input.Body.Exec.Email

	err := utils.EmailCheck(email)
	if err != nil {
		return nil, huma.Error400BadRequest(
			"Invalid mail format",
			fmt.Errorf("invalid email: %s", email),
		)
	}

	exec := models.Exec{
		ID:        input.Body.Exec.ID,
		FirstName: input.Body.Exec.FirstName,
		LastName:  input.Body.Exec.LastName,
		Email:     input.Body.Exec.Email,
		Username:  input.Body.Exec.Username,
	}

	updatedExec, err := h.execsDB.PatchExec(input.Body.Exec.ID, exec)
	if err != nil {
		return nil, err
	}
	resp := ExecPatchOutput{}
	resp.Body.Status = "Success"
	resp.Body.Data = updatedExec
	return &resp, nil
}

func (h *ExecsHandlers) ExecDeleteByIDHandler(
	ctx context.Context,
	input *struct {
		ID int `path:"id"`
	},
) (*struct {
	Body struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}
}, error,
) {
	err := h.execsDB.DeleteExec(input.ID)
	if err != nil {
		return nil, err
	}
	output := &struct {
		Body struct {
			Status string `json:"status"`
			ID     int    `json:"id"`
		}
	}{
		Body: struct {
			Status string `json:"status"`
			ID     int    `json:"id"`
		}{
			Status: "Exec deleted sucessfully",
			ID:     input.ID,
		},
	}
	return output, err
}

func (h *ExecsHandlers) ExecPasswordChangeHandler(
	ctx context.Context,
	input *struct{},
) (*struct{}, error) {
	return nil, nil
}

func (h *ExecsHandlers) ExecLoginHandler(
	ctx context.Context,
	input *ExecsLoginInput,
) (*ExecsLoginOutput, error) {
	// Data Validation

	exec := input.Body.Exec
	if exec.Username == "" || exec.Password == "" {
		h.logger.Logging.Debugf(
			"Invalid or blank username or password for username=%s and password=%s",
			exec.Username,
			exec.Password,
		)
		return nil, huma.Error400BadRequest("Invalid or blank username or password")
	}

	// Search for the user if the user actually exists
	exists, err, passFromDB := h.execsDB.SearchUsername(exec.Username)
	if err != nil {
		return nil, huma.Error404NotFound("database error:", err)
	}
	if !exists {
		return nil, huma.Error404NotFound("user not found")
	}
	inactive, err := h.execsDB.IsInactiveUser(exec.Username)
	if err != nil {
		return nil, err
	}

	if inactive {
		return nil, huma.Error403Forbidden("user inactive")
	}

	// verify password
	ps := strings.Split(passFromDB, ".")
	if len(ps) != 2 {
		h.logger.Logging.Error("invalid encoded hash format")
		return nil, huma.Error400BadRequest("invalid encoded hash format")
	}

	saltBase := ps[0]
	hashedPasswordBase64 := ps[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase)
	if err != nil {
		h.logger.Logging.Error("failed to decode the salt")
		return nil, huma.Error400BadRequest("invalid encoded hash format")

	}
	hashedDBPassword, err := base64.StdEncoding.DecodeString(hashedPasswordBase64)
	if err != nil {
		h.logger.Logging.Error("failed to decode the hashed password")
		return nil, huma.Error400BadRequest("invalid encoded hash format for password")

	}

	hashInputPassword := argon2.IDKey([]byte(exec.Password), salt, 1, 64*1024, 4, 32)

	if len(hashedDBPassword) != len(hashInputPassword) {
		h.logger.Logging.Error("incorrect password")
		return nil, huma.Error403Forbidden("incorrect password")
	}

	if subtle.ConstantTimeCompare(hashedDBPassword, hashInputPassword) == 1 {
	} else {
		h.logger.Logging.Error("incorrect password")
		return nil, huma.Error403Forbidden("incorrect password")

	}
	// generate token
	user, err := h.execsDB.GetLoginDetailsForUsername(exec.Username)
	if err != nil {
		h.logger.Logging.Errorf("error on get user details %v", err)
		return nil, huma.Error403Forbidden("incorrect user get from db")

	}
	tokenString, err := utils.SighnToken(
		fmt.Sprintf("%v", user.ID),
		user.Username,
		user.Role,
		h.conf,
	)
	if err != nil {
		h.logger.Logging.Errorf("Could not create login token %v", err)
		return nil, huma.Error500InternalServerError("Could not create login token", err)
	}

	// Send token as responce or as a cookie
	// how to make it as cookie in huma

	out := &ExecsLoginOutput{}
	out.Body.Token = tokenString
	out.Body.SetCookie = http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}

	return out, nil
}

func (h *ExecsHandlers) LogoutExecsHandler(
	ctx context.Context,
	input *struct{},
) (*struct{}, error) {
	return nil, nil
}

func (h *ExecsHandlers) ForgotpasswordExecsHandler(
	ctx context.Context,
	input *struct{},
) (*struct{}, error) {
	return nil, nil
}

func (h *ExecsHandlers) PasswordresetExecsHandler(
	ctx context.Context,
	input *struct{},
) (*struct{}, error) {
	return nil, nil
}
