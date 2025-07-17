package handler

import (
	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/http/binder"
	"Kevinmajesta/OrderManagementAPI/internal/service"
	"Kevinmajesta/OrderManagementAPI/pkg/response"
	"Kevinmajesta/OrderManagementAPI/worker"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	photoFilename := photoID.String()
	photoPath := "/assets/photos/" + photoFilename + ext

	// Kirim ke worker (gunakan pipe agar bisa dikirim ulang)
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		io.Copy(pw, src)
	}()

	worker.PhotoQueue <- worker.PhotoJob{
		Src:      pr,
		Filename: photoFilename,
		Ext:      ext,
	}

	newProduct := &entity.Products{
		Name:        input.Name,
		Description: input.Description,
		PhotoURL:    photoPath,
		Price:       input.Price,
		Stock:       input.Stock,
	}

	product, err := h.productService.CreateProduct(newProduct)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully input a new product", product))
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	var input binder.ProductUpdateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "There is an input error"))
	}

	if input.ProductID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Product ID cannot be empty"))
	}

	// Ambil data produk berdasarkan ID
	oldProduct, err := h.productService.FindProductByID(input.ProductID.String())
	if err != nil || oldProduct == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Product ID does not exist"))
	}

	newPhotoURL := oldProduct.PhotoURL // default: pakai photo lama

	// Cek apakah ada file photo baru
	file, err := c.FormFile("photo")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid photo format. Only jpg, jpeg, and png are allowed"))
		}

		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to open new photo"))
		}
		defer src.Close()

		if oldProduct.PhotoURL != "" {
			oldPath := strings.TrimPrefix(oldProduct.PhotoURL, "/")
			_ = os.Remove(oldPath)
		}

		photoID := uuid.New()
		photoFilename := fmt.Sprintf("%s%s", photoID, ext)
		photoPath := filepath.Join("assets", "photos", photoFilename)

		if err := os.MkdirAll(filepath.Dir(photoPath), os.ModePerm); err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to create directory for photo"))
		}

		dst, err := os.Create(photoPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to create photo file"))
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to save new photo"))
		}

		newPhotoURL = "/assets/photos/" + photoFilename
	}

	// Validasi data lain
	if input.Name == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Product name cannot be empty"))
	}
	if input.Description == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Product description cannot be empty"))
	}
	if input.Price <= 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Price must be greater than 0"))
	}
	if input.Stock < 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Stock must be 0 or more"))
	}

	updatedProduct := entity.UpdateProduct(
		input.ProductID,
		input.Name,
		input.Description,
		newPhotoURL,
		input.Price,
		input.Stock,
	)

	result, err := h.productService.UpdateProduct(updatedProduct)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully updated product", result))
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	var input binder.ProductDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	isDeleted, err := h.productService.DeleteProduct(input.ProductID.String())
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete product", isDeleted))
}

func (h *ProductHandler) FindAllProduct(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	search := c.QueryParam("search")

	products, err := h.productService.FindAllProduct(page, search)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data products", products))

}
