package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

var ApiRouter *gin.RouterGroup

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	DefineRoutes(r)

	return r
}

func DefineRoutes(router *gin.Engine) {
	router.POST("api/signin", Signin)
	router.POST("api/signup", Signup)

	ApiRouter = router.Group("api")
	ApiRouter.Use(Authendicate)

	RegisterEventRoutes()
	RegisterImageRoutes()
}

func Authendicate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	token, _ = strings.CutPrefix(token, "Bearer ")

	_, err := DB.Authenticate(token)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthendicated"})
	}

	tenant := GetTenant()
	c.Set("tenant", tenant)

	c.Next()
}
