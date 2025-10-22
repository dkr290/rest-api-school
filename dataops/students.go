package dataops

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/logging"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/utils"
)

type Students struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewStudentsDB(db *sql.DB, logger *logging.Logger) *Students {
	return &Students{
		db:     db,
		logger: logger,
	}
}

func (t *Students) InsertStudents(st *models.Student) (int64, error) {
	stmt, err := t.db.Prepare(utils.GenereateInsertQuery(models.Student{}, "students"))
	if err != nil {
		t.logger.Logging.Errorf("error prepare insert statement %v", err)
		return 0, t.logger.ErrorLogger(err, "error prepare insert statement")
	}
	defer stmt.Close()
	values := utils.GetStructValues(st)

	sqlResp, err := stmt.Exec(values...)
	if err != nil {

		t.logger.Logging.Errorf("error insert teacher to the database %v", err)
		return 0, t.logger.ErrorLogger(err, "error inseart teacher to the database")
	}

	lastID, err := sqlResp.LastInsertId()
	if err != nil {

		t.logger.Logging.Errorf("eror get last insert teacher %v", err)
		return 0, t.logger.ErrorLogger(err, "error get last insert teacher")
	}
	return lastID, nil
}

// TODO: to check logging and if everythging is setup like teachers

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
		return models.Student{}, fmt.Errorf("student not found %v", err)
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

func (t *Students) UpdateStudent(id int, updatedStudent models.Student) (models.Student, error) {
	var existingStudent models.Student

	row := t.db.QueryRow(
		"SELECT id ,first_name,last_name,email,class from students WHERE id = ?",
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
			return models.Student{}, huma.Error500InternalServerError("Student not found", err)
		} else {
			return models.Student{}, huma.NewError(http.StatusNotFound, "unable to retreive data", err)
		}
	}

	updatedStudent.ID = existingStudent.ID
	switch {
	}

	_, err = t.db.Exec(
		"UPDATE students SET first_name = ?, last_name = ? ,email = ? , class = ? WHERE id = ?  ",
		&updatedStudent.FirstName,
		&updatedStudent.LastName,
		&updatedStudent.Email,
		&updatedStudent.Class,
		&updatedStudent.ID,
	)
	if err != nil {
		return models.Student{}, huma.Error500InternalServerError("Error updating student", err)
	}

	return updatedStudent, nil
}

func (t *Students) PatchiStudent(id int, updatedStudent models.Student) (models.Student, error) {
	var existingStudent models.Student

	row := t.db.QueryRow(
		"SELECT id ,first_name,last_name,email,class from students WHERE id = ?",
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
			return models.Student{}, huma.Error500InternalServerError("Student not found", err)
		} else {
			return models.Student{}, huma.NewError(http.StatusNotFound, "unable to retreive data", err)
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

	studentVal := reflect.ValueOf(&existingStudent).Elem()

	updatedVal := reflect.ValueOf(updatedStudent)

	for i := 0; i < studentVal.NumField(); i++ {
		updatedField := updatedVal.Field(i)
		fieldName := studentVal.Type().Field(i).Name

		// Check if the field is a string and not empty
		if updatedField.Kind() == reflect.String && updatedField.String() != "" {
			// Find the corresponding field in existingTeacher
			existingField := studentVal.FieldByName(fieldName)

			// Check if the field exists and is settable
			if existingField.IsValid() && existingField.CanSet() {
				// Check if the field in existingStudent is also a string (for safety)
				if existingField.Kind() == reflect.String {
					// Set the value
					existingField.SetString(updatedField.String())
				}
			}
		}
	}

	_, err = t.db.Exec(
		"UPDATE students SET first_name = ?, last_name = ? ,email = ? , class = ? WHERE id = ?  ",
		existingStudent.FirstName,
		existingStudent.LastName,
		existingStudent.Email,
		existingStudent.Class,
		existingStudent.ID,
	)
	if err != nil {
		return models.Student{}, huma.Error500InternalServerError("Error updating student", err)
	}

	return existingStudent, nil
}

func (t *Students) DeleteStudent(id int) error {
	result, err := t.db.Exec("DELETE from students WHERE id = ?", id)
	if err != nil {
		return huma.Error500InternalServerError(
			"Error deleting student",
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
			"Student not found",
			err,
		)
	}

	return nil
}

func (t *Students) DeleteBulkStudents(idn []int) ([]int, error) {
	tx, err := t.db.Begin()
	if err != nil {
		t.logger.Logging.Errorf("Error starting transaction %v", err)
		return nil, t.logger.ErrorLogger(err, "Error starting Transaction")
	}
	stmt, err := tx.Prepare("DELETE from students WHERE id = ?")
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, t.logger.ErrorLogger(err, "error preparing delete statement")
	}
	defer stmt.Close()

	var deletedIds []int

	for _, id := range idn {
		res, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			t.logger.Logging.Errorf("error deleting teacher %v", err)
			return nil, t.logger.ErrorLogger(err, "error deleting teacher")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return nil, t.logger.ErrorLogger(err, "error retreiving delete result")
		}
		// if teacher was deleted then add ID to the deletedIDs slice
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			return nil, t.logger.ErrorMessage(
				fmt.Sprintf("ID %d does not exists,  doing rollback...", id),
			)
		}
	}
	// commit changes
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, t.logger.ErrorLogger(err, "error commiting the transaction")
	}
	if len(deletedIds) < 1 {
		return nil, t.logger.ErrorMessage("none of the id exists")
	}

	return deletedIds, err
}
