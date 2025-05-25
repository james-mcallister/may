package entity

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/james-mcallister/may/database"
)

// TODO: add logic to replace the manager, ipt, and comp IDs with the names
type EntityIpt struct {
	Ipts []database.Ipt
}

func Ipts(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		data := EntityIpt{}
		data.Ipts, err = database.AllIpts(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = t.ExecuteTemplate(w, "entity-ipt.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
