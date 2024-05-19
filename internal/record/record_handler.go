package record

import (
	"halosuster/internal/middleware"
	"halosuster/pkg/response"
	"halosuster/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type recordHandler struct {
	uc IRecordUsecase
}

func NewRecordHandler(uc IRecordUsecase) *recordHandler {
	return &recordHandler{
		uc: uc,
	}
}

func (h *recordHandler) Router(r *gin.RouterGroup) {
	group := r.Group("medical/record", middleware.UseJwtAuth)

	group.GET("/", h.GetAll)
}

func (h *recordHandler) GetAll(ctx *gin.Context) {
	var params RecordQueryParam

	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.GenerateResponse(ctx, 500, response.WithMessage(err.Error()))
		ctx.Abort()
		return
	}

	// Validate params
	validate := validator.New()
	errVal := validate.Struct(params)

	if errVal != nil {
		msgVal := validation.FormatValidation(errVal)

		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Request invalid"), response.WithData(msgVal))
		ctx.Abort()
		return
	}

	result, err := h.uc.GetAll(params)
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 200, response.WithData(result))
}
