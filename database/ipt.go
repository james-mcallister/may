package database

import (
	"database/sql"
	"fmt"
)

type Ipt struct {
	Id          int64  `json:"id,string"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewIpt() Ipt {
	return Ipt{}
}

func (i Ipt) Init(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS Ipt (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		description TEXT DEFAULT ''
	);
	`
	insertQuery := `
	INSERT OR IGNORE INTO Ipt
	  (name, description)
	VALUES
	  ('Software', 'A professional collective engaged in designing, building, deployment, and maintenance of IT software.'),
	  ('Test', 'A team of experts dedicated to ensuring the best quality of software products.'),
	  ('Integration', 'A team that specializes in bringing together component subsystems into a whole and ensuring that those subsystems function together.'),
	  ('Cyber', 'The cybersecurity team''s role is to protect an organization''s IT infrastructure from vulnerabilities and potential threats.'),
	  ('Network', 'The network team supports end users by configuring and maintaining computers, operating systems, and other applications that are connected to the network.'),
	  ('Logistics', 'A logistics team plays a crucial role in managing the supply chain, ensuring efficient coordination and management of resources during various operations.'),
	  ('Systems', 'Systems Engineering involves the top-down development of a system''s functional and physical requirements from a basic set of mission objectives.');
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

func IptDropdownQuery() string {
	return "SELECT name,id FROM Ipt;"
}

func GetIpt(db *sql.DB, id int64) (Ipt, error) {
	var ipt Ipt

	getQuery := `SELECT id,name,description FROM Ipt WHERE id=?;`

	row := db.QueryRow(getQuery, id)
	if err := row.Scan(&ipt.Id, &ipt.Name, &ipt.Description); err != nil {
		if err == sql.ErrNoRows {
			return ipt, fmt.Errorf("ipt id=%d: no such row", id)
		}
		return ipt, fmt.Errorf("ipt: id=%d: %v", id, err)
	}
	return ipt, nil
}

func AllIpts(db *sql.DB) ([]Ipt, error) {
	var ipts []Ipt

	getQuery := `SELECT id,name,description FROM Ipt;`

	rows, err := db.Query(getQuery)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ipt Ipt
		if err := rows.Scan(&ipt.Id, &ipt.Name, &ipt.Description); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		ipts = append(ipts, ipt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return ipts, nil
}

func UpdateIpt(db *sql.DB, ipt Ipt) (int64, error) {
	updateQuery := `
	UPDATE Ipt SET name=?, description=? WHERE id=?;
	`

	result, err := db.Exec(updateQuery, ipt.Name, ipt.Description, ipt.Id)
	if err != nil {
		return 0, fmt.Errorf("update query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update result error: %v", err)
	}
	return rows, nil
}

func InsertIpt(db *sql.DB, ipt Ipt) (int64, error) {
	insertQuery := `
	INSERT INTO Ipt (name,description) VALUES (?, ?);
	`

	result, err := db.Exec(insertQuery, ipt.Name, ipt.Description)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return rows, nil
}
