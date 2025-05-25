package entity

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/james-mcallister/may/database"
)

type EntityNetwork struct {
	Nets []database.Network
}

func Networks(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		data := EntityNetwork{}
		data.Nets, err = database.AllNetworks(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = t.ExecuteTemplate(w, "entity-network.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
