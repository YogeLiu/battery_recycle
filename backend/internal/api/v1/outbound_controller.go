package v1

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OutboundController struct {
	outboundService *services.OutboundService
}

func NewOutboundController(outboundService *services.OutboundService) *OutboundController {
	return &OutboundController{
		outboundService: outboundService,
	}
}

// GetAll godoc
// @Summary      获取所有出库订单
// @Description  分页获取出库订单列表，支持按客户和日期筛选
// @Tags         出库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(20)
// @Param        customer query string false "客户名称 (支持模糊搜索)"
// @Param        start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param        end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success      200 {object} models.Response{data=models.GetOutboundOrderResponse} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /outbound/orders [get]
func (ctrl *OutboundController) GetAll(c *gin.Context) {
	var req models.GetOutboundOrderRequest

	// 从查询参数绑定
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid query parameters",
		})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	orders, total, err := ctrl.outboundService.GetAll(&req)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	resp := models.GetOutboundOrderResponse{
		Orders: orders,
		Total:  total,
	}

	c.JSON(http.StatusOK, &models.Response{Code: models.CodeSuccess, Msg: "success", Data: resp})
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

// GetByID godoc
// @Summary      根据ID获取出库订单详情
// @Description  根据订单ID获取出库订单详情，包含订单基本信息和详细条目
// @Tags         出库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "订单ID"
// @Success      200 {object} models.Response{data=models.GetOutboundOrderDetailResp} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /outbound/orders/{id} [get]
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

	// 构建更新字段映射
	updates := make(map[string]interface{})
	if order.CustomerName != "" {
		updates["customer_name"] = order.CustomerName
	}
	if order.Status != "" {
		updates["status"] = order.Status
	}
	if order.Notes != "" {
		updates["notes"] = order.Notes
	}
	if order.TotalAmount > 0 {
		updates["total_amount"] = order.TotalAmount
	}

	if err := ctrl.outboundService.UpdateOrder(uint(id), updates); err != nil {
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
