package v1

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService services.ReportService
}

func NewReportController(reportService services.ReportService) *ReportController {
	return &ReportController{
		reportService: reportService,
	}
}

func (ctrl *ReportController) GetSummary(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	summary, err := ctrl.reportService.GetSummary(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "success",
		Data: summary,
	})
}
