package dataops

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/logging"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/utils"
)

type Execs struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewExecsDB(db *sql.DB, logger *logging.Logger) *Execs {
	return &Execs{
		db:     db,
		logger: logger,
	}
}

func (e *Execs) InsertExecs(ex *models.Exec) (int64, error) {
	stmt, err := e.db.Prepare(utils.GenereateInsertQuery(models.Exec{}, "execs"))
	if err != nil {
		e.logger.Logging.Debugf("error prepare insert statement %v", err)
		return 0, e.logger.ErrorMessage("sql database insert exec error ")
	}
	defer stmt.Close()
	values := utils.GetStructValues(ex)

	sqlResp, err := stmt.Exec(values...)
	if err != nil {

		e.logger.Logging.Debugf("error insert exec to the database %v", err)
		return 0, e.logger.ErrorMessage("sql database error")
	}

	lastID, err := sqlResp.LastInsertId()
	if err != nil {

		e.logger.Logging.Debugf("eror get last insert teacher %v", err)
		return 0, e.logger.ErrorMessage("sql database error")
	}
	return lastID, nil
}

func (e *Execs) GetAllExecs(params map[string]string, sortBy []string) (*sql.Rows, error) {
	query := "SELECT id, first_name,last_name,email, username, user_created_at, inactive_status, role FROM execs WHERE 1=1"
	var args []any
	var orderByParts []string

	// Define a whitelist of allowed sortable columns
	allowedColumns := map[string]bool{
		"first_name":      true,
		"last_name":       true,
		"email":           true,
		"username":        true,
		"user_created_at": true,
		"inactive_status": true,
		"role":            true,
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

	rows, err := e.db.Query(query, args...)
	if err != nil {
		e.logger.Logging.Debugf("error retreiving the data %v", err)
		return nil, e.logger.ErrorMessage("error retrtreiving data")
	}
	return rows, nil
}

func (e *Execs) GetExecsByID(id int) (models.Exec, error) {
	var exec models.Exec

	err := e.db.QueryRow("SELECT id, first_name, last_name ,email, username, inactive_status, role FROM execs WHERE id = ?", id).
		Scan(
			&exec.ID,
			&exec.FirstName,
			&exec.LastName,
			&exec.Email,
			&exec.Username,
			&exec.InactiveStatus,
			&exec.Role,
		)
	if err == sql.ErrNoRows {
		e.logger.Logging.Debugf("error exec not found %v", err)
		return models.Exec{}, e.logger.ErrorMessage("sql exec error")
	} else if err != nil {
		e.logger.Logging.Debugf("error quring the database %v", err)
		return models.Exec{}, e.logger.ErrorMessage("sql exec error")
	}
	return exec, nil
}

func (e *Execs) PatchExec(id int, updatedExec models.Exec) (models.Exec, error) {
	var existingExec models.Exec

	row := e.db.QueryRow(
		"SELECT id ,first_name,last_name,email, username  from execs WHERE id = ?",
		id,
	)
	err := row.Scan(
		&existingExec.ID,
		&existingExec.FirstName,
		&existingExec.LastName,
		&existingExec.Email,
		&existingExec.Username,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			e.logger.Logging.Debugf("Exec not found %v", err)
			return models.Exec{}, e.logger.ErrorMessage("database error")
		} else {
			e.logger.Logging.Debugf("unable to retreive data %v", err)
			return models.Exec{}, e.logger.ErrorMessage("database error")
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

	execVal := reflect.ValueOf(&existingExec).Elem()

	updatedVal := reflect.ValueOf(updatedExec)

	for i := 0; i < execVal.NumField(); i++ {
		updatedField := updatedVal.Field(i)
		fieldName := execVal.Type().Field(i).Name

		// Check if the field is a string and not empty
		if updatedField.Kind() == reflect.String && updatedField.String() != "" {
			// Find the corresponding field in existingTeacher
			existingField := execVal.FieldByName(fieldName)

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

	_, err = e.db.Exec(
		"UPDATE execs SET first_name = ?, last_name = ? ,email = ?, username = ?  WHERE id = ?  ",
		existingExec.FirstName,
		existingExec.LastName,
		existingExec.Email,
		existingExec.Username,
		existingExec.ID,
	)
	if err != nil {
		e.logger.Logging.Debugf("error updating exec %v", err)
		return models.Exec{}, e.logger.ErrorMessage("database error")
	}

	return existingExec, nil
}

func (e *Execs) DeleteExec(id int) error {
	result, err := e.db.Exec("DELETE from execs WHERE id = ?", id)
	if err != nil {
		e.logger.Logging.Debugf("error deleting exec %v", err)
		return e.logger.ErrorMessage("database delete error")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		e.logger.Logging.Debugf("error retreiving delete result %v", err)
		return e.logger.ErrorMessage("error database delete operation")
	}

	if rowsAffected == 0 {
		e.logger.Logging.Debugf("exec not found %v", err)
		return e.logger.ErrorMessage("exec not found")
	}

	return nil
}
