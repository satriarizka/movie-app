package handler

import (
	"fmt"
	_ "movie-app/internal/delivery/http/dto/response"
	"movie-app/internal/usecase"
	"movie-app/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportUC usecase.ReportUseCase
}

func NewReportHandler(reportUC usecase.ReportUseCase) *ReportHandler {
	return &ReportHandler{reportUC}
}

func (h *ReportHandler) GetTopMovies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	data, err := h.reportUC.GetTopMovies(limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Top movies report", data)
}

// GetRevenueReport godoc
// @Summary      Get revenue report
// @Description  See daily or monthly revenue (Admin Only)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        mode   query    string  false  "Report mode: 'day' or 'month'"
// @Success      200    {object} utils.APIResponse{data=[]response.DailyRevenueResponse}
// @Router       /reports/revenue [get]
// @Security     BearerAuth
func (h *ReportHandler) GetRevenueReport(c *gin.Context) {
	mode := c.Query("mode") // day atau month
	data, err := h.reportUC.GetRevenueReport(mode)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Revenue report", data)
}

// ExportRevenueCSV godoc
// @Summary      Export revenue CSV
// @Description  Download revenue report as CSV file (Admin Only)
// @Tags         Reports
// @Produce      text/csv
// @Param        mode   query    string  false  "Report mode: 'day' or 'month'"
// @Success      200    {file}   file
// @Router       /reports/revenue/export [get]
// @Security     BearerAuth
func (h *ReportHandler) ExportRevenueCSV(c *gin.Context) {
	mode := c.Query("mode")
	if mode == "" {
		mode = "day"
	}

	csvBytes, err := h.reportUC.GenerateRevenueCSV(mode)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate CSV", err.Error())
		return
	}

	// Set Headers untuk Download
	filename := fmt.Sprintf("revenue_report_%s.csv", mode)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/csv")

	// Kirim Bytes langsung
	c.Data(http.StatusOK, "text/csv", csvBytes)
}
