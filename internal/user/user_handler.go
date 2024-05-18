package user

import (
	localJwt "halosuster/pkg/jwt"
	"halosuster/pkg/response"
	"halosuster/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	uc IUserUsecase
}

// Constructor for user handler struct
func NewUserHandler(uc IUserUsecase) *userHandler {
	return &userHandler{
		uc: uc,
	}
}

func (h *userHandler) Router(r *gin.RouterGroup) {
	// Grouping to give URL prefix
	// ex : localhost/user
	group := r.Group("user")

	// Utillize group to use global setting on group parent (if exists)
	group.POST("nurse/login", h.NurseLogin)
	group.POST("nurse/register", h.NurseRegister)
	group.POST("nurse/:id/access", h.NurseAccess)
}

func (h *userHandler) NurseLogin(ctx *gin.Context) {
	var request NurseLoginDTO

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

	tokenData := localJwt.TokenData{
		ID: nurse.ID,	
		Name: nurse.Name,
	}

	token, er := localJwt.GenerateToken(tokenData)
	if er != nil {
		response.GenerateResponse(ctx, http.StatusInternalServerError, response.WithMessage("Failed to generate token"))
		return
	}

	res := FormatNurseLoginResponse(nurse, token)

	response.GenerateResponse(ctx, http.StatusOK, response.WithMessage("User loggedin successfully!"), response.WithData(res))
}

func (h *userHandler) NurseRegister(ctx *gin.Context) {
	var request NurseRegisterDTO

	// Parse request body to DTO
	// If error return error response
	if err := ctx.ShouldBindJSON(&request); err != nil {
		validatorMessage := validation.GenerateStructValidationError(err)
		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
		return
	}

	nurse, err := h.uc.NurseRegister(request)
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		return
	}

	res := FormatNurseRegisterResponse(nurse)

	response.GenerateResponse(ctx, http.StatusOK, response.WithMessage("User register successfully!"), response.WithData(res))
}

func (h *userHandler) NurseAccess(ctx *gin.Context) {
	id := ctx.Param("id")

	var request NurseAccessDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		validatorMessage := validation.GenerateStructValidationError(err)
		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
		return
	}

	err := h.uc.NurseAccess(request, id)
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		return
	}

	response.GenerateResponse(ctx, http.StatusOK, response.WithMessage("User given access successfully!"))
}