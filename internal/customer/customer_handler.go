package customer

import (
	"eniqlo/internal/middleware"
	"eniqlo/pkg/response"
	"eniqlo/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

type customerHandler struct {
	uc ICustomerUsecase
}

// Constructor for customer handler struct
func NewCustomerHandler(uc ICustomerUsecase) *customerHandler {
	return &customerHandler{
		uc: uc,
	}
}

// Router is required to wrap all user request by spesific path URL
func (h *customerHandler) Router(r *gin.RouterGroup) {
	// Grouping to give URL prefix
	group := r.Group("customer")

	group.Use(middleware.UseJwtAuth)

	// Utillize group to use global setting on group parent (if exists)
	group.GET("", h.FindAll)
	group.POST("register", h.register)
}

func (h *customerHandler) FindAll(c *gin.Context) {
	query := QueryParams{}
	if err := c.ShouldBindQuery(&query); err != nil {
		res := validation.FormatValidation(err)
		response.GenerateResponse(c, res.Code, response.WithMessage(res.Message))
		return
	}

	customers, err := h.uc.FindCustomers(query)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		return
	}

	res := FormatCustomersResponse(customers)

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("Customer fetched successfully!"), response.WithData(res))
}

func (h *customerHandler) register(ctx *gin.Context) {
	var request CustomerRegisterDTO

	// Parse request body to DTO
	// If error return error response
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.GenerateResponse(ctx, 400)
		ctx.Abort()
		return
	}

	// Validate request
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Register custom validation
	errValidation := validate.RegisterValidation("valid_name", validation.ValidNameValidator)
	if errValidation != nil {
		response.GenerateResponse(ctx, http.StatusInternalServerError, response.WithMessage("Failed register validation"))
	}

	// Generate error validation if not any field is not valid
	if err := validate.Struct(request); err != nil {
		validatorMessage := validation.GenerateStructValidationError(err)

		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
		ctx.Abort()
		return
	}

	// Process register via usecase
	result, err := h.uc.Register(request)
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 201, response.WithMessage("User registered successfully"), response.WithData(result))
}
