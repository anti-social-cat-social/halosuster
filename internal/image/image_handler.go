package image

import (
	"halosuster/internal/middleware"
	"halosuster/pkg/response"
	"path/filepath"
	"log"
	"github.com/gin-gonic/gin"
	"sync"
)

var wg sync.WaitGroup

type imageHandler struct {
	// uc IUserUsecase
}

// Constructor for user handler struct
func NewImageHandler() *imageHandler {
	return &imageHandler{
		// uc: uc,
	}
	// uc IImageUsecase
}

func (h *imageHandler) Router(r *gin.RouterGroup) {
	group := r.Group("image")
	// group.MaxMultipartMemory = 2 << 20
	group.POST("", middleware.UseJwtAuth, middleware.HasRoles("it","nurse"), h.Upload)
}

func (h *imageHandler) Upload(ctx *gin.Context) {
	wg.Add(1)
    go func() {
        defer wg.Done()

		file, _ := ctx.FormFile("file")

		if file == nil {
			response.GenerateResponse(ctx, 400)
			ctx.Abort()
			return
		}

		if file.Size > 2000000 {
			response.GenerateResponse(ctx, 400)
			ctx.Abort()
			return
		}

		if file.Header["Content-Type"][0] != "image/jpeg" {
			response.GenerateResponse(ctx, 400)
			ctx.Abort()
			return
		}

		
		// log.Println(file.Filename)
		// filename := filepath.Base(file.Filename)
		// log.Println(filename)

        err := UploadFileToS3(file)
        if err != nil {
            log.Println("Error uploading file:", err)
            response.GenerateResponse(ctx, 500)
			ctx.Abort()
			return
        }

		response.GenerateResponse(ctx, 200)
    }()

    wg.Wait()

}