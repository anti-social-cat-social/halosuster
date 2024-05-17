package server

import (
	"eniqlo/internal/customer"
	"eniqlo/internal/product"
	"eniqlo/internal/staff"
	"eniqlo/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRoute(engine *gin.Engine, db *sqlx.DB) {
	// Handle for not found routes
	engine.NoRoute(NoRouteHandler)
	router := engine.Group("v1")

	router.GET("ping", pingHandler)

	initializeStaffHandler(db, router)
	initializeProductHandler(db, router)
	initializeCustomerHandler(db, router)
}

func initializeStaffHandler(db *sqlx.DB, router *gin.RouterGroup) {
	// Initialize all necessary dependecies
	staffRepo := staff.NewStaffRepository(db)
	staffUc := staff.NewStaffUsecase(staffRepo)
	staffH := staff.NewStaffHandler(staffUc)

	staffH.Router(router)
}

func initializeProductHandler(db *sqlx.DB, router *gin.RouterGroup) {
	// Initialize all necessary dependecies
	productRepo := product.NewProductRepository(db)
	productUc := product.NewProductUsecase(productRepo)
	productH := product.NewProductHandler(productUc)

	productH.Router(router)
}

func initializeCustomerHandler(db *sqlx.DB, router *gin.RouterGroup) {
	// Initialize all necessary dependecies
	customerRepo := customer.NewCustomerRepository(db)
	customerUc := customer.NewCustomerUsecase(customerRepo)
	customerH := customer.NewCustomerHandler(customerUc)

	customerH.Router(router)
}

func NoRouteHandler(ctx *gin.Context) {
	response.GenerateResponse(ctx, http.StatusNotFound, response.WithMessage("Page not found"))
}

// Handler for ping request from routes
func pingHandler(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		struct {
			Data    any    `json:"data"`
			Message string `json:"message"`
			Success bool   `json:"success"`
		}{
			Success: true,
			Message: "Server is online",
			Data:    true,
		},
	)
}
