package v1

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login godoc
// @Summary      用户登录
// @Description  用户使用用户名和密码进行登录认证
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "登录请求"
// @Success      200 {object} models.Response{data=models.LoginResponse} "登录成功"
// @Failure      200 {object} models.Response "登录失败"
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeBadRequest,
			Msg:  "Invalid request data",
		})
		return
	}

	resp, err := ctrl.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{
			Code: models.CodeUnauthorized,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code: models.CodeSuccess,
		Msg:  "Login successful",
		Data: resp,
	})
}
