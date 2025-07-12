package database

import (
	"database/sql"
	"fmt"
)

type Network struct {
	Id           int64         `json:"id,string"`
	ChargeNumber string        `json:"charge_num"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Status       string        `json:"status"` // ex: 20 - Labor Only
	StartDate    string        `json:"start_date"`
	EndDate      string        `json:"end_date"`
	Proj         sql.NullInt64 `json:"proj"`
	ProjName     string        `json:"proj_name"`
}

func NewNetwork() Network {
	return Network{}
}

func (n Network) Init(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS Network (
		id INTEGER PRIMARY KEY,
		charge_number TEXT UNIQUE NOT NULL,
		title TEXT DEFAULT '',
		description TEXT DEFAULT '',
		status TEXT DEFAULT '40 - Open All', -- ex: 20 - Labor Only
		start_date TEXT DEFAULT '',
		end_date TEXT DEFAULT '',
		proj INTEGER,
		FOREIGN KEY (proj) REFERENCES Project(id)
			ON DELETE SET NULL
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

func NetworkDropdownQuery() string {
	return "SELECT charge_number || ' - ' || title AS name,id FROM Network ORDER BY name;"
}

func NetworkImportColumns() []string {
	return []string{"charge_number", "title", "description", "status",
		"start_date", "end_date"}
}

func GetNetwork(db *sql.DB, id int64) (Network, error) {
	var net Network

	getQuery := `
	SELECT
	  id,charge_number,title,description,status,start_date,end_date,proj
	FROM Network
	WHERE id=?;
	`

	row := db.QueryRow(getQuery, id)
	if err := row.Scan(&net.Id, &net.ChargeNumber, &net.Title, &net.Description, &net.Status, &net.StartDate, &net.EndDate, &net.Proj); err != nil {
		if err == sql.ErrNoRows {
			return net, fmt.Errorf("network id=%d: no such row", id)
		}
		return net, fmt.Errorf("network: id=%d: %v", id, err)
	}

	return net, nil
}

func AllNetworks(db *sql.DB) ([]Network, error) {
	var nets []Network

	getQuery := `
	SELECT
	  id,charge_number,title,description,status,start_date,end_date,proj
	FROM Network
	ORDER BY charge_number;
	`

	rows, err := db.Query(getQuery)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var net Network
		if err := rows.Scan(&net.Id, &net.ChargeNumber, &net.Title, &net.Description, &net.Status, &net.StartDate, &net.EndDate, &net.Proj); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}

		nets = append(nets, net)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return nets, nil
}

func UpdateNetwork(db *sql.DB, net Network) (int64, error) {
	updateQuery := `
	UPDATE Network SET
	  charge_number=?,title=?,description=?,status=?,start_date=?,
	  end_date=?,proj=?
	WHERE id=?;
	`

	result, err := db.Exec(updateQuery, net.ChargeNumber, net.Title, net.Description, net.Status, net.StartDate, net.EndDate, net.Proj, net.Id)
	if err != nil {
		return 0, fmt.Errorf("update query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update result error: %v", err)
	}
	return rows, nil
}

func InsertNetwork(db *sql.DB, net Network) (int64, error) {
	insertQuery := `
	INSERT INTO Network
	  (charge_number,title,description,status,start_date,end_date,proj)
	VALUES
	  (?, ?, ?, ?, ?, ?, ?);
	`

	result, err := db.Exec(insertQuery, net.ChargeNumber, net.Title, net.Description, net.Status, net.StartDate, net.EndDate, net.Proj)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return rows, nil
}
