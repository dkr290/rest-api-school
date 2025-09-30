package dataops

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/logging"
)

type DatabaseInf interface {
	InsertTeachers(*models.Teacher) (int64, error)
	GetTeacherByID(int) (models.Teacher, error)
	GetAllTeachers(map[string]string, []string) (*sql.Rows, error)
	UpdateTeacher(int, models.Teacher) (models.Teacher, error)
	PatchTeacher(int, models.Teacher) (models.Teacher, error)
	DeleteTeacher(int) error
	DeleteBulkTeachers([]int) ([]int, error)
}

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
	stmt, err := t.db.Prepare(`INSERT INTO teachers
		            (first_name,last_name,email,class,subject)
                VALUES(?,?,?,?,?)`)
	if err != nil {
		return 0, t.logger.ErrorLogger(err, "error prepare insert statement")
	}
	defer stmt.Close()
	sqlResp, err := stmt.Exec(
		tm.FirstName,
		tm.LastName,
		tm.Email,
		tm.Class,
		tm.Subject,
	)
	if err != nil {
		return 0, t.logger.ErrorLogger(err, "error inseart teacher to the database")
	}
	lastID, err := sqlResp.LastInsertId()
	if err != nil {
		return 0, t.logger.ErrorLogger(err, "error get last insert teacher")
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
		return models.Teacher{}, t.logger.ErrorLogger(err, "teacher not found")
	} else if err != nil {
		return models.Teacher{}, t.logger.ErrorLogger(err, "error quering the database")
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
		return nil, t.logger.ErrorLogger(err, "error retreiving data")
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
			return models.Teacher{}, t.logger.ErrorLogger(err, "Teacher not found")
		} else {
			return models.Teacher{}, t.logger.ErrorLogger(err, "unable to retreive data")
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
		return models.Teacher{}, t.logger.ErrorLogger(err, "error updating teacher")
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
			return models.Teacher{}, t.logger.ErrorLogger(err, "Teacher not found")
		} else {
			return models.Teacher{}, t.logger.ErrorLogger(err, "unable to retreive data")
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
		return models.Teacher{}, t.logger.ErrorLogger(err, "Error updating teacher")
	}

	return existingTeacher, nil
}

func (t *Teachers) DeleteTeacher(id int) error {
	result, err := t.db.Exec("DELETE from teachers WHERE id = ?", id)
	if err != nil {
		return t.logger.ErrorLogger(err, "Error deleting teacher")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return t.logger.ErrorLogger(err, "Error retreiving deleted teacher")
	}

	if rowsAffected == 0 {
		return t.logger.ErrorLogger(err, "Teacher not found")
	}

	return nil
}

func (t *Teachers) DeleteBulkTeachers(idn []int) ([]int, error) {
	tx, err := t.db.Begin()
	if err != nil {
		log.Println(err)
		return nil, t.logger.ErrorLogger(err, "Error starting Transaction")
	}
	stmt, err := tx.Prepare("DELETE from teachers WHERE id = ?")
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
			log.Println(err)
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
