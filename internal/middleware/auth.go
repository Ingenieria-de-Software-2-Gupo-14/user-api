package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/go-core/pkg/log"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT token authentication
func AuthMiddleware(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := getAuthToken(ctx)
		claims, err := models.ParseToken(tokenStr)
		if err != nil {
			utils.ErrorResponseWithErr(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		uId, _ := strconv.Atoi(claims.Subject)
		blocked, err := userService.IsUserBlocked(ctx.Request.Context(), uId)
		if err != nil {
			utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		if blocked {
			utils.ErrorResponse(ctx, http.StatusForbidden, "User is blocked")
			// Make Cookie Expire
			ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
			ctx.Abort()
			return
		}

		// Refresh token if it is about to expire
		if claims.ExpiresAt < time.Now().Add(time.Minute*5).Unix() {
			newToken, err := models.GenerateToken(uId, claims.Email, claims.Name, claims.Role)
			if err != nil {
				log.Error(ctx, "Error generating token", "error", err.Error())
			} else {
				ctx.SetCookie("Authorization", newToken, 3600, "/", "", false, true)
			}
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

// getAuthToken extracts the authentication token from either cookie or Authorization header
func getAuthToken(c *gin.Context) string {
	auth, _ := c.Cookie("Authorization")
	if auth == "" {
		if parts := strings.Fields(c.GetHeader("Authorization")); len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			auth = parts[1]
		}
	}
	return auth
}

// AdminOnlyMiddleware ensures only users with admin role can access the route
func AdminOnlyMiddleware(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// First authenticate the user
		tokenStr := getAuthToken(ctx)
		claims, err := models.ParseToken(tokenStr)
		if err != nil {
			utils.ErrorResponseWithErr(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		uId, _ := strconv.Atoi(claims.Subject)

		// Check if user has admin role
		if claims.Role != "admin" {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Admin access required")
			ctx.Abort()
			return
		}

		// Refresh token if it is about to expire
		if claims.ExpiresAt < time.Now().Add(time.Minute*5).Unix() {
			newToken, err := models.GenerateToken(uId, claims.Email, claims.Name, claims.Role)
			if err != nil {
				log.Error(ctx, "Error generating token", "error", err.Error())
			} else {
				ctx.SetCookie("Authorization", newToken, 3600, "/", "", false, true)
			}
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

// UserOrAdminMiddleware allows access to users who own the resource (ID matches path parameter) or admin users
func UserOrAdminMiddleware(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// First authenticate the user
		tokenStr := getAuthToken(ctx)
		claims, err := models.ParseToken(tokenStr)
		if err != nil {
			utils.ErrorResponseWithErr(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		uId, _ := strconv.Atoi(claims.Subject)
		blocked, err := userService.IsUserBlocked(ctx.Request.Context(), uId)
		if err != nil {
			utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		if blocked {
			utils.ErrorResponse(ctx, http.StatusForbidden, "User is blocked")
			ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
			ctx.Abort()
			return
		}

		// Get the ID from path parameter
		pathIdStr := ctx.Param("id")
		if pathIdStr == "" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Missing user ID in path")
			ctx.Abort()
			return
		}

		pathId, err := strconv.Atoi(pathIdStr)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID format")
			ctx.Abort()
			return
		}

		// Check if user is admin or owns the resource
		if claims.Role != "admin" && uId != pathId {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied: you can only access your own resources")
			ctx.Abort()
			return
		}

		// Refresh token if it is about to expire
		if claims.ExpiresAt < time.Now().Add(time.Minute*5).Unix() {
			newToken, err := models.GenerateToken(uId, claims.Email, claims.Name, claims.Role)
			if err != nil {
				log.Error(ctx, "Error generating token", "error", err.Error())
			} else {
				ctx.SetCookie("Authorization", newToken, 3600, "/", "", false, true)
			}
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
