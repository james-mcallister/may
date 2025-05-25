package form

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/james-mcallister/may/database"
)

type MaterialForm struct {
	Mat          database.Material
	ProjDropdown []database.Dropdown
}

func Material(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := MaterialForm{}
		if id == 0 {
			data.Mat = database.Material{}
		} else {
			data.Mat, err = database.GetMaterial(db, id)
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

		if err = t.ExecuteTemplate(w, "form-material.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func NewMaterial(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var done bool
		if r.FormValue("status") == "complete" {
			done = true
		}

		m := database.Material{
			Name:                r.FormValue("name"),
			PRDate:              r.FormValue("pr_date"),
			PODate:              r.FormValue("po_date"),
			PRNumber:            r.FormValue("pr"),
			PONumber:            r.FormValue("po"),
			Complete:            done,
			BaselineStartDate:   r.FormValue("baseline_start_date"),
			BaselineFinishDate:  r.FormValue("baseline_finish_date"),
			TentativeStartDate:  r.FormValue("tentative_start_date"),
			TentativeFinishDate: r.FormValue("tentative_finish_date"),
			ActualStartDate:     r.FormValue("actual_start_date"),
			ActualFinishDate:    r.FormValue("actual_finish_date"),
			Notes:               r.FormValue("notes"),
		}

		if r.PostForm.Has("estimated_cost") && len(r.FormValue("estimated_cost")) > 0 {
			estCost, err := strconv.ParseFloat(r.FormValue("estimated_cost"), 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m.EstimatedCost = estCost
		}

		if r.PostForm.Has("actual_cost") && len(r.FormValue("actual_cost")) > 0 {
			actCost, err := strconv.ParseFloat(r.FormValue("actual_cost"), 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m.ActualCost = actCost
		}

		if r.PostForm.Has("wp") {
			v, err := strconv.ParseInt(r.FormValue("wp"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m.WorkPackage = sql.NullInt64{Int64: v, Valid: true}
		}

		rows, err := database.InsertMaterial(db, m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func UpdateMaterial(db *sql.DB) http.Handler {
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

		var done bool
		if r.FormValue("status") == "complete" {
			done = true
		}

		m := database.Material{
			Id:                  id,
			Name:                r.FormValue("name"),
			PRDate:              r.FormValue("pr_date"),
			PODate:              r.FormValue("po_date"),
			PRNumber:            r.FormValue("pr"),
			PONumber:            r.FormValue("po"),
			Complete:            done,
			BaselineStartDate:   r.FormValue("baseline_start_date"),
			BaselineFinishDate:  r.FormValue("baseline_finish_date"),
			TentativeStartDate:  r.FormValue("tentative_start_date"),
			TentativeFinishDate: r.FormValue("tentative_finish_date"),
			ActualStartDate:     r.FormValue("actual_start_date"),
			ActualFinishDate:    r.FormValue("actual_finish_date"),
			Notes:               r.FormValue("notes"),
		}

		if r.PostForm.Has("estimated_cost") && len(r.FormValue("estimated_cost")) > 0 {
			estCost, err := strconv.ParseFloat(r.FormValue("estimated_cost"), 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m.EstimatedCost = estCost
		}

		if r.PostForm.Has("actual_cost") && len(r.FormValue("actual_cost")) > 0 {
			actCost, err := strconv.ParseFloat(r.FormValue("actual_cost"), 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m.ActualCost = actCost
		}

		if r.PostForm.Has("wp") {
			v, err := strconv.ParseInt(r.FormValue("wp"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m.WorkPackage = sql.NullInt64{Int64: v, Valid: true}
		}

		rows, err := database.UpdateMaterial(db, m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func DeleteMaterial(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rows, err := database.DeleteRow(db, "Material", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}
