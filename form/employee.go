package form

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/james-mcallister/may/database"
)

type EmployeeForm struct {
	Emp          database.Employee
	EmpDropdown  []database.Dropdown
	IptDropdown  []database.Dropdown
	CompDropdown []database.Dropdown
}

func Employee(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := EmployeeForm{}
		if id == 0 {
			currentTime := time.Now()
			data.Emp = database.Employee{
				Id:            0,
				LaborCapacity: 1.0,
				Active:        true,
				CoverageStart: currentTime.Format("2006-01-02"),
				CoverageEnd:   "2040-12-28",
			}
		} else {
			data.Emp, err = database.GetEmployee(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		data.EmpDropdown, err = database.NewDropdown(db, database.EmployeeDropdownQuery())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.IptDropdown, err = database.NewDropdown(db, database.IptDropdownQuery())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.CompDropdown, err = database.NewDropdown(db, database.CompensationDropdownQuery())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = t.ExecuteTemplate(w, "form-employee.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func NewEmployee(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var active bool
		if r.FormValue("active") == "on" {
			active = true
		}

		laborCap, err := strconv.ParseFloat(r.FormValue("labor_cap"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e := database.Employee{
			FirstName:     r.FormValue("first_name"),
			LastName:      r.FormValue("last_name"),
			Myid:          r.FormValue("myid"),
			DisplayName:   r.FormValue("last_name") + ", " + r.FormValue("first_name") + " (" + r.FormValue("myid") + ")",
			Empid:         r.FormValue("empid"),
			LaborCapacity: laborCap,
			Desk:          r.FormValue("desk"),
			Active:        active,
			CoverageStart: r.FormValue("cov_start"),
			CoverageEnd:   r.FormValue("cov_end"),
		}

		if r.PostForm.Has("comp") {
			v, err := strconv.ParseInt(r.FormValue("comp"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			e.Comp = sql.NullInt64{Int64: v, Valid: true}
		}

		if r.PostForm.Has("manager") {
			v, err := strconv.ParseInt(r.FormValue("manager"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			e.Manager = sql.NullInt64{Int64: v, Valid: true}
		}

		if r.PostForm.Has("ipt") {
			v, err := strconv.ParseInt(r.FormValue("ipt"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			e.Ipt = sql.NullInt64{Int64: v, Valid: true}
		}

		rows, err := database.InsertEmployee(db, e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func UpdateEmployee(db *sql.DB) http.Handler {
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

		var active bool
		if r.FormValue("active") == "on" {
			active = true
		}

		laborCap, err := strconv.ParseFloat(r.FormValue("labor_cap"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e := database.Employee{
			Id:            id,
			FirstName:     r.FormValue("first_name"),
			LastName:      r.FormValue("last_name"),
			Myid:          r.FormValue("myid"),
			DisplayName:   r.FormValue("last_name") + ", " + r.FormValue("first_name") + " (" + r.FormValue("myid") + ")",
			Empid:         r.FormValue("empid"),
			LaborCapacity: laborCap,
			Desk:          r.FormValue("desk"),
			Active:        active,
			CoverageStart: r.FormValue("cov_start"),
			CoverageEnd:   r.FormValue("cov_end"),
		}

		if r.PostForm.Has("comp") {
			v, err := strconv.ParseInt(r.FormValue("comp"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			e.Comp = sql.NullInt64{Int64: v, Valid: true}
		}

		if r.PostForm.Has("manager") {
			v, err := strconv.ParseInt(r.FormValue("manager"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			e.Manager = sql.NullInt64{Int64: v, Valid: true}
		}

		if r.PostForm.Has("ipt") {
			v, err := strconv.ParseInt(r.FormValue("ipt"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			e.Ipt = sql.NullInt64{Int64: v, Valid: true}
		}

		rows, err := database.UpdateEmployee(db, e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func DeleteEmployee(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rows, err := database.DeleteRow(db, "Employee", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Success: %d rows affected.", rows)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}
