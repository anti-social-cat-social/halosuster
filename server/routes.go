package server

import (
	"halosuster/internal/record"
	"halosuster/internal/user"
	"halosuster/internal/image"
	"halosuster/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRoute(engine *gin.Engine, db *sqlx.DB) {
	// Handle for not found routes
	engine.NoRoute(NoRouteHandler)
	router := engine.Group("v1")

	router.GET("ping", pingHandler)

	initializeUserHandler(db, router)
	initializeImageHandler(router)
	initializeRecordHandler(db, router)
}

func initializeUserHandler(db *sqlx.DB, router *gin.RouterGroup) {
	// Initialize all necessary dependecies
	userRepo := user.NewUserRepository(db)
	userUc := user.NewUserUsecase(userRepo)
	userH := user.NewUserHandler(userUc)

	userH.Router(router)
}

func initializeImageHandler(router *gin.RouterGroup) {
	imageH := image.NewImageHandler()

	imageH.Router(router)
}

func initializeRecordHandler(db *sqlx.DB, router *gin.RouterGroup) {
	// Initalize all dependecies
	recordRepo := record.NewRecordRepo(db)
	recordUsecase := record.NewRecordUsecase(recordRepo)
	recordHandler := record.NewRecordHandler(recordUsecase)

	recordHandler.Router(router)
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
