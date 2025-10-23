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
		t.logger.Logging.Debugf("error prepare insert statement %v", err)
		return 0, t.logger.ErrorLogger(err, "sql database insert student error ")
	}
	defer stmt.Close()
	values := utils.GetStructValues(st)

	sqlResp, err := stmt.Exec(values...)
	if err != nil {

		t.logger.Logging.Debugf("error insert student to the database %v", err)
		return 0, t.logger.ErrorLogger(err, "sql database error")
	}

	lastID, err := sqlResp.LastInsertId()
	if err != nil {

		t.logger.Logging.Debugf("eror get last insert teacher %v", err)
		return 0, t.logger.ErrorLogger(err, "sql database error")
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
		t.logger.Logging.Debugf("error student not found %v", err)
		return models.Student{}, t.logger.ErrorLogger(err, "sql student error")
	} else if err != nil {
		t.logger.Logging.Debugf("error quring the database %v", err)
		return models.Student{}, t.logger.ErrorLogger(err, "sql student error")
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
		t.logger.Logging.Debugf("error retreiving the data %v", err)
		return nil, t.logger.ErrorLogger(err, "error retrtreiving data")
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
			t.logger.Logging.Debugf("student not found %v", err)
			return models.Student{}, t.logger.ErrorLogger(err, "database retreive data error")
		} else {
			t.logger.Logging.Debugf("unable to retreive the data %v", err)
			return models.Student{}, t.logger.ErrorLogger(err, "sql error")
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
		t.logger.Logging.Debugf("error updating the student database %v", err)
		return models.Student{}, t.logger.ErrorLogger(err, "database error")
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
			t.logger.Logging.Debugf("Student not found %v", err)
			return models.Student{}, t.logger.ErrorMessage("database error")
		} else {
			t.logger.Logging.Debugf("unable to retreive data %v", err)
			return models.Student{}, t.logger.ErrorMessage("database error")
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
		t.logger.Logging.Debugf("error updating student %v", err)
		return models.Student{}, t.logger.ErrorMessage("database error")
	}

	return existingStudent, nil
}

func (t *Students) DeleteStudent(id int) error {
	result, err := t.db.Exec("DELETE from students WHERE id = ?", id)
	if err != nil {
		t.logger.Logging.Debugf("error deleting student %v", err)
		return t.logger.ErrorMessage("database delete error")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.logger.Logging.Debugf("error retreiving delete result %v", err)
		return t.logger.ErrorMessage("error database delete operation")
	}

	if rowsAffected == 0 {
		t.logger.Logging.Debugf("student not found %v", err)
		return t.logger.ErrorMessage("student not found")
	}

	return nil
}

func (t *Students) DeleteBulkStudents(idn []int) ([]int, error) {
	tx, err := t.db.Begin()
	if err != nil {
		t.logger.Logging.Errorf("Error starting transaction %v", err)
		return nil, t.logger.ErrorMessage("databae error")
	}
	stmt, err := tx.Prepare("DELETE from students WHERE id = ?")
	if err != nil {
		t.logger.Logging.Errorf("delete error and preparing delete statement %v", err)
		tx.Rollback()
		return nil, t.logger.ErrorMessage("error deleting")
	}
	defer stmt.Close()

	var deletedIds []int

	for _, id := range idn {
		res, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			t.logger.Logging.Debugf("error deleting student %v", err)
			return nil, t.logger.ErrorMessage("error database deleting student")
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
