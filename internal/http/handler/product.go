package handler

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/http/binder"
	"Kevinmajesta/OrderManagementAPI/internal/service"
	"Kevinmajesta/OrderManagementAPI/pkg/response"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	input := binder.ProductCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "There is an input error"))
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Failed to retrieve photo"))
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid photo format. Only jpg, jpeg, and png are allowed"))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to open photo"))
	}
	defer src.Close()

	photoID := uuid.New()
	photoFilename := fmt.Sprintf("%s%s", photoID, ext)
	photoPath := filepath.Join("assets", "photos", photoFilename)

	dst, err := os.Create(photoPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to create photo file"))
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to copy photo file"))
	}

	newProduct := &entity.Products{
		Name:        input.Name,
		Description: input.Description,
		PhotoURL:    "/assets/photos/" + photoFilename,
		Price:       input.Price,
		Stock:       input.Stock,
	}

	product, err := h.productService.CreateProduct(newProduct)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully input a new product", product))
}
