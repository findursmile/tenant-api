package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(M *testing.M) {
    os.Setenv("DB_NAMESPACE", "unit_test")
    loadEnv()
    InitDB()
    ImportSchema()
    router = SetupRoutes()

    exitCode := M.Run()

    RemoveSchema()
    os.Exit(exitCode)
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

func RemoveSchema() {
    client, req := getClient("sql")

    req.Body = io.NopCloser(strings.NewReader(fmt.Sprint("REMOVE NAMESPACE IF EXISTS ", os.Getenv("DB_NAMESPACE"))))

    _, err := client.Do(req)

    if err != nil {
        fmt.Println("Unable to remove schema")
        os.Exit(1)
    }
}

func TestInvaliSingupRequest(t *testing.T) {
    w := httptest.NewRecorder()

    req, _ := http.NewRequest("POST", "/api/signup", nil)

    router.ServeHTTP(w, req)

    var response map[string]string
    json.Unmarshal([]byte(w.Body.String()), &response)

    assert.Equal(t, 412, w.Code)
    assert.Equal(t, "invalid request", response["exception"])
    assert.Equal(t, "Unable to parse request", response["message"])
}

func TestSuccessfulSingup(t *testing.T) {

    w := httptest.NewRecorder()

    payload := &SignupPayload{
        Email: "praem1990@gmail.com",
        Password: "asdf",
        Country_code: "IN",
        Mobile: "0909090",
    }

    s, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/api/signup", bytes.NewReader(s))

    router.ServeHTTP(w, req)

    var response map[string]string
    json.Unmarshal([]byte(w.Body.String()), &response)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "Tenant created successfully", response["message"])
}

func TestLogin(t *testing.T) {
    w := httptest.NewRecorder()

    payload := &SigninPayload{
        Email: "praem1990@gmail.com",
        Password: "asdf",
    }

    s, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/api/signin", bytes.NewReader(s))

    router.ServeHTTP(w, req)

    var response map[string]string
    json.Unmarshal([]byte(w.Body.String()), &response)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "Logged in successfully", response["message"])

}
