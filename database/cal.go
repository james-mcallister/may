package database

// Information for the calendar view
type Cal struct {
	Period string   `json:"period"` // Month Year
	Name   string   `json:"name"`   // Emp display name or cal name
	Days   []CalDay `json:"days"`
	Outer  []int    // outer loop count (len days / 7. should be 4 or 5)
	Inner  []int    // inner loop count (7 since there are 7 days in a week)
}

type CalDay struct {
	CalDate   string  `json:"cal_date"`
	DayNum    int     `json:"day_num"`
	PlanHours float64 `json:"plan_hours"` // might not need this.
	ProdHours float64 `json:"prod_hours"`
	Disabled  bool    `json:"disabled"`
}

func NewCal() Cal {
	return Cal{}
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
