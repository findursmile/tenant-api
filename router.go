package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

var ApiRouter *gin.RouterGroup;

func DefineRoutes(router *gin.Engine) {
    router.POST("api/singin", login)
    router.POST("api/signup", CreateTenant)

    ApiRouter = router.Group("api")
    ApiRouter.Use(Authendicate)
}

func Authendicate(c *gin.Context) {
    token := c.GetHeader("Authorization")
    token, _ = strings.CutPrefix(token, "Bearer ")

    _, err := DB.Authenticate(token)

    if (err != nil) {
        c.AbortWithStatusJSON(401, gin.H{"message": "Unauthendicated"})
    }

    c.Next()
}


func login(c *gin.Context) {
}
