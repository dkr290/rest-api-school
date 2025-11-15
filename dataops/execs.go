package dataops

import (
	"database/sql"

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

		e.logger.Logging.Debugf("error insert student to the database %v", err)
		return 0, e.logger.ErrorMessage("sql database error")
	}

	lastID, err := sqlResp.LastInsertId()
	if err != nil {

		e.logger.Logging.Debugf("eror get last insert teacher %v", err)
		return 0, e.logger.ErrorMessage("sql database error")
	}
	return lastID, nil
}

func (e *Execs) GetExecsByID(id int) (models.Exec, error) {
	var exec models.Exec

	err := e.db.QueryRow("SELECT id, first_name, last_name ,email, username,inactive_status, role FROM execs WHERE id = ?", id).
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
		return models.Exec{}, e.logger.ErrorMessage("sql student error")
	}
	return exec, nil
}
