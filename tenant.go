package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type TenantPayload struct {
    Name string `json:"name"`
    Mobile string `json:"mobile"`
    Country_code string `json:"country_code"`
    Email string `json:"email"`
    Password string `json:"password"`
    Status string `json:"status"`
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

func CreateTenant(c *gin.Context) {
    var payload TenantPayload;

    err := c.ShouldBindJSON(&payload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    payload.Status = "active"

    _, err = DB.Create("tenant", payload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to create Tenant", "exception": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Tenant created successfully"})
}
