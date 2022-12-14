package middleware

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	jwtlib "github.com/arvians-id/go-apriori-microservice/adapter/third-party/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func AuthJwtMiddleware(configuration *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			c.JSON(http.StatusUnauthorized, response.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: "invalid token",
				Data:   nil,
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		jwtLibrary := jwtlib.NewJsonWebToken(configuration)
		token, err := jwtLibrary.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: err.Error(),
				Data:   nil,
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, response.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: err.Error(),
				Data:   nil,
			})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("id_user", claims["id_user"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
