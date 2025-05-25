package form

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/james-mcallister/may/database"
)

type IptForm struct {
	Ipt database.Ipt
}

func Ipt(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := IptForm{}
		if id == 0 {
			data.Ipt = database.Ipt{}
		} else {
			data.Ipt, err = database.GetIpt(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err = t.ExecuteTemplate(w, "form-ipt.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func NewIpt(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		i := database.Ipt{
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
		}

		rows, err := database.InsertIpt(db, i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func UpdateIpt(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		i := database.Ipt{
			Id:          id,
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
		}

		rows, err := database.UpdateIpt(db, i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func DeleteIpt(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rows, err := database.DeleteRow(db, "Ipt", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}
