package dataops

import (
	"database/sql"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

type TeachersInf interface {
	InsertTeachers(*models.Teacher) (int64, error)
	GetTeacherByID(int) (models.Teacher, error)
	GetAllTeachers(map[string]string, []string) (*sql.Rows, error)
	UpdateTeacher(int, models.Teacher) (models.Teacher, error)
	PatchTeacher(int, models.Teacher) (models.Teacher, error)
	DeleteTeacher(int) error
	DeleteBulkTeachers([]int) ([]int, error)
	GetStudentsByTeacherID(int) ([]models.Student, error)
}
type StudentInf interface {
	InsertStudents(*models.Student) (int64, error)
	GetStudentByID(int) (models.Student, error)
	GetAllStudents(map[string]string, []string) (*sql.Rows, error)
	UpdateStudent(int, models.Student) (models.Student, error)
	PatchiStudent(int, models.Student) (models.Student, error)
	DeleteStudent(int) error
	DeleteBulkStudents([]int) ([]int, error)
}

type ExecsInf interface {
	InsertExecs(*models.Exec) (int64, error)
	GetExecsByID(int) (models.Exec, error)
	GetAllExecs(map[string]string, []string) (*sql.Rows, error)
	PatchExec(int, models.Exec) (models.Exec, error)
	DeleteExec(int) error
	SearchUsername(string) (bool, error, string)
	IsInactiveUser(string) (bool, error)
	GetLoginDetailsForUsername(string) (models.Exec, error)
	GetUserPasswordFromId(int) (string, string, string, error)
	UpdatePasswordChange(int, string) error
	GetIdFromEmail(string) (models.Exec, error)
}
