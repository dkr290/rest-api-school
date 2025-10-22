package dataops

import (
	"database/sql"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

type DatabaseInf interface {
	InsertTeachers(*models.Teacher) (int64, error)
	GetTeacherByID(int) (models.Teacher, error)
	GetAllTeachers(map[string]string, []string) (*sql.Rows, error)
	UpdateTeacher(int, models.Teacher) (models.Teacher, error)
	PatchTeacher(int, models.Teacher) (models.Teacher, error)
	DeleteTeacher(int) error
	DeleteBulkTeachers([]int) ([]int, error)
	InsertStudents(*models.Student) (int64, error)
	GetStudentByID(int) (models.Student, error)
	GetAllStudents(map[string]string, []string) (*sql.Rows, error)
	UpdateStudent(int, models.Student) (models.Student, error)
	PatchiStudent(int, models.Student) (models.Student, error)
	DeleteStudent(int) error
}
