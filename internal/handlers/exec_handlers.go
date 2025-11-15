package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/dataops"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/utils"
)

type ExecsHandlers struct {
	mutex   sync.Mutex
	execsDB dataops.ExecsInf
}

func NewExecsHandler(tdb dataops.ExecsInf) *ExecsHandlers {
	return &ExecsHandlers{
		execsDB: tdb,
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

		exec := models.Exec{
			FirstName: newExec.FirstName,
			LastName:  newExec.LastName,
			Email:     newExec.Email,
			Username:  newExec.Username,
			Password:  newExec.Password,
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

func (h *ExecsHandlers) PatchExecsHandler(
	ctx context.Context,
	input *ExecPatchInput,
) (*ExecPatchOutput, error) {
	id := input.Body.Exec.ID
	if id <= 0 {
		return nil, huma.NewError(http.StatusBadRequest, "invalid student id", nil)
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
	}

	updatedExec, err := h.execsDB.PatchExecs(input.Body.Exec.ID, exec)
	if err != nil {
		return nil, err
	}
	resp := ExecPatchOutput{}
	resp.Body.Status = "Success"
	resp.Body.Data = updatedExec
	return &resp, nil
}

func (h *ExecsHandlers) UpdateExecHandler(
	ctx context.Context,
	input *ExecsUpdateInput,
) (*ExecsUpdateOutput, error) {
	id := input.Body.Exec.ID
	if id <= 0 {
		return nil, huma.NewError(http.StatusBadRequest, "invalid student id", nil)
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
	}

	updatedExec, err := h.execsDB.UpdateExec(input.Body.Exec.ID, exec)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, huma.Error404NotFound("not found", err)
		}
		return nil, huma.Error500InternalServerError("error update database", err)
	}
	resp := ExecsUpdateOutput{}
	resp.Body.Status = "Sucess"
	resp.Body.Data = updatedExec
	return &resp, nil
}

func (h *ExecsHandlers) ExecGetByIDHandler(
	ctx context.Context,
	input *struct{},
) (*struct{}, error) {
	return nil, nil
}

func (h *ExecsHandlers) ExecPatchByIDHandler(
	ctx context.Context,
	input *struct{},
) (*struct{}, error) {
	return nil, nil
}

func (h *ExecsHandlers) ExecDeleteByIDHandler(
	ctx context.Context,
	input *struct{},
) (*struct{}, error) {
	return nil, nil
}

func (h *ExecsHandlers) ExecPasswordChangeHandler(
	ctx context.Context,
	input *struct{},
) (*struct{}, error) {
	return nil, nil
}

func (h *ExecsHandlers) ExecLoginHandler(
	ctx context.Context,
	input *struct{},
) (*struct{}, error) {
	return nil, nil
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
