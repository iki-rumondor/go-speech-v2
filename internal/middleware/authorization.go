package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/iki-rumondor/go-speech/internal/utils"
)

func IsRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		mc := c.MustGet("map_claims")
		mapClaims := mc.(jwt.MapClaims)

		roleJwt := mapClaims["role"].(string)
		if roleJwt != role {
			utils.HandleError(c, response.UNAUTH_ERR("Hak Akses Dibatasi"))
			return
		}
		c.Next()
	}
}
