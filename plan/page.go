package plan

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/james-mcallister/may/database"
)

type LoadPlan struct {
	Plans []database.Dropdown
}

func Select(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		data := LoadPlan{}
		data.Plans, err = database.NewDropdown(db, database.PlanPageDropdownQuery())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := t.ExecuteTemplate(w, "form-plan-select.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

type NewPlanData struct {
	Emps []database.Dropdown
}

func New(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		data := NewPlanData{}
		data.Emps, err = database.NewDropdown(db, database.EmployeeDropdownQuery())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := t.ExecuteTemplate(w, "form-plan-new.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func Page(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data database.PlanPage
		if r.PostForm.Has("load_plan") && len(r.FormValue("load_plan")) > 0 {
			planPageId, err := strconv.ParseInt(r.FormValue("load_plan"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data, err = database.GetPlanPage(db, planPageId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			data = database.PlanPage{
				Title:       r.FormValue("name"),
				Description: r.FormValue("description"),
			}

			_, err = database.InsertPlanPage(db, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err = t.ExecuteTemplate(w, "plan-page.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func PlanRow(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		params := r.URL.Query()
		empId, err := strconv.ParseInt(params.Get("emp_id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		planId, err := strconv.ParseInt(params.Get("plan_id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		startDate := params.Get("start_date")
		endDate := params.Get("end_date")

		if !ValidateDateFormat(startDate) || !ValidateDateFormat(endDate) {
			http.Error(w, "Invalid query params: start/end date", http.StatusBadRequest)
			return
		}

		planHours, err := database.GetPlanRow(db, empId, planId, startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		encoder.Encode(planHours)
	})
}

func PlanHours(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		params := r.URL.Query()
		empId, err := strconv.ParseInt(params.Get("emp_id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		planId, err := strconv.ParseInt(params.Get("plan_id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		startDate := params.Get("start_date")
		endDate := params.Get("end_date")

		if !ValidateDateFormat(startDate) || !ValidateDateFormat(endDate) {
			http.Error(w, "Invalid query params: start/end date", http.StatusBadRequest)
			return
		}

		planHours, err := database.GetPlanHours(db, empId, planId, startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		encoder.Encode(planHours)
	})
}

func ProdHours(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		params := r.URL.Query()
		startDate := params.Get("start_date")
		endDate := params.Get("end_date")

		if !ValidateDateFormat(startDate) || !ValidateDateFormat(endDate) {
			http.Error(w, "Invalid query params: start/end date", http.StatusBadRequest)
			return
		}

		prodHours, err := database.GetProdHours(db, startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		encoder.Encode(prodHours)
	})
}

func ProdHoursIdx(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		params := r.URL.Query()
		startDate := params.Get("start_date")
		endDate := params.Get("end_date")

		if !ValidateDateFormat(startDate) || !ValidateDateFormat(endDate) {
			http.Error(w, "Invalid query params: start/end date", http.StatusBadRequest)
			return
		}

		prodHours, err := database.GetProdHoursIndex(db, startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		encoder.Encode(prodHours)
	})
}

func ValidateDateFormat(dateString string) bool {
	format := "2006-01-02"

	_, err := time.Parse(format, dateString)
	return err == nil
}

type PlanTable struct {
	Plan    database.Plan        `json:"plan"`
	Months  []database.PlanMonth `json:"months"`
	EmpRows []database.TableRow  `json:"emp_rows"`
}

func NewPlanTable(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !ValidateDateFormat(r.FormValue("start_date")) || !ValidateDateFormat(r.FormValue("end_date")) {
			http.Error(w, "Invalid query params: start/end date", http.StatusBadRequest)
			return
		}

		if r.PostForm.Has("name") && len(r.FormValue("name")) <= 0 {
			http.Error(w, "Invalid query params: name", http.StatusBadRequest)
			return
		}

		tab := database.Plan{
			Name:      r.FormValue("name"),
			StartDate: r.FormValue("start_date"),
			EndDate:   r.FormValue("end_date"),
		}

		tId, err := database.InsertPlan(db, tab)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tab.Id = tId

		data := PlanTable{
			Plan: tab,
		}

		m, err := database.GetPlanMonths(db, tab.StartDate, tab.EndDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.Months = m

		if err = t.ExecuteTemplate(w, "plan-table.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func NewPlanForm(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "form-plan-new-table.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func NewPlanRow(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		params := r.URL.Query()
		empId, err := strconv.ParseInt(params.Get("emp_id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		planId, err := strconv.ParseInt(params.Get("plan_id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		startDate := params.Get("start_date")
		endDate := params.Get("end_date")

		if !ValidateDateFormat(startDate) || !ValidateDateFormat(endDate) {
			http.Error(w, "Invalid query params: start/end date", http.StatusBadRequest)
			return
		}

		_, err = database.InitPlanRow(db, empId, planId, startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
