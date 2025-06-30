package v1

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetAll godoc
// @Summary      获取所有用户
// @Description  获取系统中所有用户的列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.Response{data=[]models.User} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /users [get]
func (ctrl *UserController) GetAll(c *gin.Context) {
	users, err := ctrl.userService.GetAll()
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
		Data: users,
	})
}

// Create godoc
// @Summary      创建用户
// @Description  创建新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user body models.User true "用户信息"
// @Success      200 {object} models.Response{data=models.User} "创建成功"
// @Failure      200 {object} models.Response "创建失败"
// @Router       /users [post]
func (ctrl *UserController) Create(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid request data",
		})
		return
	}

	if err := ctrl.userService.Create(&user); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "User created successfully",
		Data: user,
	})
}

// GetByID godoc
// @Summary      根据ID获取用户
// @Description  根据用户ID获取用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID"
// @Success      200 {object} models.Response{data=models.User} "获取成功"
// @Failure      200 {object} models.Response "获取失败"
// @Router       /users/{id} [get]
func (ctrl *UserController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid user ID",
		})
		return
	}

	user, err := ctrl.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeNotFound,
			Msg:  "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "success",
		Data: user,
	})
}

// Update godoc
// @Summary      更新用户
// @Description  根据ID更新用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID"
// @Param        user body models.User true "用户信息"
// @Success      200 {object} models.Response{data=models.User} "更新成功"
// @Failure      200 {object} models.Response "更新失败"
// @Router       /users/{id} [put]
func (ctrl *UserController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid user ID",
		})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid request data",
		})
		return
	}

	// 构建更新字段映射
	updates := make(map[string]interface{})
	if user.Username != "" {
		updates["username"] = user.Username
	}
	if user.RealName != "" {
		updates["real_name"] = user.RealName
	}
	if user.Role != "" {
		updates["role"] = user.Role
	}
	if user.Password != "" {
		// 密码将在 service 层自动加密
		updates["password"] = user.Password
	}

	if err := ctrl.userService.UpdateUser(uint(id), updates); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "User updated successfully",
		Data: user,
	})
}

// Delete godoc
// @Summary      删除用户
// @Description  根据ID删除用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID"
// @Success      200 {object} models.Response "删除成功"
// @Failure      200 {object} models.Response "删除失败"
// @Router       /users/{id} [delete]
func (ctrl *UserController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid user ID",
		})
		return
	}

	if err := ctrl.userService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeInternalError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "User deleted successfully",
	})
}
