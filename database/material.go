package database

import (
	"database/sql"
	"fmt"
)

type Material struct {
	Id                  int64         `json:"id,string"`
	Name                string        `json:"name"`
	EstimatedCost       float64       `json:"estimated_cost"`
	ActualCost          float64       `json:"actual_cost"`
	PRDate              string        `json:"pr_date"`
	PODate              string        `json:"po_date"`
	PRNumber            string        `json:"pr"`
	PONumber            string        `json:"po"`
	Complete            bool          `json:"complete"`
	BaselineStartDate   string        `json:"baseline_start_date"`
	BaselineFinishDate  string        `json:"baseline_finish_date"`
	TentativeStartDate  string        `json:"tentative_start_date"`
	TentativeFinishDate string        `json:"tentative_finish_date"`
	ActualStartDate     string        `json:"actual_start_date"`
	ActualFinishDate    string        `json:"actual_finish_date"`
	Notes               string        `json:"notes"`
	WorkPackage         sql.NullInt64 `json:"wp"`
	WorkPackageName     string        `json:"wp_name"`
}

func NewMaterial() Material {
	return Material{}
}

func (e Material) Init(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS Material (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		estimated_cost NUMERIC DEFAULT 0.0,
		actual_cost NUMERIC DEFAULT 0.0,
		pr_date TEXT DEFAULT '',
		po_date TEXT DEFAULT '',
		pr_number TEXT DEFAULT '',
		po_number TEXT DEFAULT '',
		complete BOOLEAN DEFAULT FALSE,
		baseline_start_date TEXT DEFAULT '',
		baseline_finish_date TEXT DEFAULT '',
		tentative_start_date TEXT DEFAULT '',
		tentative_finish_date TEXT DEFAULT '',
		actual_start_date TEXT DEFAULT '',
		actual_finish_date TEXT DEFAULT '',
		notes TEXT DEFAULT '',
		proj INTEGER,
		FOREIGN KEY (proj) REFERENCES Project(id)
			ON DELETE SET NULL,
		UNIQUE (name, proj)
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

func MaterialDropdownQuery() string {
	return "SELECT name,id FROM Material;"
}

func GetMaterial(db *sql.DB, id int64) (Material, error) {
	var mat Material

	getQuery := `
	SELECT
	    id,name,estimated_cost,actual_cost,pr_date,po_date,pr_number,
		po_number,complete,baseline_start_date,baseline_finish_date,
		tentative_start_date,tentative_finish_date,actual_start_date,
		actual_finish_date,notes,proj
	FROM Material
	WHERE id=?;
	`

	row := db.QueryRow(getQuery, id)
	if err := row.Scan(&mat.Id, &mat.Name, &mat.EstimatedCost, &mat.ActualCost, &mat.PRDate, &mat.PODate, &mat.PRNumber, &mat.PONumber, &mat.Complete, &mat.BaselineStartDate, &mat.BaselineFinishDate, &mat.TentativeStartDate, &mat.TentativeFinishDate, &mat.ActualStartDate, &mat.ActualFinishDate, &mat.Notes, &mat.WorkPackage); err != nil {
		if err == sql.ErrNoRows {
			return mat, fmt.Errorf("material id=%d: no such row", id)
		}
		return mat, fmt.Errorf("material: id=%d: %v", id, err)
	}
	return mat, nil
}

func AllMaterials(db *sql.DB) ([]Material, error) {
	var mats []Material

	getQuery := `
	SELECT
	    id,name,estimated_cost,actual_cost,pr_date,po_date,pr_number,
		po_number,complete,baseline_start_date,baseline_finish_date,
		tentative_start_date,tentative_finish_date,actual_start_date,
		actual_finish_date,notes,proj
	FROM Material;
	`

	rows, err := db.Query(getQuery)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mat Material
		if err := rows.Scan(&mat.Id, &mat.Name, &mat.EstimatedCost, &mat.ActualCost, &mat.PRDate, &mat.PODate, &mat.PRNumber, &mat.PONumber, &mat.Complete, &mat.BaselineStartDate, &mat.BaselineFinishDate, &mat.TentativeStartDate, &mat.TentativeFinishDate, &mat.ActualStartDate, &mat.ActualFinishDate, &mat.Notes, &mat.WorkPackage); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		mats = append(mats, mat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return mats, nil
}

func UpdateMaterial(db *sql.DB, mat Material) (int64, error) {
	updateQuery := `
	UPDATE Material SET
	    name=?,estimated_cost=?,actual_cost=?,pr_date=?,po_date=?,pr_number=?,
		po_number=?,complete=?,baseline_start_date=?,baseline_finish_date=?,
		tentative_start_date=?,tentative_finish_date=?,actual_start_date=?,
		actual_finish_date=?,notes=?,proj=?
	WHERE id=?;
	`

	result, err := db.Exec(updateQuery, mat.Name, mat.EstimatedCost, mat.ActualCost, mat.PRDate, mat.PODate, mat.PRNumber, mat.PONumber, mat.Complete, mat.BaselineStartDate, mat.BaselineFinishDate, mat.TentativeStartDate, mat.TentativeFinishDate, mat.ActualStartDate, mat.ActualFinishDate, mat.Notes, mat.WorkPackage, mat.Id)
	if err != nil {
		return 0, fmt.Errorf("update query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update result error: %v", err)
	}
	return rows, nil
}

func InsertMaterial(db *sql.DB, mat Material) (int64, error) {
	insertQuery := `
	INSERT INTO Material
	  (name,estimated_cost,actual_cost,pr_date,po_date,pr_number,
		po_number,complete,baseline_start_date,baseline_finish_date,
		tentative_start_date,tentative_finish_date,actual_start_date,
		actual_finish_date,notes,proj)
	VALUES
	  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	result, err := db.Exec(insertQuery, mat.Name, mat.EstimatedCost, mat.ActualCost, mat.PRDate, mat.PODate, mat.PRNumber, mat.PONumber, mat.Complete, mat.BaselineStartDate, mat.BaselineFinishDate, mat.TentativeStartDate, mat.TentativeFinishDate, mat.ActualStartDate, mat.ActualFinishDate, mat.Notes, mat.WorkPackage)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return rows, nil
}
