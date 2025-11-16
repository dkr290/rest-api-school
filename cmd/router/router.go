// Package router - all huma register routes
package router

import (
	"database/sql"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/dataops"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/handlers"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/logging"
)

func Router(db *sql.DB, debugFlag bool) *http.ServeMux {
	llogger := logging.Init(debugFlag)
	router := http.NewServeMux()
	teachersDB := dataops.NewTeachersDB(db, llogger)
	studentsDB := dataops.NewStudentsDB(db, llogger)
	execDB := dataops.NewExecsDB(db, llogger)

	teacherHandler := handlers.NewTeachersHandler(teachersDB)
	studetnsHandler := handlers.NewStudentsHandler(studentsDB)
	execHandler := handlers.NewExecsHandler(execDB)

	api := humago.New(router, huma.DefaultConfig("My API", "1.0.0"))

	huma.Get(api, "/", teacherHandler.RootHandler)

	routesTeachers(api, teacherHandler)

	routesStudents(api, studetnsHandler)

	routesExec(api, execHandler)

	return router
}

func routesTeachers(api huma.API, teacherHandler *handlers.TeacherHandlers) {
	huma.Register(api, huma.Operation{
		OperationID: "get-teacher",
		Method:      http.MethodGet,
		Path:        "/teachers/{id}",
		Summary:     "Get a teacher",
		Description: "Get a teacher by ID.",
		Tags:        []string{"Teachers"},
	}, teacherHandler.TeacherGet)

	huma.Register(api, huma.Operation{
		OperationID: "post-teachers",
		Method:      http.MethodPost,
		Path:        "/teachers",
		Summary:     "Create teachers",
		Description: "Create teachers.",
		Tags:        []string{"Teachers"},
	}, teacherHandler.TeachersAdd)

	huma.Register(api, huma.Operation{
		OperationID: "get-teachers",
		Method:      http.MethodGet,
		Path:        "/teachers",
		Summary:     "Get all teachers",
		Description: "Get all teachers or with filtering.",
		Tags:        []string{"Teachers"},
	}, teacherHandler.TeachersGet)

	huma.Register(api, huma.Operation{
		OperationID: "update-teacher",
		Method:      http.MethodPut,
		Path:        "/teachers/{id}",
		Summary:     "Update all fields of a teacher",
		Description: "Update all fields of a teacher mandatory.",
		Tags:        []string{"Teachers"},
	}, teacherHandler.UpdateTeacherHandler)

	huma.Register(api, huma.Operation{
		OperationID: "patch-teacher",
		Method:      http.MethodPatch,
		Path:        "/teachers/{id}",
		Summary:     "Patch teacher",
		Description: "Patch some teacher fields only.",
		Tags:        []string{"Teachers"},
	}, teacherHandler.PatchTeacherHandler)

	huma.Register(api, huma.Operation{
		OperationID: "delete-teacher",
		Method:      http.MethodDelete,
		Path:        "/teachers/{id}",
		Summary:     "Delete Teacher by ID",
		Description: "Delete a teacher record by ID.",
		Tags:        []string{"Teachers"},
	}, teacherHandler.DeleteTeacherHandler)

	huma.Register(api, huma.Operation{
		OperationID: "patch-teachers",
		Method:      http.MethodPatch,
		Path:        "/teachers",
		Summary:     "Patch teachers",
		Description: "Patch bulk many teachers fields.",
		Tags:        []string{"Teachers"},
	}, teacherHandler.PatchTeachersHandler)

	huma.Register(api, huma.Operation{
		OperationID: "students-by-teacher-id",
		Method:      http.MethodGet,
		Path:        "/teachers/{id}/students",
		Summary:     "Get students by id of the teacher",
		Description: "Get students by teacher id",
		Tags:        []string{"Teachers"},
	}, teacherHandler.GetStudentsByTeacherId)

	huma.Register(api, huma.Operation{
		OperationID: "delete-teachers",
		Method:      http.MethodDelete,
		Path:        "/teachers",
		Summary:     "Delete teachers",
		Description: "Delete bulk many teachers fields.",
		Tags:        []string{"Teachers"},
	}, teacherHandler.DeleteTeachersHandler)
}

func routesStudents(api huma.API, studentHandler *handlers.StudentHandlers) {
	huma.Register(api, huma.Operation{
		OperationID: "get-student",
		Method:      http.MethodGet,
		Path:        "/students/{id}",
		Summary:     "Get a student",
		Description: "Get a student by ID.",
		Tags:        []string{"Students"},
	}, studentHandler.StudentGet)

	huma.Register(api, huma.Operation{
		OperationID: "post-students",
		Method:      http.MethodPost,
		Path:        "/students",
		Summary:     "Create students",
		Description: "Create students.",
		Tags:        []string{"Students"},
	}, studentHandler.StudentsAdd)

	huma.Register(api, huma.Operation{
		OperationID: "get-students",
		Method:      http.MethodGet,
		Path:        "/students",
		Summary:     "Get all students",
		Description: "Get all students or with filtering.",
		Tags:        []string{"Students"},
	}, studentHandler.StudentsGet)

	huma.Register(api, huma.Operation{
		OperationID: "update-student",
		Method:      http.MethodPut,
		Path:        "/students/{id}",
		Summary:     "Update all fields of a student",
		Description: "Update all fields of a student mandatory.",
		Tags:        []string{"Students"},
	}, studentHandler.UpdateStudentHandler)

	huma.Register(api, huma.Operation{
		OperationID: "patch-student",
		Method:      http.MethodPatch,
		Path:        "/students/{id}",
		Summary:     "Patch student",
		Description: "Patch some student fields only.",
		Tags:        []string{"Students"},
	}, studentHandler.PatchStudentHandler)

	huma.Register(api, huma.Operation{
		OperationID: "delete-student",
		Method:      http.MethodDelete,
		Path:        "/students/{id}",
		Summary:     "Delete Student by ID",
		Description: "Delete a student record by ID.",
		Tags:        []string{"Students"},
	}, studentHandler.DeleteStudentHandler)

	huma.Register(api, huma.Operation{
		OperationID: "patch-students",
		Method:      http.MethodPatch,
		Path:        "/students",
		Summary:     "Patch students",
		Description: "Patch bulk many students fields.",
		Tags:        []string{"Students"},
	}, studentHandler.PatchStudentsHandler)

	huma.Register(api, huma.Operation{
		OperationID: "delete-students",
		Method:      http.MethodDelete,
		Path:        "/students",
		Summary:     "Delete students",
		Description: "Delete bulk many students fields.",
		Tags:        []string{"Students"},
	}, studentHandler.DeleteStudentsHandler)
}

func routesExec(api huma.API, execHandler *handlers.ExecsHandlers) {
	huma.Register(api, huma.Operation{
		OperationID: "get-exec",
		Method:      http.MethodGet,
		Path:        "/exec/{id}",
		Summary:     "Get exec",
		Description: "Get exec.",
		Tags:        []string{"Exec"},
	}, execHandler.ExecGetHandler)

	huma.Register(api, huma.Operation{
		OperationID: "add-exec",
		Method:      http.MethodPost,
		Path:        "/execs",
		Summary:     "Add exec",
		Description: "Add exec.",
		Tags:        []string{"Exec"},
	}, execHandler.ExecAddHandler)

	huma.Register(api, huma.Operation{
		OperationID: "get-execs",
		Method:      http.MethodGet,
		Path:        "/execs",
		Summary:     "Get execs",
		Description: "Get execs.",
		Tags:        []string{"Exec"},
	}, execHandler.ExecsGetHandler)

	huma.Register(api, huma.Operation{
		OperationID: "patch-execs",
		Method:      http.MethodPatch,
		Path:        "/execs",
		Summary:     "Patch execs",
		Description: "Patch bulk many execs.",
		Tags:        []string{"Exec"},
	}, execHandler.PatchExecsHandler)

	huma.Register(api, huma.Operation{
		OperationID: "delete-exec-by-id",
		Method:      http.MethodDelete,
		Path:        "/execs/{id}",
		Summary:     "Delete exec by id",
		Description: "Delete exec by id.",
		Tags:        []string{"Exec"},
	}, execHandler.ExecDeleteByIDHandler)

	huma.Register(api, huma.Operation{
		OperationID: "password-change-exec",
		Method:      http.MethodPost,
		Path:        "/execs/{id}/updatepassword",
		Summary:     "Change exec password by id",
		Description: "Change exec password by id.",
		Tags:        []string{"Exec"},
	}, execHandler.ExecPasswordChangeHandler)

	huma.Register(api, huma.Operation{
		OperationID: "login-exec",
		Method:      http.MethodPost,
		Path:        "/execs/login",
		Summary:     "Logiun exec",
		Description: "Login execs.",
		Tags:        []string{"Exec"},
	}, execHandler.ExecLoginHandler)

	huma.Register(api, huma.Operation{
		OperationID: "logout-execs",
		Method:      http.MethodPost,
		Path:        "/execs/logout",
		Summary:     "Logout exec",
		Description: "Logout execs.",
		Tags:        []string{"Exec"},
	}, execHandler.LogoutExecsHandler)

	huma.Register(api, huma.Operation{
		OperationID: "forgotpassword-execs",
		Method:      http.MethodPost,
		Path:        "/execs/forgotpassword",
		Summary:     "Forgot password exec",
		Description: "Forgot password execs.",
		Tags:        []string{"Exec"},
	}, execHandler.ForgotpasswordExecsHandler)

	huma.Register(api, huma.Operation{
		OperationID: "resetpassword-execs",
		Method:      http.MethodPost,
		Path:        "/execs/resetpassword/reset/{resetcode}",
		Summary:     "Password reset exec",
		Description: "Password reset execs.",
		Tags:        []string{"Exec"},
	}, execHandler.PasswordresetExecsHandler)
}
