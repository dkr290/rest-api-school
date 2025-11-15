package handlers

import (
	"context"
	"sync"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/dataops"
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

func (h *ExecsHandlers) ExecAddHandler(ctx context.Context, input *struct{}) (*struct{}, error) {
	return nil, nil
}

func (h *ExecsHandlers) PatchExecsHandler(ctx context.Context, input *struct{}) (*struct{}, error) {
	return nil, nil
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
