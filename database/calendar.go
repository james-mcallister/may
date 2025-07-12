package database

import (
	"database/sql"
	"fmt"
)

type Calendar struct {
	Id          int64  `json:"id,string"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewCalendar() Calendar {
	return Calendar{}
}

func (c Calendar) Init(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS Calendar (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		description TEXT DEFAULT ''
	);
	`
	insertQuery := `
	INSERT OR IGNORE INTO Calendar 
	  (id, name, description) 
	VALUES 
	  (1, '9/80A', 'First Friday of pay period off.'),
	  (2, '9/80B', 'Second Friday of pay period off'),
	  (3, '4/10', 'All Fridays off'),
	  (4, '5/40', '8 hours each working day');
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

func CalendarDropdownQuery() string {
	return "SELECT name,id FROM Calendar;"
}

func GetCalendar(db *sql.DB, id int64) (Calendar, error) {
	var cal Calendar

	getQuery := `SELECT id,name,description FROM Calendar WHERE id=?;`

	row := db.QueryRow(getQuery, id)
	if err := row.Scan(&cal.Id, &cal.Name, &cal.Description); err != nil {
		if err == sql.ErrNoRows {
			return cal, fmt.Errorf("calendar id=%d: no such row", id)
		}
		return cal, fmt.Errorf("calendar: id=%d: %v", id, err)
	}
	return cal, nil
}

func AllCalendars(db *sql.DB) ([]Calendar, error) {
	var cals []Calendar

	getQuery := `SELECT id,name,description FROM Calendar ORDER BY name;`

	rows, err := db.Query(getQuery)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cal Calendar
		if err := rows.Scan(&cal.Id, &cal.Name, &cal.Description); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		cals = append(cals, cal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return cals, nil
}

func UpdateCalendar(db *sql.DB, cal Calendar) (int64, error) {
	updateQuery := `
	UPDATE Calendar SET name=?, description=? WHERE id=?;
	`

	result, err := db.Exec(updateQuery, cal.Name, cal.Description, cal.Id)
	if err != nil {
		return 0, fmt.Errorf("update query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update result error: %v", err)
	}
	return rows, nil
}

func InsertCalendar(db *sql.DB, cal Calendar) (int64, error) {
	insertQuery := `
	INSERT INTO Calendar (name,description) VALUES (?, ?);
	`

	result, err := db.Exec(insertQuery, cal.Name, cal.Description)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return rows, nil
}
