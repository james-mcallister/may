package database

import (
	"database/sql"
	"fmt"
)

type Employee struct {
	Id            int64         `json:"id,string"`
	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	DisplayName   string        `json:"display_name"`
	Myid          string        `json:"myid"`
	Empid         string        `json:"empid"`
	LaborCapacity float64       `json:"labor_cap,string"`
	Desk          string        `json:"desk"`
	Active        bool          `json:"active"`
	CoverageStart string        `json:"cov_start"`
	CoverageEnd   string        `json:"cov_end"`
	Comp          sql.NullInt64 `json:"comp"`
	Manager       sql.NullInt64 `json:"manager"`
	Ipt           sql.NullInt64 `json:"ipt"`
}

func NewEmployee() Employee {
	return Employee{}
}

func (e Employee) Init(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS Employee (
		id INTEGER PRIMARY KEY,
		first_name TEXT DEFAULT '',
		last_name TEXT DEFAULT '',
		display_name TEXT DEFAULT '',
		myid TEXT UNIQUE NOT NULL,
		empid TEXT DEFAULT '00000',
		labor_capacity NUMERIC DEFAULT 1.0,
		desk TEXT DEFAULT '',
		active BOOLEAN DEFAULT TRUE,
		coverage_start TEXT DEFAULT CURRENT_DATE,
		coverage_end TEXT DEFAULT '2040-12-28',
		comp INTEGER,
		reports_to INTEGER,
		ipt INTEGER,
		FOREIGN KEY (comp) REFERENCES Compensation(id)
			ON DELETE SET NULL,
		FOREIGN KEY (reports_to) REFERENCES Employee(id)
			ON DELETE SET NULL,
		FOREIGN KEY (ipt) REFERENCES IPT(id)
			ON DELETE SET NULL,
		CHECK (labor_capacity > 0.0 AND labor_capacity <= 1.0)
	);
	`

	insertQuery := `
	INSERT OR IGNORE INTO Employee
	  (id, first_name, last_name, myid, empid, comp, display_name)
	VALUES
	  (1, 'M1', 'TBD', 'x00001', '00001', 1, 'TBD, M1 (x00001)'),
	  (2, 'M2', 'TBD', 'x00002', '00002', 2, 'TBD, M2 (x00002)'),
	  (3, 'M3', 'TBD', 'x00003', '00003', 3, 'TBD, M3 (x00003)'),
	  (4, 'T1', 'TBD', 'x00004', '00004', 4, 'TBD, T1 (x00004)'),
	  (5, 'T2', 'TBD', 'x00005', '00005', 5, 'TBD, T2 (x00005)'),
	  (6, 'T3', 'TBD', 'x00006', '00006', 6, 'TBD, T3 (x00006)'),
	  (7, 'T4', 'TBD', 'x00007', '00007', 7, 'TBD, T4 (x00007)'),
	  (8, 'T5', 'TBD', 'x00008', '00008', 8, 'TBD, T5 (x00008)'),
	  (9, 'T6', 'TBD', 'x00009', '00009', 9, 'TBD, T6 (x00009)'),
	  (10, 'T7', 'TBD', 'x00016', '00016', 10, 'TBD, T7 (x00010)'),
	  (11, 'A1', 'TBD', 'x00010', '00010', 11, 'TBD, A1 (x00011)'),
	  (12, 'A2', 'TBD', 'x00011', '00011', 12, 'TBD, A2 (x00012)'),
	  (13, 'A3', 'TBD', 'x00012', '00012', 13, 'TBD, A3 (x00013)'),
	  (14, 'A4', 'TBD', 'x00013', '00013', 14, 'TBD, A4 (x00014)'),
	  (15, 'A5', 'TBD', 'x00014', '00014', 15, 'TBD, A5 (x00015)'),
	  (16, 'Intern', 'TBD', 'x00015', '00015', 16, 'TBD, Intern (x00015)');
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

	_, err = tx.Exec(insertQuery)
	if err != nil {
		return fmt.Errorf("error executing INSERT transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func EmployeeDropdownQuery() string {
	return "SELECT display_name AS name,id FROM Employee;"
}

func EmployeeImportColumns() []string {
	return []string{"first_name", "last_name", "display_name", "myid", "empid",
		"ipt", "manager", "grade"}
}

func GetEmployee(db *sql.DB, id int64) (Employee, error) {
	var emp Employee

	getQuery := `
	SELECT
	  id,first_name,last_name,myid,empid,labor_capacity,
	  desk,active,coverage_start,coverage_end,comp,reports_to,ipt
	FROM Employee
	WHERE id=?;
	`

	row := db.QueryRow(getQuery, id)
	if err := row.Scan(&emp.Id, &emp.FirstName, &emp.LastName, &emp.Myid, &emp.Empid, &emp.LaborCapacity, &emp.Desk, &emp.Active, &emp.CoverageStart, &emp.CoverageEnd, &emp.Comp, &emp.Manager, &emp.Ipt); err != nil {
		if err == sql.ErrNoRows {
			return emp, fmt.Errorf("employee id=%d: no such row", id)
		}
		return emp, fmt.Errorf("employee: id=%d: %v", id, err)
	}
	return emp, nil
}

func AllEmployees(db *sql.DB) ([]Employee, error) {
	var emps []Employee

	getQuery := `
	SELECT
	  id,first_name,last_name,myid,empid,labor_capacity,
	  desk,active,coverage_start,coverage_end,comp,reports_to,ipt
	FROM Employee;
	`

	rows, err := db.Query(getQuery)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.Id, &emp.FirstName, &emp.LastName, &emp.Myid, &emp.Empid, &emp.LaborCapacity, &emp.Desk, &emp.Active, &emp.CoverageStart, &emp.CoverageEnd, &emp.Comp, &emp.Manager, &emp.Ipt); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		emps = append(emps, emp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return emps, nil
}

func UpdateEmployee(db *sql.DB, emp Employee) (int64, error) {
	updateQuery := `
	UPDATE Employee SET
	  first_name=?,last_name=?,myid=?,empid=?,labor_capacity=?,desk=?,
	  active=?,coverage_start=?,coverage_end=?,comp=?,reports_to=?,ipt=?
	WHERE id=?;
	`

	result, err := db.Exec(updateQuery, emp.FirstName, emp.LastName, emp.Myid, emp.Empid, emp.LaborCapacity, emp.Desk, emp.Active, emp.CoverageStart, emp.CoverageEnd, emp.Comp, emp.Manager, emp.Ipt, emp.Id)
	if err != nil {
		return 0, fmt.Errorf("update query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update result error: %v", err)
	}
	return rows, nil
}

func InsertEmployee(db *sql.DB, emp Employee) (int64, error) {
	insertQuery := `
	INSERT INTO Employee
	  (first_name,last_name,myid,empid,labor_capacity,
	  desk,active,coverage_start,coverage_end,comp,reports_to,ipt)
	VALUES
	  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	result, err := db.Exec(insertQuery, emp.FirstName, emp.LastName, emp.Myid, emp.Empid, emp.LaborCapacity, emp.Desk, emp.Active, emp.CoverageStart, emp.CoverageEnd, emp.Comp, emp.Manager, emp.Ipt)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return rows, nil
}
