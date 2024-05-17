package auth

import (
	localJwt "1-cat-social/pkg/jwt"
	"1-cat-social/pkg/response"
	"1-cat-social/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Auth handler structure for auth
type authHandler struct {
	uc IAuthUsecase
}

// Constructor for auth handler struct
func NewAuthHandler(uc IAuthUsecase) *authHandler {
	return &authHandler{
		uc: uc,
	}
}

// Router is required to wrap all user request by spesific path URL
func (h *authHandler) Router(r *gin.RouterGroup) {
	// Grouping to give URL prefix
	// ex : localhost/user
	group := r.Group("user")

	// Utillize group to use global setting on group parent (if exists)
	group.POST("nurse/login", h.NurseLogin)
}

func (h *authHandler) NurseLogin(ctx *gin.Context) {
	var request NurseLoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		validatorMessage := validation.GenerateStructValidationError(err)
		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
		return
	}

	nurse, err := h.uc.NurseLogin(request)
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		return
	}

	token, er := localJwt.GenerateToken(nurse)
	if er != nil {
		response.GenerateResponse(ctx, http.StatusInternalServerError, response.WithMessage("Failed to generate token"))
		return
	}

	res := FormatNurseLoginResponse(nurse, token)

	response.GenerateResponse(ctx, http.StatusOK, response.WithMessage("User loggedin successfully!"), response.WithData(res))
}
