// Package router - all huma register routes
package router

import (
	"database/sql"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/dataops"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/handlers"
)

func Router(db *sql.DB) *http.ServeMux {
	router := http.NewServeMux()
	teachersDB := dataops.NewTeachersDB(db)
	teacherHandler := handlers.NewTeachersHandler(teachersDB)

	api := humago.New(router, huma.DefaultConfig("My API", "1.0.0"))

	huma.Register(api, huma.Operation{
		OperationID: "get-teacher",
		Method:      http.MethodGet,
		Path:        "/teachers/{id}",
		Summary:     "Get a teacher",
		Description: "Get a teacher by ID.",
	}, teacherHandler.TeacherGet)

	huma.Get(api, "/", teacherHandler.RootHandler)

	huma.Register(api, huma.Operation{
		OperationID: "post-teachers",
		Method:      http.MethodPost,
		Path:        "/teachers",
		Summary:     "Create teachers",
		Description: "Create teachers.",
	}, teacherHandler.TeachersAdd)

	huma.Register(api, huma.Operation{
		OperationID: "get-teachers",
		Method:      http.MethodGet,
		Path:        "/teachers",
		Summary:     "Get all teachers",
		Description: "Get all teachers or with filtering.",
	}, teacherHandler.TeachersGet)

	huma.Register(api, huma.Operation{
		OperationID: "update-teacher",
		Method:      http.MethodPut,
		Path:        "/teachers/{id}",
		Summary:     "Update all fields of a teacher",
		Description: "Update all fields of a teacher mandatory.",
	}, teacherHandler.UpdateTeacherHandler)

	huma.Register(api, huma.Operation{
		OperationID: "patch-teacher",
		Method:      http.MethodPatch,
		Path:        "/teachers/{id}",
		Summary:     "Patch teacher",
		Description: "Patch some teacher fields only.",
	}, teacherHandler.PatchTeacherHandler)

	huma.Register(api, huma.Operation{
		OperationID: "delete-teacher",
		Method:      http.MethodDelete,
		Path:        "/teachers/{id}",
		Summary:     "Delete Teacher by ID",
		Description: "Delete a teacher record by ID.",
	}, teacherHandler.DeleteTeacherHandler)

	huma.Register(api, huma.Operation{
		OperationID: "patch-teachers",
		Method:      http.MethodPatch,
		Path:        "/teachers",
		Summary:     "Patch teachers",
		Description: "Patch bulk many teachers fields.",
	}, teacherHandler.PatchTeachersHandler)

	huma.Register(api, huma.Operation{
		OperationID: "delete-teachers",
		Method:      http.MethodDelete,
		Path:        "/teachers",
		Summary:     "Delete teachers",
		Description: "Delete bulk many teachers fields.",
	}, teacherHandler.DeleteTeachersHandler)
	return router
}
