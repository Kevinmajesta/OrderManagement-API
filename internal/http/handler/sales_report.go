package handler

import (
	"net/http"
	"strconv"
	"time"

	"Kevinmajesta/OrderManagementAPI/internal/http/binder"
	"Kevinmajesta/OrderManagementAPI/internal/service"
	"Kevinmajesta/OrderManagementAPI/pkg/response"

	"github.com/labstack/echo/v4"
)

type SalesReportHandler struct {
	salesReportService service.SalesReportService
}

func NewSalesReportHandler(salesReportService service.SalesReportService) *SalesReportHandler {
	return &SalesReportHandler{salesReportService: salesReportService}
}

func (h *SalesReportHandler) GetSalesReportByDateRange(c echo.Context) error {
	var req binder.SalesReportDateRangeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid request"))
	}

	if req.StartDate == "" || req.EndDate == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "start_date and end_date are required"))
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD"))
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD"))
	}

	report, err := h.salesReportService.GetSalesReportByDateRange(startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sales report generated", report))
}

func (h *SalesReportHandler) GetDailySalesReport(c echo.Context) error {
	dateStr := c.QueryParam("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid date format, use YYYY-MM-DD"))
	}

	report, err := h.salesReportService.GetDailySalesReport(date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "daily sales report generated", report))
}

func (h *SalesReportHandler) GetMonthlySalesReport(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err == nil {
			year = y
		}
	}

	if monthStr != "" {
		m, err := strconv.Atoi(monthStr)
		if err == nil && m >= 1 && m <= 12 {
			month = m
		}
	}

	report, err := h.salesReportService.GetMonthlySalesReport(year, time.Month(month))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "monthly sales report generated", report))
}
