package middleware

import (
	"halosuster/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Check if user has a role from the token.
// Use this on router as middleware
func HasRoles(r string) gin.HandlerFunc {
	// Check from user HasRoles
	return func(ctx *gin.Context) {
		if ctx.GetString("userRole") != string(r) {
			response.GenerateResponse(ctx, http.StatusUnauthorized, response.WithMessage("User not permitted to do this operation"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
