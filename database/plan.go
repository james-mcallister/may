package database

import (
	"database/sql"
	"fmt"
)

type Plan struct {
	Id        int64  `json:"id,string"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func NewPlan() Plan {
	return Plan{}
}

func (p Plan) Init(db *sql.DB) error {
	// data for the plan table
	tableQuery := `
	CREATE TABLE IF NOT EXISTS Plan (
	    id INTEGER PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		start_date TEXT DEFAULT '',
		end_date TEXT DEFAULT '',
		plan INTEGER,
		FOREIGN KEY (plan) REFERENCES PlanPage(id)
			ON DELETE CASCADE
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

func PlanDropdownQuery() string {
	return "SELECT name,id FROM Plan ORDER BY name;"
}

func GetPlan(db *sql.DB, planId int64) (Plan, error) {
	var t Plan

	getQuery := `
	SELECT id,name,start_date,end_date FROM Plan WHERE plan=?;
	`

	row := db.QueryRow(getQuery, planId)
	if err := row.Scan(&t.Id, &t.Name, &t.StartDate, &t.EndDate); err != nil {
		if err == sql.ErrNoRows {
			return t, fmt.Errorf("plan table id=%d: no such row", planId)
		}
		return t, fmt.Errorf("plan table: id=%d: %v", planId, err)
	}
	return t, nil
}

func InsertPlan(db *sql.DB, t Plan) (int64, error) {
	insertQuery := `
	INSERT INTO Plan
	  (name,start_date,end_date)
	VALUES
	  (?, ?, ?);
	`

	result, err := db.Exec(insertQuery, t.Name, t.StartDate, t.EndDate)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return id, nil
}
