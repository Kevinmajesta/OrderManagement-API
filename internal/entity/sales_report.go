package entity

import (
	"time"

	"github.com/google/uuid"
)

type SalesReport struct {
	ReportID                uuid.UUID           `json:"report_id"`
	ReportDate              time.Time           `json:"report_date"`
	PeriodStartDate         time.Time           `json:"period_start_date"`
	PeriodEndDate           time.Time           `json:"period_end_date"`
	TotalSales              float64             `json:"total_sales"`
	TotalTransactions       int64               `json:"total_transactions"`
	TotalTax                float64             `json:"total_tax"`
	AverageTransactionValue float64             `json:"average_transaction_value"`
	CashAmount              float64             `json:"cash_amount"`
	MidtransAmount          float64             `json:"midtrans_amount"`
	TotalCustomers          int64               `json:"total_customers"`
	PaymentMethodBreakdown  []PaymentMethodStat `json:"payment_method_breakdown"`
	TopProducts             []TopProductStat    `json:"top_products"`
	CreatedAt               time.Time           `json:"created_at"`
}

type PaymentMethodStat struct {
	PaymentMethod string  `json:"payment_method"`
	TotalAmount   float64 `json:"total_amount"`
	Count         int64   `json:"count"`
	Percentage    float64 `json:"percentage"`
}

type TopProductStat struct {
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	QuantitySold int       `json:"quantity_sold"`
	TotalRevenue float64   `json:"total_revenue"`
}

type SalesReportRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
