package main

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var ApiRouter *gin.RouterGroup;

func DefineRoutes(router *gin.Engine) {
    router.POST("api/signin", login)
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

type LoginPayload struct {
    Email string `json:"email"`
    Password string `json:"password"`
    NS string `json:"NS"`
    DB string `json:"DB"`
    SC string `json:"SC"`
}

func login(c *gin.Context) {
    var loginPayload LoginPayload;

    err := c.ShouldBindJSON(&loginPayload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    loginPayload.NS = os.Getenv("DB_NAMESPACE")
    loginPayload.DB = os.Getenv("DB_DATABASE")
    loginPayload.SC = "tenant"

    token, err := DB.Signin(loginPayload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to login", "exception": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Logged in successfully", "token": token})
}
