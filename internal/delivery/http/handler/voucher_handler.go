package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/request"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/response"
	"github.com/shoelfikar/voucher-management-system/internal/domain/service"
)

type VoucherHandler struct {
	voucherService service.VoucherService
}

func NewVoucherHandler(voucherService service.VoucherService) *VoucherHandler {
	return &VoucherHandler{
		voucherService: voucherService,
	}
}

// GetAll handles GET /api/vouchers
// @Summary Get all vouchers
// @Description Get all vouchers with pagination, search, and sorting
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by voucher code"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order (asc/desc)" default(desc)
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.VoucherListResponse}
// @Failure 500 {object} response.Response
// @Router /api/vouchers [get]
func (h *VoucherHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	vouchers, total, err := h.voucherService.GetAll(page, limit, search, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	voucherListResponse := response.BuildVoucherListResponse(vouchers, page, limit, total)

	c.JSON(http.StatusOK, response.SuccessResponse(voucherListResponse))
}

// GetByID handles GET /api/vouchers/:id
// @Summary Get voucher by ID
// @Description Get a single voucher by its ID
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param id path int true "Voucher ID"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.VoucherResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/vouchers/{id} [get]
func (h *VoucherHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid voucher ID"))
		return
	}

	voucher, err := h.voucherService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err.Error()))
		return
	}

	voucherResponse := response.ToVoucherResponse(voucher)

	c.JSON(http.StatusOK, response.SuccessResponse(voucherResponse))
}

// Create handles POST /api/vouchers
// @Summary Create a new voucher
// @Description Create a new voucher with the provided details
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param request body request.CreateVoucherRequest true "Voucher details"
// @Security BearerAuth
// @Success 201 {object} response.Response{data=response.VoucherResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/vouchers [post]
func (h *VoucherHandler) Create(c *gin.Context) {
	var req request.CreateVoucherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	voucher, err := h.voucherService.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	voucherResponse := response.ToVoucherResponse(voucher)

	c.JSON(http.StatusCreated, response.SuccessResponseWithMessage("Voucher created successfully", voucherResponse))
}

// Update handles PUT /api/vouchers/:id
// @Summary Update a voucher
// @Description Update an existing voucher with the provided details
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param id path int true "Voucher ID"
// @Param request body request.UpdateVoucherRequest true "Voucher details"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.VoucherResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/vouchers/{id} [put]
func (h *VoucherHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid voucher ID"))
		return
	}

	var req request.UpdateVoucherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	voucher, err := h.voucherService.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	voucherResponse := response.ToVoucherResponse(voucher)

	c.JSON(http.StatusOK, response.SuccessResponseWithMessage("Voucher updated successfully", voucherResponse))
}

// Delete handles DELETE /api/vouchers/:id
// @Summary Delete a voucher
// @Description Soft delete a voucher by its ID
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param id path int true "Voucher ID"
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/vouchers/{id} [delete]
func (h *VoucherHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid voucher ID"))
		return
	}

	err = h.voucherService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseWithMessage("Voucher deleted successfully", nil))
}

// ImportCSV handles POST /api/vouchers/upload-csv
// @Summary Import vouchers from CSV
// @Description Upload a CSV file to bulk import vouchers
// @Tags Vouchers
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=service.ImportResult}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/vouchers/upload-csv [post]
func (h *VoucherHandler) ImportCSV(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("File is required"))
		return
	}
	defer file.Close()


	if !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Only CSV files are allowed"))
		return
	}

	// Validate file size (max 5MB)
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("File size exceeds 5MB"))
		return
	}

	result, err := h.voucherService.ImportVouchers(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseWithMessage("CSV import completed", result))
}

// UploadBatch handles POST /api/vouchers/upload-batch
// @Summary Upload batch of vouchers
// @Description Upload a batch of vouchers with duplicate checking
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param request body request.BatchUploadRequest true "Batch vouchers"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=service.BatchImportResult}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/vouchers/upload-batch [post]
func (h *VoucherHandler) UploadBatch(c *gin.Context) {
	var req request.BatchUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request"))
		return
	}

	// Validate batch size
	if len(req.Vouchers) > 1000 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Batch size exceeds 1000"))
		return
	}

	result, err := h.voucherService.ImportBatch(req.Vouchers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(result))
}

// ExportCSV handles GET /api/vouchers/export
// @Summary Export vouchers to CSV
// @Description Download all vouchers as a CSV file
// @Tags Vouchers
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {file} file
// @Failure 500 {object} response.Response
// @Router /api/vouchers/export [get]
func (h *VoucherHandler) ExportCSV(c *gin.Context) {
	data, err := h.voucherService.ExportVouchers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=vouchers.csv")
	c.Data(http.StatusOK, "text/csv", data)
}
