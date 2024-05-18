package user

import (
	localJwt "halosuster/pkg/jwt"
	"halosuster/pkg/response"
	"halosuster/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	itAuth := group.Group("it")

	// Utillize group to use global setting on group parent (if exists)
	group.POST("nurse/login", h.NurseLogin)

	// Auth route for IT
	itAuth.POST("login")
}

func (h *userHandler) ITLogin(ctx *gin.Context) {
	var request ITLoginDTO

	// Create validator instance
	validator := validator.New(validator.WithRequiredStructEnabled())

	// Parse request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(err.Error()))
		return
	}

	// Validate input
	err := validator.Struct(request)
	if err != nil {
		validationErr := validation.FormatValidation(err)

		response.GenerateResponse(ctx, validationErr.Code, response.WithData(validationErr.Message))
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 200, response.WithData(request))
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
		ID:   nurse.ID,
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
