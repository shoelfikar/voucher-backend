package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/request"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/response"
	"github.com/shoelfikar/voucher-management-system/internal/domain/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles POST /api/login
// @Summary User login
// @Description Authenticate user with email and password (dummy validation)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=response.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid credentials"))
		return
	}

	loginResponse := response.LoginResponse{
		Token: token,
		User: response.UserInfo{
			Email: user.Email,
		},
	}

	c.JSON(http.StatusOK, response.SuccessResponse(loginResponse))
}
