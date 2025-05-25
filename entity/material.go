package entity

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/james-mcallister/may/database"
)

type EntityMaterial struct {
	Mats []database.Material
}

func Material(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		data := EntityMaterial{}
		data.Mats, err = database.AllMaterials(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = t.ExecuteTemplate(w, "entity-material.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
