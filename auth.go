package main

import "github.com/gin-gonic/gin"

func RegisterRoutes() {
    auth := ApiRouter.Group("auth")
    auth.POST("singin", login)
    auth.POST("signup", CreateTenant)
}

func login(c *gin.Context) {
}
