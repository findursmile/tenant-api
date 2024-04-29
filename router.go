package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)

var ApiRouter *gin.RouterGroup;

func SetupRoutes() *gin.Engine {
    r := gin.Default()

    DefineRoutes(r)

    return r
}

func DefineRoutes(router *gin.Engine) {
    router.POST("api/signin", Signin)
    router.POST("api/signup", Signup)

    ApiRouter = router.Group("api")
    ApiRouter.Use(Authendicate)
    ApiRouter.GET("events", GetEvents)
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

func GetEvents(c *gin.Context) {
    data, err := DB.Select("event")

    var userEvents []map[string]interface{};

    surrealdb.Unmarshal(data, &userEvents)

    if err != nil {
        c.AbortWithStatusJSON(412, gin.H{"message": "Unable to fetch events"})
    }

    c.JSON(200, gin.H{"events": userEvents})
}
