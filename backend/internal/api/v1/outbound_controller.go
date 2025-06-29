package v1

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OutboundController struct {
	outboundService services.OutboundService
}

func NewOutboundController(outboundService services.OutboundService) *OutboundController {
	return &OutboundController{
		outboundService: outboundService,
	}
}

func (ctrl *OutboundController) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	orders, total, err := ctrl.outboundService.GetAll(limit, offset)
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
		Data: map[string]interface{}{
			"orders": orders,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (ctrl *OutboundController) Create(c *gin.Context) {
	var req models.CreateOutboundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid request data",
		})
		return
	}

	// Get current user from context
	user, _ := c.Get("user")
	userModel := user.(*models.User)

	order, err := ctrl.outboundService.Create(&req, userModel.ID)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "Outbound order created successfully",
		Data: order,
	})
}

func (ctrl *OutboundController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid order ID",
		})
		return
	}

	order, err := ctrl.outboundService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeNotFound,
			Msg:  "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "success",
		Data: order,
	})
}

func (ctrl *OutboundController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid order ID",
		})
		return
	}

	var order models.OutboundOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid request data",
		})
		return
	}

	order.ID = uint(id)
	if err := ctrl.outboundService.Update(&order); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "Order updated successfully",
		Data: order,
	})
}

func (ctrl *OutboundController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid order ID",
		})
		return
	}

	if err := ctrl.outboundService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "Order deleted successfully",
	})
}
