package staff

import (
	localJwt "1-cat-social/pkg/jwt"
	"1-cat-social/pkg/response"
	"1-cat-social/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type staffHandler struct {
	uc IStaffUsecase
}

// Constructor for staff handler struct
func NewStaffHandler(uc IStaffUsecase) *staffHandler {
	return &staffHandler{
		uc: uc,
	}
}

// Router is required to wrap all user request by spesific path URL
func (h *staffHandler) Router(r *gin.RouterGroup) {
	// Grouping to give URL prefix
	group := r.Group("staff")

	// Utillize group to use global setting on group parent (if exists)
	group.POST("nurse/login", h.NurseLogin)
}

func (h *authHandler) NurseLogin(ctx *gin.Context) {
	var request NurseLoginRequest

	// Parse request body to DTO
	// If error return error response
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
