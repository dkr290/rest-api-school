package dataops

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

type Students struct {
	db *sql.DB
}

func NewStudentsDB(db *sql.DB) *Students {
	return &Students{
		db: db,
	}
}

func (t *Students) Insertstudents(st *models.Student) (int64, error) {
	stmt, err := t.db.Prepare(`INSERT INTO students
		            (first_name,last_name,email,class)
                VALUES(?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	sqlResp, err := stmt.Exec(
		st.FirstName,
		st.LastName,
		st.Email,
		st.Class,
	)
	if err != nil {
		return 0, err
	}
	lastID, err := sqlResp.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

func (t *Students) GetStudentByID(id int) (models.Student, error) {
	var student models.Student
	err := t.db.QueryRow("SELECT id, first_name, last_name ,email, class FROM students WHERE id = ?", id).
		Scan(
			&student.ID,
			&student.FirstName,
			&student.LastName,
			&student.Email,
			&student.Class,
		)
	if err == sql.ErrNoRows {
		return models.Student{}, fmt.Errorf("teacher not found %v", err)
	} else if err != nil {
		return models.Student{}, fmt.Errorf("error quering the database %v", err)
	}
	return student, nil
}

func (t *Students) GetAllStudents(params map[string]string, sortBy []string) (*sql.Rows, error) {
	query := "SELECT id, first_name,last_name,email,class FROM students WHERE 1=1"
	var args []any
	var orderByParts []string

	// Define a whitelist of allowed sortable columns
	allowedColumns := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
	}
	// filtering by map of params
	for param, dbField := range params {
		if dbField != "" {
			query += " AND " + param + " = ?"
			args = append(args, dbField)
		}
	}

	for _, criteria := range sortBy {
		parts := strings.Split(criteria, ":")
		if len(parts) == 2 {
			sortColumn := parts[0]
			sortOrder := strings.ToUpper(parts[1])

			if allowedColumns[sortColumn] && (sortOrder == "ASC" || sortOrder == "DESC") {
				orderByParts = append(orderByParts, fmt.Sprintf("%s %s", sortColumn, sortOrder))
			}
		}
	}

	if len(orderByParts) > 0 {
		query += " ORDER BY " + strings.Join(orderByParts, ", ")
	}

	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query error %v", err)
	}
	return rows, nil
}

func (t *Students) UpdateTeacher(id int, updatedTeacher models.Student) (models.Student, error) {
	var existingStudent models.Student

	row := t.db.QueryRow(
		"SELECT id ,first_name,last_name,email,class from teachers WHERE id = ?",
		id,
	)
	err := row.Scan(
		&existingStudent.ID,
		&existingStudent.FirstName,
		&existingStudent.LastName,
		&existingStudent.Email,
		&existingStudent.Class,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return models.Student{}, huma.Error500InternalServerError("Teacher not found", err)
		} else {
			return models.Student{}, huma.NewError(http.StatusNotFound, "unable to retreive data", err)
		}
	}

	updatedTeacher.ID = existingStudent.ID
	switch {
	}

	_, err = t.db.Exec(
		"UPDATE teachers SET first_name = ?, last_name = ? ,email = ? , class = ? WHERE id = ?  ",
		&updatedTeacher.FirstName,
		&updatedTeacher.LastName,
		&updatedTeacher.Email,
		&updatedTeacher.Class,
		&updatedTeacher.ID,
	)
	if err != nil {
		return models.Student{}, huma.Error500InternalServerError("Error updating teacher", err)
	}

	return updatedTeacher, nil
}

func (t *Teachers) PatchTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	var existingTeacher models.Teacher

	row := t.db.QueryRow(
		"SELECT id ,first_name,last_name,email,class,subject from teachers WHERE id = ?",
		id,
	)
	err := row.Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return models.Teacher{}, huma.Error500InternalServerError("Teacher not found", err)
		} else {
			return models.Teacher{}, huma.NewError(http.StatusNotFound, "unable to retreive data", err)
		}
	}

	// if updatedTeacher.FirstName != "" {
	// 	existingTeacher.FirstName = updatedTeacher.FirstName
	// }
	// if updatedTeacher.LastName != "" {
	// 	existingTeacher.LastName = updatedTeacher.LastName
	// }
	//
	// if updatedTeacher.Email != "" {
	// 	existingTeacher.Email = updatedTeacher.Email
	// }
	// if updatedTeacher.Class != "" {
	// 	existingTeacher.Class = updatedTeacher.Class
	// }
	// if updatedTeacher.Subject != ""  {
	// 	existingTeacher.Subject = updatedTeacher.Subject
	// }

	// apply updates using reflect package

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()

	updatedVal := reflect.ValueOf(updatedTeacher)

	for i := 0; i < teacherVal.NumField(); i++ {
		updatedField := updatedVal.Field(i)
		fieldName := teacherVal.Type().Field(i).Name

		// Check if the field is a string and not empty
		if updatedField.Kind() == reflect.String && updatedField.String() != "" {
			// Find the corresponding field in existingTeacher
			existingField := teacherVal.FieldByName(fieldName)

			// Check if the field exists and is settable
			if existingField.IsValid() && existingField.CanSet() {
				// Check if the field in existingTeacher is also a string (for safety)
				if existingField.Kind() == reflect.String {
					// Set the value
					existingField.SetString(updatedField.String())
				}
			}
		}
	}

	_, err = t.db.Exec(
		"UPDATE teachers SET first_name = ?, last_name = ? ,email = ? , class = ?,subject = ? WHERE id = ?  ",
		existingTeacher.FirstName,
		existingTeacher.LastName,
		existingTeacher.Email,
		existingTeacher.Class,
		existingTeacher.Subject,
		existingTeacher.ID,
	)
	if err != nil {
		return models.Teacher{}, huma.Error500InternalServerError("Error updating teacher", err)
	}

	return existingTeacher, nil
}

func (t *Teachers) DeleteTeacher(id int) error {
	result, err := t.db.Exec("DELETE from teachers WHERE id = ?", id)
	if err != nil {
		return huma.Error500InternalServerError(
			"Error deleting teacher",
			err,
		)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return huma.Error500InternalServerError(
			"Error retreiving delete result",
			err,
		)
	}

	if rowsAffected == 0 {
		return huma.Error404NotFound(
			"Teacher not found",
			err,
		)
	}

	return nil
}
