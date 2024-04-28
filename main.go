package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-envparse"
	"github.com/surrealdb/surrealdb.go"
)

var DB *surrealdb.DB

func main() {
	r := gin.Default()
    ApiRouter = r.Group("/api")

    r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

    buf, err := os.ReadFile(".env")

    if err != nil {
        panic(err)
    }

    env, err := envparse.Parse(bytes.NewReader(buf))

    if err != nil {
        panic(err)
    }

    websocket.DefaultDialer.TLSClientConfig = &tls.Config{
        InsecureSkipVerify: true,
    }

    endpoint := fmt.Sprint("wss://", env["DB_HOST"], ":", env["DB_PORT"], "/rpc")
    DB, err := surrealdb.New(endpoint)

    if err != nil {
        panic(err)
    }

    _, err = DB.Signin(map[string]interface{}{
        "user": env["DB_USER"],
        "pass": env["DB_PASS"],
    })

    if err != nil {
        panic(err)
    }

    RegisterRoutes()

	r.Run() // listen and serve on 0.0.0.0:8080
}
