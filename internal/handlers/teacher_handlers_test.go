package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

// Mock implementation of TeachersInf interface
type mockTeachersDB struct {
	teachers []models.Teacher
	teacher  models.Teacher
	err      error
}

// Implement all required interface methods
func (m *mockTeachersDB) InsertTeachers(t *models.Teacher) (int64, error) {
	return 0, m.err
}

func (m *mockTeachersDB) GetTeacherByID(id int) (models.Teacher, error) {
	return m.teacher, m.err
}

func (m *mockTeachersDB) GetAllTeachers(
	params map[string]string,
	sortBy []string,
) (*sql.Rows, error) {
	// For testing, you'll need to use a library like sqlmock or restructure to avoid sql.Rows
	// This is a limitation - sql.Rows cannot be easily mocked
	// Better approach: return []models.Teacher instead of *sql.Rows from the dataops layer
	return nil, m.err
}

func (m *mockTeachersDB) UpdateTeacher(id int, t models.Teacher) (models.Teacher, error) {
	return t, m.err
}

func (m *mockTeachersDB) PatchTeacher(id int, t models.Teacher) (models.Teacher, error) {
	return models.Teacher{}, m.err
}

func (m *mockTeachersDB) DeleteTeacher(id int) error {
	return m.err
}

func (m *mockTeachersDB) DeleteBulkTeachers(ids []int) ([]int, error) {
	return nil, m.err
}

func (m *mockTeachersDB) GetStudentsByTeacherID(id int) ([]models.Student, error) {
	return nil, m.err
}

func TestTeacherGetById(t *testing.T) {
	_, api := humatest.New(t)
	mockDB := &mockTeachersDB{
		teacher: models.Teacher{
			ID:        42,
			FirstName: "Jane",
			LastName:  "Small",
			Email:     "janesmall@example.com",
			Class:     "12C",
			Subject:   "History",
		},
		err: nil,
	}
	h := NewTeachersHandler(mockDB)
	huma.Register(api, huma.Operation{
		OperationID: "get-teacher",
		Method:      http.MethodGet,
		Path:        "/teachers/{id}",
		Summary:     "Get a teacher by ID",
		Description: "Get a teacher by ID.",
		Tags:        []string{"Teachers"},
	}, h.TeacherGet)

	resp := api.Get("/teachers/42")
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %v", resp.Code)
	}
	body := resp.Body.String()
	if !strings.Contains(body, "Jane") {
		t.Fatalf("Expected response to contain 'Jane', got: %s", body)
	}
	if !strings.Contains(body, "Small") {
		t.Fatalf("Expected response to contain 'Small', got: %s", body)
	}
	if !strings.Contains(body, "42") {
		t.Fatalf("Expected response to contain ID '42', got: %s", body)
	}
}

func TestUpdateTeacherHandler(t *testing.T) {
	_, api := humatest.New(t)
	mockDB := &mockTeachersDB{
		teacher: models.Teacher{
			ID:        42,
			FirstName: "Jane",
			LastName:  "Small",
			Email:     "janesmall@example.com",
			Class:     "12C",
			Subject:   "History",
		},
		err: nil,
	}
	h := NewTeachersHandler(mockDB)
	huma.Register(api, huma.Operation{
		OperationID: "update-teacher",
		Method:      http.MethodPut,
		Path:        "/teachers/{id}",
		Summary:     "Update all fields of a teacher",
		Description: "Update all fields of a teacher mandatory.",
		Tags:        []string{"Teachers"},
	}, h.UpdateTeacherHandler)

	resp := api.Put("/teachers/42", map[string]any{
		"teacher": map[string]any{
			"first_name": "Jane",
			"last_name":  "Small",
			"id":         42,
			"email":      "janesmall@example.com",
			"class":      "12C",
			"subject":    "History",
		},
	})
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %v", resp.Code)
	}

	body := resp.Body.String()
	if !strings.Contains(body, "Jane") {
		t.Fatalf("Expected response to contain 'Jane', got: %s", body)
	}
	if !strings.Contains(body, "Small") {
		t.Fatalf("Expected response to contain 'Small', got: %s", body)
	}
	if !strings.Contains(body, "42") {
		t.Fatalf("Expected response to contain ID '42', got: %s", body)
	}
}
