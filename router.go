package main

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ApiRouter *gin.RouterGroup

func SetupRoutes() *gin.Engine {
	r := gin.Default()
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AddAllowHeaders("Authorization")

    r.Use(cors.New(config))

	DefineRoutes(r)

	return r
}

func DefineRoutes(router *gin.Engine) {
	router.POST("api/signin", Signin)
	router.POST("api/signup", Signup)

	ApiRouter = router.Group("api")
	ApiRouter.Use(Authendicate)

	RegisterEventRoutes(router)
	RegisterImageRoutes(router)

    router.Static("/data", "./data")
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
