package database

import (
	"database/sql"
	"fmt"
)

type Compensation struct {
	Id            int64   `json:"id,string"`
	ResourceCode  string  `json:"resource_code"`
	Grade         string  `json:"grade"`
	LaborCategory string  `json:"labor_category"`
	HourlyRate    float64 `json:"hourly_rate,string"`
}

func NewCompensation() Compensation {
	return Compensation{}
}

func (c Compensation) Init(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS Compensation (
		id INTEGER PRIMARY KEY,
		resource_code TEXT UNIQUE NOT NULL,
		grade TEXT NOT NULL,
		labor_category TEXT DEFAULT '',
		hourly_rate NUMERIC DEFAULT 0
	);
	`
	insertQuery := `
	INSERT OR IGNORE INTO Compensation
	  (id, resource_code, grade, labor_category, hourly_rate)
	VALUES
	  (1, 'L_VS1A_X_H', 'MGR01', 'Manager 1', '171'),
	  (2, 'L_VS2A_X_H', 'MGR02', 'Manager 2', '211'),
	  (3, 'L_VS3A_X_H', 'MGR03', 'Manager 3', '245'),
	  (4, 'L_VT1A_X_H', 'ENG01', 'Associate Engineer', '108'),
	  (5, 'L_VT2A_X_H', 'ENG02', 'Engineer', '131'),
	  (6, 'L_VT3A_X_H', 'ENG03', 'Principal Engineer', '155'),
	  (7, 'L_VT4A_X_H', 'ENG04', 'Sr. Principal Engineer', '198'),
	  (8, 'L_VT5A_X_H', 'ENG05', 'Staff Engineer', '242'),
	  (9, 'L_VT6A_X_H', 'ENG06', 'Sr. Staff Engineer', '270'),
	  (10, 'L_VT7A_X_H', 'ENG07', 'Consultant', '298'),
	  (11, 'L_VA1A_X_H', 'ADM01', 'Associate Administrator', '88'),
	  (12, 'L_VA2A_X_H', 'ADM02', 'Administrator', '105'),
	  (13, 'L_VA3A_X_H', 'ADM03', 'Principal Administrator', '127'),
	  (14, 'L_VA4A_X_H', 'ADM04', 'Sr. Principal Administrator', '160'),
	  (15, 'L_VA5A_X_H', 'ADM05', 'Staff Administrator', '202'),
	  (16, 'L_VSHA_X_H', 'TEC01', 'College Intern Technical', '89');
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

func CompensationDropdownQuery() string {
	return "SELECT grade || ' - ' || labor_category AS name,id FROM Compensation ORDER BY name;"
}

func GetCompensation(db *sql.DB, id int64) (Compensation, error) {
	var comp Compensation

	getQuery := `
	SELECT 
	  id,resource_code,grade,labor_category,hourly_rate 
	FROM Compensation 
	WHERE id=?;
	`

	row := db.QueryRow(getQuery, id)
	if err := row.Scan(&comp.Id, &comp.ResourceCode, &comp.Grade, &comp.LaborCategory, &comp.HourlyRate); err != nil {
		if err == sql.ErrNoRows {
			return comp, fmt.Errorf("compensation id=%d: no such row", id)
		}
		return comp, fmt.Errorf("compensation: id=%d: %v", id, err)
	}
	return comp, nil
}

func AllCompensation(db *sql.DB) ([]Compensation, error) {
	var comps []Compensation

	getQuery := `SELECT id,resource_code,grade,labor_category,hourly_rate FROM Compensation;`

	rows, err := db.Query(getQuery)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comp Compensation
		if err := rows.Scan(&comp.Id, &comp.ResourceCode, &comp.Grade, &comp.LaborCategory, &comp.HourlyRate); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		comps = append(comps, comp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return comps, nil
}

func UpdateCompensation(db *sql.DB, comp Compensation) (int64, error) {
	updateQuery := `
	UPDATE Compensation SET resource_code=?, grade=?, labor_category=?, hourly_rate=? WHERE id=?;
	`

	result, err := db.Exec(updateQuery, comp.ResourceCode, comp.Grade, comp.LaborCategory, comp.HourlyRate, comp.Id)
	if err != nil {
		return 0, fmt.Errorf("update query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update result error: %v", err)
	}
	return rows, nil
}

func InsertCompensation(db *sql.DB, comp Compensation) (int64, error) {
	insertQuery := `
	INSERT INTO Compensation (resource_code,grade,labor_category,hourly_rate) VALUES (?, ?, ?, ?);
	`

	result, err := db.Exec(insertQuery, comp.ResourceCode, comp.Grade, comp.LaborCategory, comp.HourlyRate)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return rows, nil
}
