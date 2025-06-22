package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type PlanDay struct {
	PlanHours   float64 `json:"plan_hours"`
	CalDate     string  `json:"cal_date"`
	Description string  `json:"description"`
}

func NewPlanDay() PlanDay {
	return PlanDay{}
}

func (p PlanDay) Init(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS PlanDay (
		id INTEGER PRIMARY KEY,
		planned_hours NUMERIC DEFAULT 0.0,
		cal_date TEXT NOT NULL,
		description TEXT DEFAULT '',
		updated_at TEXT DEFAULT CURRENT_DATE,
		emp INTEGER,
		plan INTEGER,
		FOREIGN KEY (emp) REFERENCES Employee(id)
			ON DELETE CASCADE,
		FOREIGN KEY (plan) REFERENCES Plan(id)
			ON DELETE CASCADE,
		UNIQUE (cal_date, emp, plan) ON CONFLICT ABORT
	);
	`

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(tableQuery)
	if err != nil {
		return fmt.Errorf("error executing CREATE TABLE transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func InitPlanRow(db *sql.DB, empId, planId int64, startDate, endDate string) (int64, error) {
	dates, err := GetDateList(db, startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("date list error: %v", err)
	}

	var sb strings.Builder
	e := strconv.FormatInt(empId, 10)
	p := strconv.FormatInt(planId, 10)
	l := len(dates) - 1

	sb.WriteString("INSERT INTO PlanDay (planned_hours,cal_date,emp,plan) VALUES ")
	for i := 0; i < l; i++ {
		sb.WriteString("(0.0,")
		sb.WriteString("'")
		sb.WriteString(dates[i])
		sb.WriteString("'")
		sb.WriteString(",")
		sb.WriteString(e)
		sb.WriteString(",")
		sb.WriteString(p)
		sb.WriteString("),")
	}
	// last row is semicolon terminated
	sb.WriteString("(0.0,")
	sb.WriteString("'")
	sb.WriteString(dates[l])
	sb.WriteString("'")
	sb.WriteString(",")
	sb.WriteString(e)
	sb.WriteString(",")
	sb.WriteString(p)
	sb.WriteString(");")

	insertQuery := sb.String()

	result, err := db.Exec(insertQuery)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return rows, nil
}

// generates a list of ISO formatted dates from start to end inclusive
func GetDateList(db *sql.DB, startDate, endDate string) ([]string, error) {
	dates := make([]string, 0)

	q := `
	SELECT cal_date
	FROM CalendarHours
	WHERE cal_date BETWEEN ? AND ?;
	`

	rows, err := db.Query(q, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("no such row: %v", err)
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		dates = append(dates, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database rows error: %v", err)
	}
	return dates, nil
}

type PlanRow struct {
	EmpId  int64              `json:"emp_id"`
	PlanId int64              `json:"plan_id"`
	Hours  map[string]float64 `json:"hours"`
}

func GetPlanRow(db *sql.DB, empId, planId int64, startDate, endDate string) (PlanRow, error) {
	hours := make(map[string]float64)
	row := PlanRow{
		EmpId:  empId,
		PlanId: planId,
	}

	getQuery := `
	SELECT cal_date,planned_hours
	FROM PlanDay
	WHERE cal_date BETWEEN ? AND ?
	  AND emp = ?
	  AND plan = ?;
	`

	rows, err := db.Query(getQuery, startDate, endDate, empId, planId)
	if err != nil {
		return row, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d string
		var p float64
		if err := rows.Scan(&d, &p); err != nil {
			if err == sql.ErrNoRows {
				return row, fmt.Errorf("error: no rows")
			}
			return row, fmt.Errorf("row scan error: %v", err)
		}
		hours[d] = p
	}
	if err := rows.Err(); err != nil {
		return row, fmt.Errorf("rows error: %v", err)
	}
	row.Hours = hours
	return row, nil
}

func GetPlanHours(db *sql.DB, empId, planId int64, startDate, endDate string) ([]float64, error) {
	var h []float64

	getQuery := `
	SELECT planned_hours
	FROM PlanDay
	WHERE cal_date BETWEEN ? AND ?
	  AND emp = ?
	  AND plan = ?;
	`

	rows, err := db.Query(getQuery, startDate, endDate, empId, planId)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p float64
		if err := rows.Scan(&p); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		h = append(h, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return h, nil
}

type TableRow struct {
	EmpId     int64
	ScopeId   int64
	EmpName   string
	ScopeName string
	LaborRate string
}

// func LoadPlanRows(db *sql.DB, empId, planId int64, startDate, endDate string) ([]TableRow, error) {
// 	var planRows []TableRow

// 	// TODO: plan hours by day is a separate query. load here should be able to use
// 	// sum of each month with group by for the template
// 	// try to have zeros for cal_date not in PlanDay table
// 	getQuery := `
// 	SELECT c.fiscal_year,c.fiscal_month,sum(p.planned_hours)
// 	FROM CalendarHours c
// 	JOIN PlanDay p ON c.cal_date = p.cal_date
// 	WHERE c.cal_date BETWEEN ? AND ?
// 	  AND emp = ?
// 	  AND plan = ?
// 	GROUP BY c.fiscal_year,c.fiscal_month;
// 	`

// 	rows, err := db.Query(getQuery, startDate, endDate, empId, planId)
// 	if err != nil {
// 		return planRows, fmt.Errorf("query error: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var r Row
// 		if err := rows.Scan(); err != nil {
// 			if err == sql.ErrNoRows {
// 				return planRows, fmt.Errorf("error: no rows")
// 			}
// 			return planRows, fmt.Errorf("row scan error: %v", err)
// 		}
// 		hours[d] = p
// 	}
// 	if err := rows.Err(); err != nil {
// 		return planRows, fmt.Errorf("rows error: %v", err)
// 	}
// 	return planRows, nil
// }

func UpdatePlanRow(db *sql.DB, empId, planId int64, rows []PlanDay) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("database transaction error: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE PlanDay SET updated_at=CURRENT_DATE, planned_hours=?, description=? WHERE cal_date=? AND emp=? AND plan=?;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range rows {
		if _, err := stmt.Exec(v.PlanHours, v.Description, v.CalDate, empId, planId); err != nil {
			return fmt.Errorf("stmt exec error: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %v", err)
	}
	return nil
}

type PlanMonth struct {
	FiscalPeriod string
	StartDate    string
	EndDate      string
	DisplayName  string
}

// the column name format is MMM-YYYY (ex: Oct-2024)
func GetPlanMonths(db *sql.DB, startDate, endDate string) ([]PlanMonth, error) {
	var months []PlanMonth

	q := `
	SELECT fiscal_year,fiscal_month,fiscal_period,min(cal_date),max(cal_date)
	FROM CalendarHours
	WHERE fiscal_period IN (
	    SELECT DISTINCT fiscal_period
		FROM CalendarHours
		WHERE cal_date BETWEEN ? AND ?
	)
	GROUP BY fiscal_period;
	`

	rows, err := db.Query(q, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fy, fm int
		var month PlanMonth
		if err := rows.Scan(&fy, &fm, &month.FiscalPeriod, &month.StartDate, &month.EndDate); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no such row: %v", err)
			}
			return nil, fmt.Errorf("rows scan error: %v", err)
		}
		formatted := FormatDate(fy, fm)
		month.DisplayName = formatted
		months = append(months, month)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	return months, nil
}

// convert int values fiscal_month=2, fiscal_year=2024 to
// string formatted column for planner table/lookups Mar-2024
func FormatDate(fy, fm int) string {
	monthLookup := []string{
		"",
		"Jan",
		"Feb",
		"Mar",
		"Apr",
		"May",
		"Jun",
		"Jul",
		"Aug",
		"Sep",
		"Oct",
		"Nov",
		"Dec",
	}
	return monthLookup[fm] + "-" + strconv.Itoa(fy)
}
