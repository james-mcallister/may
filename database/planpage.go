package database

import (
	"database/sql"
	"fmt"
)

type PlanPage struct {
	Id          int64   `json:"id,string"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	TargetCost  float64 `json:"target_cost"`
	TargetHours float64 `json:"target_hours"`
}

func NewLaborPlan() PlanPage {
	return PlanPage{}
}

func (p PlanPage) Init(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS PlanPage (
	    id INTEGER PRIMARY KEY,
		title TEXT UNIQUE NOT NULL,
		description TEXT DEFAULT '',
		target_cost NUMERIC DEFAULT 0.0,
		target_hours NUMERIC DEFAULT 0.0
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

func PlanPageDropdownQuery() string {
	return "SELECT title AS name,id FROM PlanPage ORDER BY name;"
}

func GetPlanPage(db *sql.DB, id int64) (PlanPage, error) {
	var plan PlanPage

	getQuery := `SELECT id,title,description,target_cost,target_hours FROM PlanPage WHERE id=?;`

	row := db.QueryRow(getQuery, id)
	if err := row.Scan(&plan.Id, &plan.Title, &plan.Description, &plan.TargetCost, &plan.TargetHours); err != nil {
		if err == sql.ErrNoRows {
			return plan, fmt.Errorf("plan page id=%d: no such row", id)
		}
		return plan, fmt.Errorf("plan page: id=%d: %v", id, err)
	}
	return plan, nil
}

func InsertPlanPage(db *sql.DB, plan PlanPage) (int64, error) {
	insertQuery := `
	INSERT INTO PlanPage
	  (title,description)
	VALUES
	  (?, ?);
	`

	result, err := db.Exec(insertQuery, plan.Title, plan.Description)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return rows, nil
}

func GetTargetValues(db *sql.DB, id int64) (float64, float64, error) {
	getQuery := `SELECT target_cost,target_hours FROM PlanPage WHERE id=?;`

	row := db.QueryRow(getQuery, id)
	var targetCost, targetHours float64
	if err := row.Scan(&targetCost, &targetHours); err != nil {
		if err == sql.ErrNoRows {
			return 0.0, 0.0, fmt.Errorf("target values id=%d: no such row", id)
		}
		return 0.0, 0.0, fmt.Errorf("target values: id=%d: %v", id, err)
	}
	return targetCost, targetHours, nil
}

func UpdateTargetValues(db *sql.DB, id int64, targetCost, targetHours float64) (int64, error) {
	updateQuery := `
	UPDATE PlanPage SET target_cost=?,target_hours=? WHERE id=?;
	`

	result, err := db.Exec(updateQuery, targetCost, targetHours, id)
	if err != nil {
		return 0, fmt.Errorf("update query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update result error: %v", err)
	}
	return rows, nil
}
