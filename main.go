package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/hashicorp/go-envparse"
	"github.com/surrealdb/surrealdb.go"
)

var DB *surrealdb.DB

func main() {
    loadEnv()

	r := SetupRoutes()

    InitDB()

	r.Run() // listen and serve on 0.0.0.0:8080
    defer Close()
}

func loadEnv() {
    buf, err := os.ReadFile(".env")

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    env, err := envparse.Parse(bytes.NewReader(buf))

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    for key, value := range env {
        os.Setenv(key, value)
    }
}
