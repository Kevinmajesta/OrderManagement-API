package binder

type SalesReportDateRangeRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type SalesReportDailyRequest struct {
	Date string `json:"date"`
}

type SalesReportMonthlyRequest struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}
