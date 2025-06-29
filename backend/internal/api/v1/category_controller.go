package v1

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// GetAll godoc
// @Summary      获取所有电池分类
// @Description  获取系统中所有电池分类的列表
// @Tags         电池分类管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.Response{data=[]models.BatteryCategory} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /categories [get]
func (ctrl *CategoryController) GetAll(c *gin.Context) {
	categories, err := ctrl.categoryService.GetAll()
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
		Data: categories,
	})
}

// Create godoc
// @Summary      创建电池分类
// @Description  创建新的电池分类
// @Tags         电池分类管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category body models.BatteryCategory true "电池分类信息"
// @Success      200 {object} models.Response{data=models.BatteryCategory} "创建成功"
// @Failure      200 {object} models.Response "创建失败"
// @Router       /categories [post]
func (ctrl *CategoryController) Create(c *gin.Context) {
	var category models.BatteryCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid request data",
		})
		return
	}

	if err := ctrl.categoryService.Create(&category); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "Category created successfully",
		Data: category,
	})
}

// GetByID godoc
// @Summary      根据ID获取电池分类
// @Description  根据分类ID获取电池分类信息
// @Tags         电池分类管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "分类ID"
// @Success      200 {object} models.Response{data=models.BatteryCategory} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /categories/{id} [get]
func (ctrl *CategoryController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid category ID",
		})
		return
	}

	category, err := ctrl.categoryService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeNotFound,
			Msg:  "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "success",
		Data: category,
	})
}

// Update godoc
// @Summary      更新电池分类
// @Description  根据ID更新电池分类信息
// @Tags         电池分类管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "分类ID"
// @Param        category body models.BatteryCategory true "电池分类信息"
// @Success      200 {object} models.Response{data=models.BatteryCategory} "更新成功"
// @Failure      200 {object} models.Response "更新失败"
// @Router       /categories/{id} [put]
func (ctrl *CategoryController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid category ID",
		})
		return
	}

	var category models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid request data",
		})
		return
	}

	// 构建更新字段映射
	updates := make(map[string]interface{})
	if category.Name != "" {
		updates["name"] = category.Name
	}
	if category.Description != "" {
		updates["description"] = category.Description
	}
	if category.UnitPrice > 0 {
		updates["unit_price"] = category.UnitPrice
	}

	if err := ctrl.categoryService.UpdateCategory(uint(id), updates); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "Category updated successfully",
		Data: category,
	})
}

// Delete godoc
// @Summary      删除电池分类
// @Description  根据ID删除电池分类
// @Tags         电池分类管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "分类ID"
// @Success      200 {object} models.Response "删除成功"
// @Failure      200 {object} models.Response "删除失败"
// @Router       /categories/{id} [delete]
func (ctrl *CategoryController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid category ID",
		})
		return
	}

	if err := ctrl.categoryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "Category deleted successfully",
	})
}
