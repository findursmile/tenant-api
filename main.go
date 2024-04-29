package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-envparse"
	"github.com/surrealdb/surrealdb.go"
)

var DB *surrealdb.DB

func main() {
    loadEnv()

	r := SetupRoutes()

    InitDB()

	r.Run() // listen and serve on 0.0.0.0:8080
}

func InitDB() {
    websocket.DefaultDialer.TLSClientConfig = &tls.Config{
        InsecureSkipVerify: true,
    }

    var err error
    schema := "ws"
    if os.Getenv("DB_SECURED") == "true" {
        schema = "wss"
    }
    endpoint := fmt.Sprint(schema, "://", os.Getenv( "DB_HOST" ), ":", os.Getenv( "DB_PORT" ), "/rpc")
    fmt.Println("WS endpoint: ", endpoint)
    DB, err = surrealdb.New(endpoint)

    if err != nil {
        panic(err)
    }

    DB.Use(os.Getenv( "DB_NAMESPACE" ), os.Getenv( "DB_DATABASE" ))
}

func loadEnv() {
    buf, err := os.ReadFile(".env")

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    env, err := envparse.Parse(bytes.NewReader(buf))

    if err != nil {
        fmt.Println(string(buf[:]))
        fmt.Println(err.Error())
        return
    }

    for key, value := range env {
        fmt.Println("Setting ENV ", key, ": ", value)
        os.Setenv(key, value)
    }
}
