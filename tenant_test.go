package main

import (
	"bytes"
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
    req.Header.Add("Content-Type", "application/json")

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
        Name: "Mohan Raj",
        Email: "praem19902@gmail.com",
        Password: "asdf",
        CountryCode: "IN",
        Mobile: "0909090",
    }

    s, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/api/signup", bytes.NewReader(s))
    req.Header.Add("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    var response map[string]string
    json.Unmarshal([]byte(w.Body.String()), &response)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "", response["exception"])
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
    req.Header.Add("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    var response map[string]string
    json.Unmarshal([]byte(w.Body.String()), &response)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "Logged in successfully", response["message"])

}
