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

// Create godoc
// @Summary      创建出库订单
// @Description  创建新的出库订单，支持批量创建
// @Tags         出库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateOutboundOrderRequest true "创建出库订单请求"
// @Success      200 {object} models.Response "创建成功"
// @Failure      200 {object} models.Response "创建失败"
// @Router       /outbound/orders [post]
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

// Update godoc
// @Summary      更新出库订单
// @Description  根据订单ID更新出库订单
// @Tags         出库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "订单ID"
// @Param        request body models.UpdateOutboundOrderRequest true "更新出库订单请求"
// @Success      200 {object} models.Response "更新成功"
// @Failure      200 {object} models.Response "更新失败"
// @Router       /outbound/orders/{id} [put]
func (ctrl *OutboundController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid order ID",
		})
		return
	}

	var req models.UpdateOutboundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid request data: " + err.Error(),
		})
		return
	}

	// 判断是否需要更新订单项
	if len(req.Items) > 0 {
		// 完整更新（包括订单项）
		if err := ctrl.outboundService.UpdateOrderComplete(uint(id), &req); err != nil {
			c.JSON(http.StatusOK, &models.Response{
				Code: models.CodeInternalError,
				Msg:  err.Error(),
			})
			return
		}
	} else {
		// 仅更新基本信息
		if err := ctrl.outboundService.UpdateOrderBasic(uint(id), &req); err != nil {
			c.JSON(http.StatusOK, &models.Response{
				Code: models.CodeInternalError,
				Msg:  err.Error(),
			})
			return
		}
	}

	// 获取更新后的订单详情
	updatedOrder, err := ctrl.outboundService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeSuccess,
			Msg:  "Order updated successfully, but failed to retrieve updated data",
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "Order updated successfully",
		Data: updatedOrder,
	})
}

// Delete godoc
// @Summary      删除出库订单
// @Description  根据订单ID删除出库订单
// @Tags         出库管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "订单ID"
// @Success      200 {object} models.Response "删除成功"
// @Failure      200 {object} models.Response "删除失败"
// @Router       /outbound/orders/{id} [delete]
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
