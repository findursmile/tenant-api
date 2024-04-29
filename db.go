package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/surrealdb/surrealdb.go"
)

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

func ImportSchema() {
    reader, _ := os.Open("schema.surql")

    client, req := getClient("import")

    req.Body = reader

    _, err := client.Do(req)

    if err != nil {
        fmt.Println("Unable to import schema")
        os.Exit(1)
    }
}

func getClient(uri string) (*http.Client, *http.Request) {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
    }

    client := &http.Client{Transport: tr}

    schema := "http"

    if os.Getenv("DB_SECURED") == "true" {
        schema = "https"
    }

    url := fmt.Sprintf("%s://%s:%s/%s", schema, os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), uri)
    fmt.Println("DB Endpoint ", url)
    req, _ := http.NewRequest("POST", url, nil)

    req.Header.Add("NS", os.Getenv("DB_NAMESPACE"))
    req.Header.Add("DB", os.Getenv("DB_DATABASE"))
    req.Header.Add("Accept", "application/json")
    req.SetBasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASS"))

    return client, req
}


