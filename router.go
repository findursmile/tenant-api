package main

import "github.com/gin-gonic/gin"

var ApiRouter *gin.RouterGroup;

func DefineRoutes(router *gin.Engine) {
    ApiRouter = router.Group("api")
    ApiRouter.POST("login")
}


