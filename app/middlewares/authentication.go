package middlewares

import (
	"adamnasrudin03/movie-festival/app/entity"
	"adamnasrudin03/movie-festival/pkg/database"
	"adamnasrudin03/movie-festival/pkg/helpers"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := helpers.VerifyToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, helpers.APIResponse(err.Error(), http.StatusUnauthorized, nil))
			c.Abort()
			return
		}

		c.Set("userData", claims)
		c.Next()
	}
}

func AuthorizationMustBeAdmin() gin.HandlerFunc {

	return func(c *gin.Context) {
		db := database.GetDB()
		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))
		userEmail := userData["email"].(string)
		user := entity.User{}

		err := db.Select("role").Where("id = ? AND email = ?", userID, userEmail).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) || user.ID == 0 {
			c.JSON(http.StatusUnauthorized, helpers.APIResponse("Log in again with registered user", http.StatusUnauthorized, nil))
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, helpers.APIResponse("Failed to check user log in", http.StatusUnauthorized, nil))
			c.Abort()
			return
		}

		if user.Role != "ADMIN" {
			c.JSON(http.StatusForbidden, helpers.APIResponse("You are not allowed to access this resources", http.StatusForbidden, nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
