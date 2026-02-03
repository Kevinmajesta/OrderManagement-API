package service

import (
	"errors"
	"time"

	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"

	"github.com/google/uuid"
)

type SalesReportService interface {
	GetSalesReportByDateRange(startDate, endDate time.Time) (*entity.SalesReport, error)
	GetDailySalesReport(date time.Time) (*entity.SalesReport, error)
	GetMonthlySalesReport(year int, month time.Month) (*entity.SalesReport, error)
}

type salesReportService struct {
	salesReportRepo repository.SalesReportRepository
}

func NewSalesReportService(salesReportRepo repository.SalesReportRepository) *salesReportService {
	return &salesReportService{
		salesReportRepo: salesReportRepo,
	}
}

func (s *salesReportService) GetSalesReportByDateRange(startDate, endDate time.Time) (*entity.SalesReport, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start_date must be before end_date")
	}

	reportData, err := s.salesReportRepo.GetSalesReportByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	paymentBreakdown, _ := s.salesReportRepo.GetPaymentMethodBreakdown(startDate, endDate)
	topProducts, _ := s.salesReportRepo.GetTopProducts(startDate, endDate, 10)

	report := &entity.SalesReport{
		ReportID:                uuid.New(),
		ReportDate:              time.Now(),
		PeriodStartDate:         startDate,
		PeriodEndDate:           endDate,
		TotalSales:              reportData["total_sales"].(float64),
		TotalTransactions:       reportData["total_transactions"].(int64),
		TotalTax:                reportData["total_tax"].(float64),
		AverageTransactionValue: reportData["average_transaction_value"].(float64),
		CashAmount:              reportData["cash_amount"].(float64),
		MidtransAmount:          reportData["midtrans_amount"].(float64),
		TotalCustomers:          reportData["total_customers"].(int64),
		PaymentMethodBreakdown:  paymentBreakdown,
		TopProducts:             topProducts,
		CreatedAt:               time.Now(),
	}

	return report, nil
}

func (s *salesReportService) GetDailySalesReport(date time.Time) (*entity.SalesReport, error) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.Add(24 * time.Hour)

	return s.GetSalesReportByDateRange(startDate, endDate)
}

func (s *salesReportService) GetMonthlySalesReport(year int, month time.Month) (*entity.SalesReport, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	return s.GetSalesReportByDateRange(startDate, endDate)
}
