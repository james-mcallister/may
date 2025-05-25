package form

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/james-mcallister/may/database"
)

type NetworkForm struct {
	Net          database.Network
	ProjDropdown []database.Dropdown
}

func Network(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := NetworkForm{}
		if id == 0 {
			data.Net = database.Network{}
		} else {
			data.Net, err = database.GetNetwork(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		data.ProjDropdown, err = database.NewDropdown(db, database.ProjectDropdownQuery())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = t.ExecuteTemplate(w, "form-network.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func NewNetwork(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		n := database.Network{
			ChargeNumber: r.FormValue("charge_num"),
			Title:        r.FormValue("title"),
			Description:  r.FormValue("description"),
			Status:       r.FormValue("status"),
			StartDate:    r.FormValue("start_date"),
			EndDate:      r.FormValue("end_date"),
		}

		if r.PostForm.Has("proj") {
			v, err := strconv.ParseInt(r.FormValue("proj"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			n.Proj = sql.NullInt64{Int64: v, Valid: true}
		}

		rows, err := database.InsertNetwork(db, n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func UpdateNetwork(db *sql.DB) http.Handler {
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

		n := database.Network{
			Id:           id,
			ChargeNumber: r.FormValue("charge_num"),
			Title:        r.FormValue("title"),
			Description:  r.FormValue("description"),
			Status:       r.FormValue("status"),
			StartDate:    r.FormValue("start_date"),
			EndDate:      r.FormValue("end_date"),
		}

		if r.PostForm.Has("proj") {
			v, err := strconv.ParseInt(r.FormValue("proj"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			n.Proj = sql.NullInt64{Int64: v, Valid: true}
		}

		rows, err := database.UpdateNetwork(db, n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func DeleteNetwork(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rows, err := database.DeleteRow(db, "Network", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}
