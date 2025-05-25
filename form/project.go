package form

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/james-mcallister/may/database"
)

type ProjectForm struct {
	Proj         database.Project
	ProjDropdown []database.Dropdown
}

func Project(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := ProjectForm{}
		if id == 0 {
			data.Proj = database.Project{}
		} else {
			data.Proj, err = database.GetProject(db, id)
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

		if err = t.ExecuteTemplate(w, "form-project.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func NewProject(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := database.Project{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			WbsId:       r.FormValue("wbsid"),
			StmtOfWork:  r.FormValue("stmt_of_work"),
			StartDate:   r.FormValue("start_date"),
			EndDate:     r.FormValue("end_date"),
			Evt:         r.FormValue("evt"),
		}

		if r.PostForm.Has("ims_uid") && len(r.FormValue("ims_uid")) > 0 {
			imsUid, err := strconv.Atoi(r.FormValue("ims_uid"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			p.ImsUid = imsUid
		}

		if r.PostForm.Has("wad_lineid") && len(r.FormValue("wad_lineid")) > 0 {
			wadLineId, err := strconv.Atoi(r.FormValue("wad_lineid"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			p.WadLineId = wadLineId
		}

		if r.PostForm.Has("parent_proj") {
			v, err := strconv.ParseInt(r.FormValue("parent_proj"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			p.ParentProject = sql.NullInt64{Int64: v, Valid: true}
		}

		rows, err := database.InsertProject(db, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func UpdateProject(db *sql.DB) http.Handler {
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

		p := database.Project{
			Id:          id,
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			WbsId:       r.FormValue("wbsid"),
			StmtOfWork:  r.FormValue("stmt_of_work"),
			StartDate:   r.FormValue("start_date"),
			EndDate:     r.FormValue("end_date"),
			Evt:         r.FormValue("evt"),
		}

		if r.PostForm.Has("ims_uid") && len(r.FormValue("ims_uid")) > 0 {
			imsUid, err := strconv.Atoi(r.FormValue("ims_uid"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			p.ImsUid = imsUid
		}

		if r.PostForm.Has("wad_lineid") && len(r.FormValue("wad_lineid")) > 0 {
			wadLineId, err := strconv.Atoi(r.FormValue("wad_lineid"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			p.WadLineId = wadLineId
		}

		if r.PostForm.Has("parent_proj") {
			v, err := strconv.ParseInt(r.FormValue("parent_proj"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			p.ParentProject = sql.NullInt64{Int64: v, Valid: true}
		}

		rows, err := database.UpdateProject(db, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func DeleteProject(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rows, err := database.DeleteRow(db, "Project", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}
