package repository

import (
	"time"

	"Kevinmajesta/OrderManagementAPI/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SalesReportRepository interface {
	GetSalesReportByDateRange(startDate, endDate time.Time) (map[string]interface{}, error)
	GetDailySalesReport(date time.Time) (map[string]interface{}, error)
	GetMonthlySalesReport(year int, month time.Month) (map[string]interface{}, error)
	GetPaymentMethodBreakdown(startDate, endDate time.Time) ([]entity.PaymentMethodStat, error)
	GetTopProducts(startDate, endDate time.Time, limit int) ([]entity.TopProductStat, error)
}

type salesReportRepository struct {
	db *gorm.DB
}

func NewSalesReportRepository(db *gorm.DB) SalesReportRepository {
	return &salesReportRepository{db: db}
}

func (r *salesReportRepository) GetSalesReportByDateRange(startDate, endDate time.Time) (map[string]interface{}, error) {
	var result map[string]interface{}

	// Get basic sales metrics
	var totalSales float64
	var totalTransactions int64
	var totalCustomers int64
	var cashAmount float64
	var midtransAmount float64

	r.db.Model(&entity.Order{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "paid").
		Select("SUM(total_price) as total, COUNT(DISTINCT order_id) as count, COUNT(DISTINCT user_id) as users").
		Row().
		Scan(&totalSales, &totalTransactions, &totalCustomers)

	// Get cash vs midtrans breakdown
	r.db.Model(&entity.Order{}).
		Where("created_at BETWEEN ? AND ? AND status = ? AND payment_method = ?", startDate, endDate, "paid", "cash").
		Select("COALESCE(SUM(total_price), 0)").Row().Scan(&cashAmount)

	r.db.Model(&entity.Order{}).
		Where("created_at BETWEEN ? AND ? AND status = ? AND payment_method = ?", startDate, endDate, "paid", "midtrans").
		Select("COALESCE(SUM(total_price), 0)").Row().Scan(&midtransAmount)

	result = map[string]interface{}{
		"total_sales":        totalSales,
		"total_transactions": totalTransactions,
		"total_customers":    totalCustomers,
		"cash_amount":        cashAmount,
		"midtrans_amount":    midtransAmount,
		"period_start_date":  startDate,
		"period_end_date":    endDate,
	}

	if totalTransactions > 0 {
		result["average_transaction_value"] = totalSales / float64(totalTransactions)
		result["total_tax"] = totalSales * 0.10
	}

	return result, nil
}

func (r *salesReportRepository) GetDailySalesReport(date time.Time) (map[string]interface{}, error) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.Add(24 * time.Hour)

	return r.GetSalesReportByDateRange(startDate, endDate)
}

func (r *salesReportRepository) GetMonthlySalesReport(year int, month time.Month) (map[string]interface{}, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	return r.GetSalesReportByDateRange(startDate, endDate)
}

func (r *salesReportRepository) GetPaymentMethodBreakdown(startDate, endDate time.Time) ([]entity.PaymentMethodStat, error) {
	var stats []entity.PaymentMethodStat

	var totalSales float64
	r.db.Model(&entity.Order{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "paid").
		Select("COALESCE(SUM(total_price), 0)").Row().Scan(&totalSales)

	rows, err := r.db.Model(&entity.Order{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "paid").
		Select("payment_method, COALESCE(SUM(total_price), 0) as total_amount, COUNT(*) as count").
		Group("payment_method").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var paymentMethod string
		var totalAmount float64
		var count int64
		rows.Scan(&paymentMethod, &totalAmount, &count)

		percentage := 0.0
		if totalSales > 0 {
			percentage = (totalAmount / totalSales) * 100
		}

		stats = append(stats, entity.PaymentMethodStat{
			PaymentMethod: paymentMethod,
			TotalAmount:   totalAmount,
			Count:         count,
			Percentage:    percentage,
		})
	}

	return stats, nil
}

func (r *salesReportRepository) GetTopProducts(startDate, endDate time.Time, limit int) ([]entity.TopProductStat, error) {
	var topProducts []entity.TopProductStat

	rows, err := r.db.Model(&entity.OrderItem{}).
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status = ?", startDate, endDate, "paid").
		Select("order_items.product_id, products.name, SUM(order_items.quantity) as qty, SUM(order_items.total_price) as revenue").
		Group("order_items.product_id, products.name").
		Order("SUM(order_items.quantity) DESC").
		Limit(limit).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productID uuid.UUID
		var productName string
		var qtySold int
		var revenue float64
		rows.Scan(&productID, &productName, &qtySold, &revenue)

		topProducts = append(topProducts, entity.TopProductStat{
			ProductID:    productID,
			ProductName:  productName,
			QuantitySold: qtySold,
			TotalRevenue: revenue,
		})
	}

	return topProducts, nil
}
