package Interfaces

import "github.com/gin-gonic/gin"

type IAuthenticationMiddleware interface {
	RequireAuth(IToken IToken) gin.HandlerFunc
}
