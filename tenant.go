package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)

type TenantPayload struct {
    Name string `json:"name"`
    Mobile string `json:"mobile"`
    Country_code string `json:"country_code"`
    Email string `json:"email"`
    Password string `json:"password"`
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
        panic(err)
    }
    fmt.Print(payload)
    data, err := DB.Create("tenant", map[string]interface{} {
        "name": payload.Name,
        "mobile": payload.Mobile,
        "email": payload.Email,
        "password": payload.Password,
        "country_code": payload.Country_code,
    })

    if err != nil {
        panic(err)
    }

    tenant := make([]Tenant, 1)

    err = surrealdb.Unmarshal(data, &tenant)

    if err != nil {
        panic(err)
    }


    fmt.Printf("Tenant :%s", tenant[0].Id)

    c.JSON(200, gin.H{"message": "Tenant created successfully"})
}
