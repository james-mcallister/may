package database

import (
	"database/sql"
	"fmt"
)

type Project struct {
	Id            int64         `json:"id,string"`
	Title         string        `json:"title"`        // Software Maintenance (SM)
	Description   string        `json:"description"`  // CLIN, Control Account, Work Package, etc.
	WbsId         string        `json:"wbsid"`        // 1.1002.5.2
	StmtOfWork    string        `json:"stmt_of_work"` // statement of work or WP description
	StartDate     string        `json:"start_date"`
	EndDate       string        `json:"end_date"`
	ImsUid        int           `json:"ims_uid,string"`    // Assigned by the scheduler
	WadLineId     int           `json:"wad_lineid,string"` // Is this for the PLATO tab?
	Evt           string        `json:"evt"`
	ParentProject sql.NullInt64 `json:"parent_proj"`
}

func NewProject() Project {
	return Project{}
}

func (p Project) Init(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS Project (
		id INTEGER PRIMARY KEY,
		title TEXT UNIQUE NOT NULL,
		description TEXT DEFAULT '',
		wbs_id TEXT DEFAULT '',
		stmt_of_work TEXT DEFAULT '',
		start_date TEXT DEFAULT '',
		end_date TEXT DEFAULT '',
		ims_uid INTEGER DEFAULT 0,
		wad_line_id INTEGER DEFAULT 0,
		evt TEXT DEFAULT '',
		parent_project INTEGER,
		FOREIGN KEY (parent_project) REFERENCES Project(id)
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

func ProjectDropdownQuery() string {
	return "SELECT wbs_id || ' - ' || title AS name,id FROM Project;"
}

func ProjectImportColumns() []string {
	return []string{"title", "description", "wbs_id", "stmt_of_work",
		"start_date", "end_date", "ims_uid", "evt"}
}

func GetProject(db *sql.DB, id int64) (Project, error) {
	var proj Project

	getQuery := `
	SELECT
	  id,title,description,wbs_id,stmt_of_work,start_date,
	  end_date,ims_uid,wad_line_id,evt,parent_project
	FROM Project
	WHERE id=?;
	`

	row := db.QueryRow(getQuery, id)
	if err := row.Scan(&proj.Id, &proj.Title, &proj.Description, &proj.WbsId, &proj.StmtOfWork, &proj.StartDate, &proj.EndDate, &proj.ImsUid, &proj.WadLineId, &proj.Evt, &proj.ParentProject); err != nil {
		if err == sql.ErrNoRows {
			return proj, fmt.Errorf("project id=%d: no such row", id)
		}
		return proj, fmt.Errorf("project: id=%d: %v", id, err)
	}

	return proj, nil
}

func AllProjects(db *sql.DB) ([]Project, error) {
	var projects []Project

	getQuery := `
	SELECT
	  id,title,description,wbs_id,stmt_of_work,start_date,
	  end_date,ims_uid,wad_line_id,evt,parent_project
	FROM Project;
	`

	rows, err := db.Query(getQuery)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var proj Project
		if err := rows.Scan(&proj.Id, &proj.Title, &proj.Description, &proj.WbsId, &proj.StmtOfWork, &proj.StartDate, &proj.EndDate, &proj.ImsUid, &proj.WadLineId, &proj.Evt, &proj.ParentProject); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}

		projects = append(projects, proj)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return projects, nil
}

func UpdateProject(db *sql.DB, proj Project) (int64, error) {
	updateQuery := `
	UPDATE Project SET
	  title=?,description=?,wbs_id=?,stmt_of_work=?,start_date=?,
	  end_date=?,ims_uid=?,wad_line_id=?,evt=?,parent_project=?
	WHERE id=?;
	`

	result, err := db.Exec(updateQuery, proj.Title, proj.Description, proj.WbsId, proj.StmtOfWork, proj.StartDate, proj.EndDate, proj.ImsUid, proj.WadLineId, proj.Evt, proj.ParentProject, proj.Id)
	if err != nil {
		return 0, fmt.Errorf("update query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("update result error: %v", err)
	}
	return rows, nil
}

func InsertProject(db *sql.DB, proj Project) (int64, error) {
	insertQuery := `
	INSERT INTO Project
	  (title,description,wbs_id,stmt_of_work,start_date,
	   end_date,ims_uid,wad_line_id,evt,parent_project)
	VALUES
	  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	result, err := db.Exec(insertQuery, proj.Title, proj.Description, proj.WbsId, proj.StmtOfWork, proj.StartDate, proj.EndDate, proj.ImsUid, proj.WadLineId, proj.Evt, proj.ParentProject)
	if err != nil {
		return 0, fmt.Errorf("insert query error: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("insert result error: %v", err)
	}
	return rows, nil
}
