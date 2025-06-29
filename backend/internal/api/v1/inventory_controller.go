package v1

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InventoryController struct {
	inventoryService services.InventoryService
}

func NewInventoryController(inventoryService services.InventoryService) *InventoryController {
	return &InventoryController{
		inventoryService: inventoryService,
	}
}

// GetAll godoc
// @Summary      获取所有库存
// @Description  获取所有电池分类的库存信息
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.Response{data=[]models.Inventory} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /inventory [get]
func (ctrl *InventoryController) GetAll(c *gin.Context) {
	inventories, err := ctrl.inventoryService.GetAll()
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
		Data: inventories,
	})
}

// GetByCategoryID godoc
// @Summary      根据分类ID获取库存
// @Description  根据电池分类ID获取该分类的库存信息
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        categoryId path int true "分类ID"
// @Success      200 {object} models.Response{data=models.Inventory} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /inventory/{categoryId} [get]
func (ctrl *InventoryController) GetByCategoryID(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("categoryId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid category ID",
		})
		return
	}

	inventory, err := ctrl.inventoryService.GetByCategoryID(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeNotFound,
			Msg:  "Inventory not found",
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "success",
		Data: inventory,
	})
}
