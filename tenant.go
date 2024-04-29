package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type SignupPayload struct {
    Name string `json:"name"`
    Mobile string `json:"mobile"`
    Country_code string `json:"country_code"`
    Email string `json:"email"`
    Password string `json:"password"`
    Status string `json:"status"`
    NS string `json:"NS"`
    DB string `json:"DB"`
    SC string `json:"SC"`
}

type SigninPayload struct {
    Email string `json:"email"`
    Password string `json:"password"`
    NS string `json:"NS"`
    DB string `json:"DB"`
    SC string `json:"SC"`
}

type Tenant struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Mobile string `json:"mobile"`
    Country_code string `json:"country_code"`
    Email string `json:"email"`
    Password string `json:"password"`
    Status string `json:"status"`
    Created time.Time `json:"created"`
    Updated time.Time `json:"updated"`
}

func Signup(c *gin.Context) {
    var payload SignupPayload;

    err := c.ShouldBindJSON(&payload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    payload.Status = "active"
    payload.SC = "tenant"
    payload.NS = os.Getenv("DB_NAMESPACE")
    payload.DB = os.Getenv("DB_DATABASE")

    _, err = DB.Signup(payload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to create Tenant", "exception": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Tenant created successfully"})
}

func Signin(c *gin.Context) {
    var signinPayload SigninPayload;

    err := c.ShouldBindJSON(&signinPayload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    signinPayload.NS = os.Getenv("DB_NAMESPACE")
    signinPayload.DB = os.Getenv("DB_DATABASE")
    signinPayload.SC = "tenant"

    token, err := DB.Signin(signinPayload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to login", "exception": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Logged in successfully", "token": token})
}
