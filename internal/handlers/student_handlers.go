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

type StudentHandlers struct {
	mutex      sync.Mutex
	studentsDB dataops.StudentInf
}

func NewStudentsHandler(sdb dataops.StudentInf) *StudentHandlers {
	return &StudentHandlers{
		studentsDB: sdb,
	}
}

func (h *StudentHandlers) StudentGet(ctx context.Context, input *struct {
	ID int `path:"id"`
},
) (*StudentIDResponse, error) {
	resp := StudentIDResponse{}

	student, err := h.studentsDB.GetStudentByID(input.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Error quering database", err)
	}

	resp.Body.Data = student
	return &resp, nil
}

func (h *StudentHandlers) StudentsGet(
	ctx context.Context,
	input *models.StudentsQueryInput,
) (*StudentsOutput, error) {
	response := StudentsOutput{}

	params := map[string]string{
		"first_name": input.FirstName,
		"last_name":  input.LastName,
		"email":      input.Email,
		"class":      input.Class,
	}

	sortBy := input.SortBy
	// filtering by params basically with query parameters anf filtering
	rows, err := h.studentsDB.GetAllStudents(params, sortBy)
	if err != nil {
		return nil, huma.Error500InternalServerError("Error quering database", err)
	}

	studentsList := make([]models.Student, 0)

	for rows.Next() {
		var student models.Student
		err = rows.Scan(
			&student.ID,
			&student.FirstName,
			&student.LastName,
			&student.Email,
			&student.Class,
		)
		if err != nil {
			return nil, huma.Error500InternalServerError("Error scanning database results", err)
		}
		studentsList = append(studentsList, student)
	}
	defer rows.Close()

	response.Body.Status = "Sucess"
	response.Body.Count = len(studentsList)
	response.Body.Data = studentsList
	return &response, nil
}

func (h *StudentHandlers) StudentsAdd(
	ctx context.Context,
	input *StudentsInput,
) (*StudentsOutput, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	addedStudents := make([]models.Student, len(input.Body.Students))

	for i, newStudent := range input.Body.Students {

		err := utils.EmailCheck(newStudent.Email)
		if err != nil {
			return nil, huma.Error400BadRequest(
				"Invalid mail format",
				fmt.Errorf("invalid email: %s", newStudent.Email),
			)
		}

		student := models.Student{
			FirstName: newStudent.FirstName,
			LastName:  newStudent.LastName,
			Email:     newStudent.Email,
			Class:     newStudent.Class,
		}
		id, err := h.studentsDB.InsertStudents(&student)
		if err != nil {
			return nil, huma.Error500InternalServerError(
				"Error adding to the database",
				err,
			)
		}
		student.ID = int(id)
		addedStudents[i] = student
	}

	resp := &StudentsOutput{}
	resp.Body.Status = "Success"
	resp.Body.Count = len(addedStudents)
	resp.Body.Data = addedStudents
	return resp, nil
}

func (h *StudentHandlers) UpdateStudentHandler(
	ctx context.Context,
	input *StudentsUpdateInput,
) (*StudentsUpdateOutput, error) {
	id := input.Body.Student.ID
	if id <= 0 {
		return nil, huma.NewError(http.StatusBadRequest, "invalid student id", nil)
	}
	email := input.Body.Student.Email
	err := utils.EmailCheck(email)
	if err != nil {
		return nil, huma.Error400BadRequest(
			"Invalid mail format",
			fmt.Errorf("invalid email: %s", email),
		)
	}

	student := models.Student{
		ID:        input.Body.Student.ID,
		FirstName: input.Body.Student.FirstName,
		LastName:  input.Body.Student.LastName,
		Email:     input.Body.Student.Email,
		Class:     input.Body.Student.Class,
	}

	updatedStudent, err := h.studentsDB.UpdateStudent(input.Body.Student.ID, student)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, huma.Error404NotFound("not found", err)
		}
		return nil, huma.Error500InternalServerError("error update database", err)
	}
	resp := StudentsUpdateOutput{}
	resp.Body.Status = "Sucess"
	resp.Body.Data = updatedStudent
	return &resp, nil
}

func (h *StudentHandlers) PatchStudentHandler(
	ctx context.Context,
	input *StudentPatchInput,
) (*StudentPatchOutput, error) {
	id := input.Body.Student.ID
	if id <= 0 {
		return nil, huma.NewError(http.StatusBadRequest, "invalid student id", nil)
	}

	email := input.Body.Student.Email

	err := utils.EmailCheck(email)
	if err != nil {
		return nil, huma.Error400BadRequest(
			"Invalid mail format",
			fmt.Errorf("invalid email: %s", email),
		)
	}

	student := models.Student{
		ID:        input.Body.Student.ID,
		FirstName: input.Body.Student.FirstName,
		LastName:  input.Body.Student.LastName,
		Email:     input.Body.Student.Email,
		Class:     input.Body.Student.Class,
	}
	// TODO: to add better error handling here
	updatedStudent, err := h.studentsDB.PatchiStudent(input.Body.Student.ID, student)
	if err != nil {
		return nil, err
	}
	resp := StudentPatchOutput{}
	resp.Body.Status = "Success"
	resp.Body.Data = updatedStudent
	return &resp, nil
}

func (h *StudentHandlers) DeleteStudentHandler(ctx context.Context, input *struct {
	ID int `path:"id"`
},
) (*struct {
	Body struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}
}, error,
) {
	err := h.studentsDB.DeleteStudent(input.ID)
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
			Status: "Student deleted sucessfully",
			ID:     input.ID,
		},
	}

	return output, err
}

func (h *StudentHandlers) PatchStudentsHandler(
	ctx context.Context,
	input *StudentsPatchInput,
) (*StudentsPatchOutput, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	patchedStudents := make([]models.Student, len(input.Body.Students))

	for i, newStudent := range input.Body.Students {

		err := utils.EmailCheck(newStudent.Email)
		if err != nil {
			return nil, huma.Error400BadRequest(
				"Invalid mail format",
				fmt.Errorf("invalid email: %s", newStudent.Email),
			)
		}

		student := models.Student(newStudent)
		t, err := h.studentsDB.PatchiStudent(newStudent.ID, student)
		if err != nil {
			return nil, err
		}
		patchedStudents[i] = t
	}

	resp := &StudentsPatchOutput{}
	resp.Body.Status = "Success"
	resp.Body.Data = patchedStudents
	return resp, nil
}

func (h *StudentHandlers) DeleteStudentsHandler(
	ctx context.Context,
	input *DeleteStudentsInput,
) (*DeleteStudentsOutput, error) {
	respIDn, err := h.studentsDB.DeleteBulkStudents(input.IDn)
	if err != nil {
		return nil, err
	}

	resp := &DeleteStudentsOutput{}
	resp.Body.Status = "Success"
	resp.Body.ID = respIDn

	return resp, nil
}
