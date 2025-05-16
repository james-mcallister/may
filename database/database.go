package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Connstring string
}

func NewDB() (Database, error) {
	d := Database{}
	p, err := dbPath()
	if err != nil {
		return d, fmt.Errorf("error creating db file: %v", err)
	}
	d.Connstring = "file:" + p + "?mode=rwc&_foreign_keys=true"
	return d, nil
}

func (d Database) Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", d.Connstring)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func dbPath() (string, error) {
	baseDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("dbpath error: %v", err)
	}
	appDir := "may_app"
	dbFile := "may.db"
	fullPath := filepath.Join(baseDir, appDir)

	err = os.MkdirAll(fullPath, 0755)
	if err != nil {
		return "", fmt.Errorf("dbpath mkdir error: %v", err)
	}
	return filepath.Join(fullPath, dbFile), nil
}

type DBTable interface {
	Init(*sql.DB) error
}

func InitDB(db *sql.DB) error {
	tables := []DBTable{
		Calendar{},
		CalendarHours{},
		Compensation{},
		Ipt{},
		Employee{},
		Project{},
		Network{},
	}

	for _, t := range tables {
		err := t.Init(db)
		if err != nil {
			return fmt.Errorf("database init error: %v", err)
		}
	}
	return nil
}

type Dropdown struct {
	Id   int64  `json:"id,string"`
	Name string `json:"name"`
}

func NewDropdown(db *sql.DB, query string) ([]Dropdown, error) {
	var names []Dropdown

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("dropdown query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name Dropdown
		if err := rows.Scan(&name.Name, &name.Id); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("dropdown query error: No rows")
			}
			return nil, fmt.Errorf("dropdown query error: %v", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dropdown query error: %v", err)
	}
	return names, nil
}

func DeleteRow(db *sql.DB, table string, id int64) (int64, error) {
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id=?;", table)

	result, err := db.Exec(deleteQuery, id)
	if err != nil {
		return 0, fmt.Errorf("delete query exec error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("delete query result error: %v", err)
	}
	return rows, nil
}

func DeleteAll(db *sql.DB, table string) (int64, error) {
	deleteAllQuery := fmt.Sprintf("DELETE FROM %s;", table)

	result, err := db.Exec(deleteAllQuery)
	if err != nil {
		return 0, fmt.Errorf("delete all query exec error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("delete all query result error: %v", err)
	}
	return rows, nil
}
