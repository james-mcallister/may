package entity

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/james-mcallister/may/database"
)

// TODO: add logic to replace the manager, ipt, and comp IDs with the names
type EntityEmployee struct {
	Emps []database.Employee
}

func Employees(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		data := EntityEmployee{}
		data.Emps, err = database.AllEmployees(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = t.ExecuteTemplate(w, "entity-employee.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
