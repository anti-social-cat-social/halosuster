package user

import (
	"halosuster/internal/middleware"
	"halosuster/pkg/jwt"
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
	nurseAuth := group.Group("nurse")

	// Utillize group to use global setting on group parent (if exists)
	group.POST("nurse/login", h.NurseLogin)

	// Auth route for IT
	itAuth.POST("login", h.ITLogin)
	group.POST("nurse/register", h.NurseRegister)
	group.POST("nurse/:id/access", h.NurseAccess)

	group.GET("", middleware.UseJwtAuth, middleware.HasRoles(string(IT)), h.GetUsers)
	nurseAuth.DELETE("/:id", middleware.UseJwtAuth, middleware.HasRoles(string(IT)), h.Delete)
}

func (h *userHandler) ITLogin(ctx *gin.Context) {
	var request ITLoginDTO

	// Create validator instance
	validator := validator.New(validator.WithRequiredStructEnabled())

	// Parse request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(err.Error()))
		ctx.Abort()
		return
	}

	// Validate input
	validator.RegisterValidation("valid_nip", validation.ValidNIP(string(ITPrefix)))

	err := validator.Struct(request)
	if err != nil {
		validationErr := validation.FormatValidation(err)

		response.GenerateResponse(ctx, validationErr.Code, response.WithData(validationErr.Message))
		ctx.Abort()
		return
	}

	resp, respError := h.uc.ITLogin(request)
	if respError != nil {
		response.GenerateResponse(ctx, respError.Code, response.WithMessage(respError.Error.Error()))
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 200, response.WithData(*resp))
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

	tokenData := jwt.TokenData{
		ID:   nurse.ID,
		Name: nurse.Name,
		Role: string(nurse.Role),
	}

	token, er := jwt.GenerateToken(tokenData)
	if er != nil {
		response.GenerateResponse(ctx, http.StatusInternalServerError, response.WithMessage("Failed to generate token"))
		return
	}

	res := FormatLoginResponse(nurse, token)

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

func (h *userHandler) GetUsers(ctx *gin.Context) {
	var queryParam UserQueryParams
	if err := ctx.ShouldBindQuery(&queryParam); err != nil {
		validatorMessage := validation.GenerateStructValidationError(err)
		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
		return
	}

	users, err := h.uc.GetUsers(queryParam)
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		return
	}

	res := FormatUsersResponse(users)
	response.GenerateResponse(ctx, http.StatusOK, response.WithMessage("success"), response.WithData(res))
}

func (h *userHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.uc.Delete(id)
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		return
	}

	response.GenerateResponse(ctx, 200)
}
