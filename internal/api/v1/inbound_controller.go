package v1

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InboundController struct {
	inboundService *services.InboundService
}

func NewInboundController(inboundService *services.InboundService) *InboundController {
	return &InboundController{
		inboundService: inboundService,
	}
}

// GetAll godoc
// @Summary      获取所有入库订单
// @Description  分页获取入库订单列表，支持按供应商和日期筛选
// @Tags         入库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        req body models.GetInboundOrderRequest true "入库订单请求"
// @Success      200 {object} models.Response{data=models.GetInboundOrderResponse} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /inbound/orders/search [post]
func (ctrl *InboundController) GetAll(c *gin.Context) {
	var req models.GetInboundOrderRequest

	// 从查询参数绑定
	if err := c.Bind(&req); err != nil {
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

	orders, total, err := ctrl.inboundService.GetAll(&req)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	resp := models.GetInboundOrderResponse{
		Orders: orders,
		Total:  total,
	}

	c.JSON(http.StatusOK, &models.Response{Code: models.CodeSuccess, Msg: "success", Data: resp})
}

// Create godoc
// @Summary      创建入库订单
// @Description  创建新的入库订单
// @Tags         入库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order body models.CreateInboundOrderRequest true "入库订单信息"
// @Success      200 {object} models.Response{data=models.InboundOrder} "创建成功"
// @Failure      200 {object} models.Response "创建失败"
// @Router       /inbound/orders [post]
func (ctrl *InboundController) Create(c *gin.Context) {
	var req models.CreateInboundOrderRequest
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

	order, err := ctrl.inboundService.Create(&req, userModel.ID)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "Inbound order created successfully",
		Data: order,
	})
}

// GetByID godoc
// @Summary      根据ID获取入库订单详情
// @Description  根据订单ID获取入库订单详情，包含订单基本信息和详细条目
// @Tags         入库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "订单ID"
// @Success      200 {object} models.Response{data=models.GetInboudOrderDetailResp} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /inbound/orders/{id} [get]
func (ctrl *InboundController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid order ID",
		})
		return
	}

	order, err := ctrl.inboundService.GetByID(uint(id))
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

// Update godoc
// @Summary      更新入库订单
// @Description  根据ID更新入库订单信息
// @Tags         入库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "订单ID"
// @Param        order body models.InboundOrder true "入库订单信息"
// @Success      200 {object} models.Response{data=models.InboundOrder} "更新成功"
// @Failure      200 {object} models.Response "更新失败"
// @Router       /inbound/orders/{id} [put]
func (ctrl *InboundController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid order ID",
		})
		return
	}

	var order models.InboundOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid request data",
		})
		return
	}

	// 构建更新字段映射
	updates := make(map[string]interface{})
	if order.SupplierName != "" {
		updates["supplier_name"] = order.SupplierName
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

	if err := ctrl.inboundService.UpdateOrder(uint(id), updates); err != nil {
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

// Delete godoc
// @Summary      删除入库订单
// @Description  根据ID删除入库订单
// @Tags         入库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "订单ID"
// @Success      200 {object} models.Response "删除成功"
// @Failure      200 {object} models.Response "删除失败"
// @Router       /inbound/orders/{id} [delete]
func (ctrl *InboundController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid order ID",
		})
		return
	}

	if err := ctrl.inboundService.Delete(uint(id)); err != nil {
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
