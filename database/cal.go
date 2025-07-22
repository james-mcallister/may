package database

import (
	"database/sql"
	"fmt"
)

// Information for the calendar view
type Cal struct {
	Period string   `json:"period"` // Month Year
	Name   string   `json:"name"`   // Emp display name or cal name
	Days   []CalDay `json:"days"`
	Outer  []int    // outer loop count (len days / 7. should be 4 or 5)
	Inner  []int    // inner loop count (7 since there are 7 days in a week)
}

func NewCal() Cal {
	return Cal{}
}

type CalDay struct {
	DayId     int64   `json:"id,string"`
	CalDate   string  `json:"cal_date"`
	DayNum    string  `json:"day_num"`
	ProdHours float64 `json:"prod_hours"`
	Disabled  bool    `json:"disabled"`
}

func GetCalData(db *sql.DB, popStart, popEnd, fiscalPeriod string) ([]CalDay, error) {
	var days []CalDay

	getQuery := `
	SELECT id,cal_date,productive_hours
	FROM CalendarHours
	WHERE fiscal_period=?
	ORDER BY cal_date;
	`

	rows, err := db.Query(getQuery, fiscalPeriod)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d CalDay
		if err := rows.Scan(&d.DayId, &d.CalDate, &d.ProdHours); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		d.DayNum = GetDayNum(d.CalDate)
		days = append(days, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	if err := disableDays(db, days, popStart, popEnd); err != nil {
		return nil, fmt.Errorf("server error in disable days: %v", err)
	}
	return days, nil
}

func GetDayNum(isodate string) string {
	var dayNum string
	if isodate[8] == '0' {
		dayNum = isodate[9:10]
	} else {
		dayNum = isodate[8:10]
	}
	return dayNum
}

func disableDays(db *sql.DB, days []CalDay, popStart, popEnd string) error {
	// get the id for the pop start date and pop end date. The id is a monotonically
	// increasing integer assigned by sqlite. Lower id numbers should be earlier
	// dates simplifying date comparison for disabling days in the calendar view
	popDates, err := checkPop(db, popStart, popEnd)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}
	// result should have pop start and pop end
	if len(popDates) != 2 {
		return fmt.Errorf("query error: len pop dates incorrect")
	}

	last := len(days) - 1
	startId := popDates[0].DateId
	endId := popDates[1].DateId
	if (days[0].DayId >= startId) && (days[last].DayId <= endId) {
		return nil
	}

	for i, v := range days {
		if (v.DayId < startId) || (v.DayId > endId) {
			days[i].Disabled = true
		}
	}
	return nil
}

type DateId struct {
	CalDate string
	DateId  int64
}

func checkPop(db *sql.DB, popStart, popEnd string) ([]DateId, error) {
	var popDates []DateId

	getQuery := `
	SELECT id,cal_date
	FROM CalendarHours
	WHERE cal_date IN (?, ?)
	ORDER BY id;
	`

	rows, err := db.Query(getQuery, popStart, popEnd)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d DateId
		if err := rows.Scan(&d.DateId, &d.CalDate); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: no rows")
			}
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		popDates = append(popDates, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return popDates, nil
}
