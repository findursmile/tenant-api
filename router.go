package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)

var ApiRouter *gin.RouterGroup;

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

type Image struct {
    ImageUri string `json:"image_uri"`
    Status string `json:"status"`
}

func GetEvents(c *gin.Context) {
    event := map[string]interface{} {
        "cover_photo": "test",
        "event_date": time.Now(),
        "event_end_at": time.Now().Add(time.Hour),
        "name": "New event",
        "status": "draft",
        "title": "My new event",
    }

    _, err := DB.Create("event", event)
    if err != nil {
        panic(err)
    }

    image := Image{
        ImageUri: "TEst",
        Status: "pending",
    }

    _, err = DB.Create("image", image)

    if err != nil {
        panic(err)
    }

    data, err := DB.Select("image")

    var userEvents []Image;

    surrealdb.Unmarshal(data, &userEvents)

    if err != nil {
        c.AbortWithStatusJSON(412, gin.H{"message": "Unable to fetch events"})
    }

    c.JSON(200, gin.H{"events": userEvents})
}
