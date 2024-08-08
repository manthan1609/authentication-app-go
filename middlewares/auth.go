package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manthan1609/authentication-app/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// clientToken := c.Request.Header.Get("token")
		clientToken, err := c.Cookie("access_token")

		if err != nil || clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token not present",
			})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)

		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.UserType)

		c.Next()
	}

}
