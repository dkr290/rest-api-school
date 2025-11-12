package dataops

import (
	"database/sql"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/logging"
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
