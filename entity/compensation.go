package entity

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/james-mcallister/may/database"
)

type EntityCompensation struct {
	Comps []database.Compensation
}

func Compensation(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		data := EntityCompensation{}
		data.Comps, err = database.AllCompensation(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = t.ExecuteTemplate(w, "entity-compensation.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
