// Package dataops for database operations
package dataops

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/logging"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/utils"
)

type Teachers struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewTeachersDB(db *sql.DB, logger *logging.Logger) *Teachers {
	return &Teachers{
		db:     db,
		logger: logger,
	}
}

func (t *Teachers) InsertTeachers(tm *models.Teacher) (int64, error) {
	stmt, err := t.db.Prepare(utils.GenereateInsertQuery(models.Teacher{}, "teachers"))
	if err != nil {
		t.logger.Logging.Debugf("error prepare insert statement %v", err)
		return 0, t.logger.ErrorMessage("error database insert statement")
	}
	defer stmt.Close()

	values := utils.GetStructValues(tm)
	sqlResp, err := stmt.Exec(values...)
	if err != nil {
		t.logger.Logging.Errorf("error insert teacher to the database %v", err)
		return 0, t.logger.ErrorMessage("error database teacher insert")
	}
	lastID, err := sqlResp.LastInsertId()
	if err != nil {
		t.logger.Logging.Errorf("eror get last insert teacher %v", err)
		return 0, t.logger.ErrorMessage("error database")
	}
	return lastID, nil
}

func (t *Teachers) GetTeacherByID(id int) (models.Teacher, error) {
	var teacher models.Teacher
	err := t.db.QueryRow("SELECT id, first_name, last_name ,email, class, subject FROM teachers WHERE id = ?", id).
		Scan(
			&teacher.ID,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.Email,
			&teacher.Class,
			&teacher.Subject,
		)
	if err == sql.ErrNoRows {
		t.logger.Logging.Debugf("teacher not found %v", err)
		return models.Teacher{}, t.logger.ErrorMessage("teacher not found")
	} else if err != nil {
		t.logger.Logging.Debugf("error quering the database %v", err)
		return models.Teacher{}, t.logger.ErrorMessage("error quering the database error")
	}
	return teacher, nil
}

func (t *Teachers) GetAllTeachers(params map[string]string, sortBy []string) (*sql.Rows, error) {
	query := "SELECT id, first_name,last_name,email,class,subject FROM teachers WHERE 1=1"
	var args []any
	var orderByParts []string

	// Define a whitelist of allowed sortable columns
	allowedColumns := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
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
		t.logger.Logging.Debugf("error retreiving data %v", err)
		return nil, t.logger.ErrorMessage("error retreiving data")
	}
	return rows, nil
}

func (t *Teachers) UpdateTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
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
			t.logger.Logging.Debugf("techer not found %v", err)
			return models.Teacher{}, t.logger.ErrorMessage("sql error")
		} else {
			t.logger.Logging.Debugf("unable to retreive the data %v", err)
			return models.Teacher{}, t.logger.ErrorMessage("sql error")
		}
	}

	updatedTeacher.ID = existingTeacher.ID
	switch {
	}

	_, err = t.db.Exec(
		"UPDATE teachers SET first_name = ?, last_name = ? ,email = ? , class = ?,subject = ? WHERE id = ?  ",
		&updatedTeacher.FirstName,
		&updatedTeacher.LastName,
		&updatedTeacher.Email,
		&updatedTeacher.Class,
		&updatedTeacher.Subject,
		&updatedTeacher.ID,
	)
	if err != nil {
		t.logger.Logging.Debugf("errr updating the teacher database %v", err)
		return models.Teacher{}, t.logger.ErrorMessage("error teacher database error")
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
			t.logger.Logging.Warnf("teacher not found %v", err)
			return models.Teacher{}, t.logger.ErrorMessage("Teacher not found")
		} else {
			t.logger.Logging.Debugf("unable to retreive data %v", err)
			return models.Teacher{}, t.logger.ErrorMessage("unable to retreive data")
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
		t.logger.Logging.Debugf("Error updating teacher %v", err)
		return models.Teacher{}, t.logger.ErrorMessage("Error updating teacher")
	}

	return existingTeacher, nil
}

func (t *Teachers) DeleteTeacher(id int) error {
	result, err := t.db.Exec("DELETE from teachers WHERE id = ?", id)
	if err != nil {
		t.logger.Logging.Debugf("error deleting teacher -  %v", err)
		return t.logger.ErrorMessage("Error deleting teacher")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.logger.Logging.Debugf("error retreiving deleted teacher %v", err)
		return t.logger.ErrorMessage("Error retreiving deleted teacher")
	}

	if rowsAffected == 0 {
		t.logger.Logging.Debugf("teacher not found %v", err)
		return t.logger.ErrorMessage("Teacher not found")
	}

	return nil
}

func (t *Teachers) DeleteBulkTeachers(idn []int) ([]int, error) {
	tx, err := t.db.Begin()
	if err != nil {
		t.logger.Logging.Errorf("Error starting transaction %v", err)
		return nil, t.logger.ErrorLogger(err, "Error starting Transaction")
	}
	stmt, err := tx.Prepare("DELETE from teachers WHERE id = ?")
	if err != nil {
		t.logger.Logging.Debugf("error preparing delete statement %v", err)
		tx.Rollback()
		return nil, t.logger.ErrorMessage("error preparing delete statement")
	}
	defer stmt.Close()

	var deletedIds []int

	for _, id := range idn {
		res, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			t.logger.Logging.Errorf("error deleting teacher %v", err)
			return nil, t.logger.ErrorMessage("error deleting teacher")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			t.logger.Logging.Debugf("error retreiving delete result %v", err)
			log.Println(err)
			return nil, t.logger.ErrorMessage("error retreiving delete result")
		}
		// if teacher was deleted then add ID to the deletedIDs slice
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
			t.logger.Logging.Debugf("ID %d does not exists", id)
			tx.Rollback()
			return nil, t.logger.ErrorMessage("ID does not exists,  doing rollback...")

		}
	}
	// commit changes
	err = tx.Commit()
	if err != nil {
		t.logger.Logging.Debugf("error commiting the transaction %v", err)
		return nil, t.logger.ErrorMessage("error commiting the transaction")
	}
	if len(deletedIds) < 1 {
		return nil, t.logger.ErrorMessage("none of the id exists")
	}

	return deletedIds, err
}

func (t *Teachers) GetStudentsByTeacherID(id int) ([]models.Student, error) {
	query := `SELECT id,first_name,last_name,email,class FROM students where class=(SELECT class from teachers WHERE id = ?)`

	var students []models.Student
	rows, err := t.db.Query(query, id)
	if err != nil {
		t.logger.Logging.Debugf("error while execute query %v", err)
		return nil, t.logger.ErrorMessage("error execute query")
	}

	defer rows.Close()
	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.ID,
			&student.FirstName,
			&student.LastName,
			&student.Email,
			&student.Class,
		)
		if err != nil {
			return nil, t.logger.ErrorMessage("error fetching the database")
		}
		students = append(students, student)

	}
	err = rows.Err()
	if err != nil {
		return nil, t.logger.ErrorLogger(err, "rows error")
	}
	return students, nil
}
