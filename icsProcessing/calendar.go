package icsProcessing

type CalendarTemplate struct {
	ProdId      string          `json:"prod_id"`
	Version     string          `json:"version"`
	CalScale    string          `json:"cal_scale"`
	Method      string          `json:"method"`
	CalName     string          `json:"cal_name"`
	CalTimeZone string          `json:"cal_timezone"`
	Events      []EventTemplate `json:"events"`
}

type EventTemplate struct {
	DateTimeStart  string `json:"date_time_start"`
	DateTimeEnd    string `json:"date_time_end"`
	DateTimeStamp  string `json:"date_time_stamp"`
	Uid            string `json:"uid"`
	Classification string `json:"classification"`
	Sequence       string `json:"sequence"`
	Status         string `json:"status"`
	Summary        string `json:"summary"`
	Transparency   string `json:"transparency"`
}
