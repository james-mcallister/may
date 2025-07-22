package plan

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func PlanRow(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		// expecting emp_ids and plan_ids to be a comma separated list of int64 nums
		params := r.URL.Query()

		empIdsParam := params.Get("emp_ids")
		empIdsStr := strings.Split(empIdsParam, ",")
		empIds := make([]int64, len(empIdsStr))
		for i, val := range empIdsStr {
			empIds[i], err = strconv.ParseInt(val, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		planIdsParam := params.Get("plan_ids")
		planIdsStr := strings.Split(planIdsParam, ",")
		planIds := make([]int64, len(planIdsStr))
		for i, val := range planIdsStr {
			planIds[i], err = strconv.ParseInt(val, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		startDate := params.Get("start_date")
		endDate := params.Get("end_date")

		if !ValidateDateFormat(startDate) || !ValidateDateFormat(endDate) {
			http.Error(w, "Invalid query params: start/end date", http.StatusBadRequest)
			return
		}

		planRow, err := database.GetPlanRows(db, empIds, planIds, startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := t.ExecuteTemplate(w, "plan-emp-row.html", planRow); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

		// TODO: need to add param for PlanPage id (int64)
		// this might actually not be needed until save is clicked.
		// then save will associate all the plans on the page with the page
		// this gives flexibility for which tables are used on which page
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

type RowId struct {
	EmpId  int64 `json:"emp_id,string"`
	PlanId int64 `json:"plan_id,string"`
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

		rowId := RowId{
			EmpId:  empId,
			PlanId: planId,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		encoder.Encode(rowId)
	})
}

func NewPlanRowForm(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		data, err := database.NewDropdown(db, database.EmployeeDropdownQuery())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = t.ExecuteTemplate(w, "emp-row-select.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func DeleteRow(db *sql.DB) http.Handler {
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

		_, err = database.DeletePlanRow(db, empId, planId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func UpdateRow(db *sql.DB) http.Handler {
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

		body, err := io.ReadAll(io.LimitReader(r.Body, 10<<20)) // 10MB limit
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var data map[string]float64
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		i := 0
		rows := make([]database.PlanDay, len(data))
		for k, v := range data {
			r := database.PlanDay{
				CalDate:   k,
				PlanHours: v,
			}
			rows[i] = r
			i++
		}

		if err = database.UpdatePlanRow(db, empId, planId, rows); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func Calendar(t *template.Template, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		params := r.URL.Query()
		popStart := params.Get("start_date")
		popEnd := params.Get("end_date")
		fiscalPeriod := params.Get("fiscal_period")

		if !ValidateDateFormat(popStart) || !ValidateDateFormat(popEnd) {
			http.Error(w, "Invalid query params: start/end date", http.StatusBadRequest)
			return
		}

		if len(fiscalPeriod) != 6 {
			http.Error(w, "Invalid query params: start/end date", http.StatusBadRequest)
			return
		}

		data := database.NewCal()
		data.Period = MonthDay(fiscalPeriod)
		data.Days, err = database.GetCalData(db, popStart, popEnd, fiscalPeriod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// init loop counters for template
		numWeeks := len(data.Days) / 7
		numDays := 7
		data.Outer = make([]int, numWeeks)
		data.Inner = make([]int, numDays)

		if err = t.ExecuteTemplate(w, "plan-calendar.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func MonthDay(fiscalPeriod string) string {
	// fiscal period format is YYYYMM
	y, m := fiscalPeriod[:4], fiscalPeriod[4:]
	lookup := map[string]string{
		"01": "January",
		"02": "February",
		"03": "March",
		"04": "April",
		"05": "May",
		"06": "June",
		"07": "July",
		"08": "August",
		"09": "September",
		"10": "October",
		"11": "November",
		"12": "December",
	}
	return lookup[m] + " " + y
}
