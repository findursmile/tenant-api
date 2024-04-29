package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvaliSingupRequest(t *testing.T) {
    router := SetupRoutes()

    w := httptest.NewRecorder()

    req, _ := http.NewRequest("POST", "/api/signup", nil)

    router.ServeHTTP(w, req)

    var response map[string]string
    json.Unmarshal([]byte(w.Body.String()), &response)
    fmt.Print(response)
    assert.Equal(t, 412, w.Code)
    assert.Equal(t, "invalid request", response["exception"])
    assert.Equal(t, "Unable to parse request", response["message"])
}
