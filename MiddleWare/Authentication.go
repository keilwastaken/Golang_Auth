package MiddleWare

import (
	"Clarity_go/Interfaces"
	"github.com/gin-gonic/gin"
)

type Authentication struct {
}

var _ Interfaces.IAuthenticationMiddleware = (*Authentication)(nil)

func (auth Authentication) RequireAuth(IToken Interfaces.IToken) gin.HandlerFunc {
	return func(c *gin.Context) {

		user, err := IToken.ValidateToken(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(err.Status, gin.H{"error": err.Message})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
