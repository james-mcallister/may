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

	insertStmt := "INSERT INTO PlanDay (planned_hours,cal_date,emp,plan) VALUES "

	var sb strings.Builder
	e := strconv.FormatInt(empId, 10)
	p := strconv.FormatInt(planId, 10)
	l := len(dates) - 1

	// number of bytes of hardcoded values in the loop (date is standard length)
	// not including the emp/plan ids
	// (0.0,'2025-01-01',,),
	lenRow := 21
	numRows := len(dates)
	stringBytes := ((len(e) + len(p) + lenRow) * numRows) + len(insertStmt)

	sb.Grow(stringBytes)
	sb.WriteString(insertStmt)
	for i := 0; i < l; i++ {
		sb.WriteString("(0.0,'")
		sb.WriteString(dates[i])
		sb.WriteString("',")
		sb.WriteString(e)
		sb.WriteByte(',')
		sb.WriteString(p)
		sb.WriteString("),")
	}
	// last row is semicolon terminated
	sb.WriteString("(0.0,'")
	sb.WriteString(dates[l])
	sb.WriteString("',")
	sb.WriteString(e)
	sb.WriteByte(',')
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

type TableRow struct {
	EmpId     int64
	ScopeId   int64
	EmpName   string
	ScopeName string
	LaborRate string
	Months    []PlanMonth
}

func GetPlanRows(db *sql.DB, empId, planId []int64, startDate, endDate string) ([]TableRow, error) {
	var data []TableRow
	var sb strings.Builder

	stmt1 := `
	SELECT e.id,e.display_name,
		   p.id,p.name,c.hourly_rate
	FROM PlanDay pd
	JOIN Employee e ON pd.emp=e.id
	JOIN Compensation c ON e.comp=c.id
	JOIN Plan p ON pd.plan=p.id
	WHERE e.id IN (`

	stmt2 := `)
	AND p.id IN (`

	stmt3 := `)
	GROUP BY e.display_name,p.name,c.hourly_rate
	ORDER BY e.display_name;
	`

	// minimize allocations. 50 is a rough estimate for id string size
	growSize := len(stmt1) + len(stmt2) + len(stmt3) + 50
	sb.Grow(growSize)

	sb.WriteString(stmt1)
	num := len(empId) - 1
	for i := 0; i < num; i++ {
		sb.WriteString(strconv.FormatInt(empId[i], 10))
		sb.WriteByte(',')
	}
	sb.WriteString(strconv.FormatInt(empId[num], 10))

	sb.WriteString(stmt2)
	num = len(planId) - 1
	for i := 0; i < num; i++ {
		sb.WriteString(strconv.FormatInt(planId[i], 10))
		sb.WriteByte(',')
	}
	sb.WriteString(strconv.FormatInt(planId[num], 10))
	sb.WriteString(stmt3)

	getQuery := sb.String()

	rows, err := db.Query(getQuery)
	if err != nil {
		return data, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t TableRow
		if err := rows.Scan(&t.EmpId, &t.EmpName, &t.ScopeId, &t.ScopeName, &t.LaborRate); err != nil {
			if err == sql.ErrNoRows {
				return data, fmt.Errorf("error: no rows")
			}
			return data, fmt.Errorf("row scan error: %v", err)
		}

		months, err := GetPlanMonths(db, startDate, endDate)
		if err != nil {
			return data, fmt.Errorf("error getting months: %v", err)
		}
		t.Months = months
		data = append(data, t)
	}
	if err := rows.Err(); err != nil {
		return data, fmt.Errorf("rows error: %v", err)
	}

	return data, nil
}

func GetPlanHours(db *sql.DB, empId, planId int64, startDate, endDate string) ([]float64, error) {
	var h []float64

	getQuery := `
	SELECT IFNULL(planned_hours,0)
	FROM CalendarHours c
	FULL JOIN PlanDay p ON p.cal_date = c.cal_date
		AND emp = ?
		AND plan = ?
	WHERE c.cal_date BETWEEN ? AND ?;
	`

	rows, err := db.Query(getQuery, empId, planId, startDate, endDate)
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
	MonthHours   float64
}

// the column name format is MMM-YYYY (ex: Oct-2024)
func GetPlanMonths(db *sql.DB, startDate, endDate string) ([]PlanMonth, error) {
	var months []PlanMonth

	q := `
	SELECT fiscal_year,fiscal_month,fiscal_period,min(cal_date),max(cal_date),sum(productive_hours)
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
		if err := rows.Scan(&fy, &fm, &month.FiscalPeriod, &month.StartDate, &month.EndDate, &month.MonthHours); err != nil {
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

func DeletePlanRow(db *sql.DB, empId, planId int64) (int64, error) {
	deleteQuery := `
	DELETE FROM PlanDay
	WHERE emp=?
	  AND plan=?;
	`

	result, err := db.Exec(deleteQuery, empId, planId)
	if err != nil {
		return 0, fmt.Errorf("delete query exec error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("delete query result error: %v", err)
	}
	return rows, nil
}
