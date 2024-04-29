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
    loadEnv()

	r := gin.Default()
    ApiRouter = r.Group("/api")
    ApiRouter.Use(Authendicate)

    r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

    websocket.DefaultDialer.TLSClientConfig = &tls.Config{
        InsecureSkipVerify: true,
    }

    endpoint := fmt.Sprint("wss://", os.Getenv( "DB_HOST" ), ":", os.Getenv( "DB_PORT" ), "/rpc")
    var err error
    DB, err = surrealdb.New(endpoint)

    if err != nil {
        panic(err)
    }

    DB.Use(os.Getenv( "DB_NAMESPACE" ), os.Getenv( "DB_DATABASE" ))

    DefineRoutes(r)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func loadEnv() {
    buf, err := os.ReadFile(".env")

    if err != nil {
        return
    }

    env, err := envparse.Parse(bytes.NewReader(buf))

    if err != nil {
        return
    }

    for key, value := range env {
        os.Setenv(key, value)
    }
}
