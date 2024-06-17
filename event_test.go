package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/surrealdb/surrealdb.go"
)

func GetToken() string {
    payload := &SignupPayload{
        Email: "praem1990@gmail.com",
        Password: "asdf",
        Country_code: "IN",
        Mobile: "0909090",
        Status : "active",
        SC : "tenant",
        NS : os.Getenv("DB_NAMESPACE"),
        DB : os.Getenv("DB_DATABASE"),
    }

    _, err := DB.Signup(payload)

    if err != nil {
        panic(err)
    }

    signInPayload := &SigninPayload{
        Email: "praem1990@gmail.com",
        Password: "asdf",
        SC : "tenant",
        NS : os.Getenv("DB_NAMESPACE"),
        DB : os.Getenv("DB_DATABASE"),
    }

    token, err := DB.Signin(&signInPayload)
    if err != nil {
        panic(err)
    }

    return token.(string)
}

func TestCreateEvent(t *testing.T) {
    w := httptest.NewRecorder()

    payload := &map[string]interface{} {
        "name": "Test Event",
        "title": "Event that you never missed",
        "event_date": "2024-05-05",
        "event_end_at": "2024-05-15",
    }

    s, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/api/events", bytes.NewReader(s))

    req.Header.Add("Authorization", "Bearer " + GetToken())

    router.ServeHTTP(w, req)

    var response map[string]string
    json.Unmarshal([]byte(w.Body.String()), &response)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "", response["exception"])
    assert.Equal(t, "Event was created successfully", response["message"])
}

func TestUpdateEvent(t *testing.T) {
    w := httptest.NewRecorder()

    payload := &map[string]interface{} {
        "name": "Test Event",
        "title": "Event that you never missed",
        "event_date": time.Now(),
        "event_end_at": time.Now().Add(time.Hour * 24),
        "status": "pending",
        "tenant": GetTenant().Id,
    }

    data, err := DB.Create("event", &payload)

    if err != nil {
        panic(err)
    }

    events := make([]Event, 1)

    err = surrealdb.Unmarshal(data, &events)

    if err != nil  {
        panic(err)
    }

    payload = &map[string]interface{} {
        "name": "Test Event 2",
        "title": "Event that you never missed",
    }

    s, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/api/events/" + events[0].Id, bytes.NewReader(s))

    req.Header.Add("Authorization", "Bearer " + GetToken())

    router.ServeHTTP(w, req)

    var response map[string]string
    json.Unmarshal([]byte(w.Body.String()), &response)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "", response["exception"])
    assert.Equal(t, "Event was updated successfully", response["message"])
}
