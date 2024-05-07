package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)

func RegisterEventRoutes() {
    ApiRouter.GET("events", GetEvents)
    ApiRouter.POST("events", CreateEvent)
    ApiRouter.POST("events/:id", UpdateEvent)
}

type JsonDate time.Time

// Implement Marshaler and Unmarshaler interface
func (j *JsonDate) UnmarshalJSON(b []byte) error {
    s := strings.Trim(string(b), "\"")
    t, err := time.Parse("2006-01-02", s)
    if err != nil {
        return err
    }
    *j = JsonDate(t)
    return nil
}

func (j JsonDate) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Time(j))
}

// Maybe a Format function for printing your date
func (j JsonDate) Format(s string) string {
    t := time.Time(j)
    return t.Format(s)
}

type Event struct {
    Id string `json:"id"`
    Name string `json:"name"`
    CoverPhoto string `json:"cover_photo"`
    Title string `json:"title"`
    EventDate *time.Time `json:"event_date"`
    EventEndAt *time.Time `json:"event_end_at"`

    Status string `json:"status"`

    Created *time.Time `json:"created"`
    Updated *time.Time `json:"updated"`
}

type CreateEventPayload struct {
    Name string `json:"name" binding: "required"`
    Title string `json:"title" binding: "required"`
    CoverPhoto string `json:"cover_photo"`
    EventDate *JsonDate `json:"event_date" binding: "required"`
    EventEndAt *JsonDate `json:"event_end_at" binding: "required"`
    Tenant string `json:"tenant"`
    Status string `json:"status"`
}

func GetEvents(c *gin.Context) {
    data, err := DB.Select("event")

    var userEvents []map[string]interface{};

    surrealdb.Unmarshal(data, &userEvents)

    if err != nil {
        c.AbortWithStatusJSON(412, gin.H{"message": "Unable to fetch events"})
    }

    c.JSON(200, gin.H{"events": userEvents})
}

func CreateEvent(c *gin.Context) {
    var payload CreateEventPayload;

    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    payload.Status = "pending"

    authUser, exists := c.Get("user")
    if exists == false {
        c.JSON(412, gin.H{"message": "Unable to get the user"})
        return
    }

    tenant, ok := authUser.(*Tenant)
    if ok == false {
        c.JSON(412, gin.H{"message": "Unable cast the user"})
        return
    }
    payload.Tenant = tenant.Id

    data, err := DB.Create("event", &payload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to create event", "exception": err.Error()})
        return
    }

    events := make([]Event, 1)

    err = surrealdb.Unmarshal(data, &events)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to Unmarshal event", "exception": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Event was created successfully", "event": events[0].Id})
}

func UpdateEvent(c *gin.Context) {
    var payload CreateEventPayload;

    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    events := make([]Event, 1)

    data, _ := DB.Select(c.Param("id"))

    err := surrealdb.Unmarshal(data, &events)

    if err != nil || len(events) == 0 {
        c.JSON(412, gin.H{"message": "Unable to find event", "exception": err.Error()})
        return
    }

    events[0].Name = payload.Name
    events[0].Title = payload.Title

    eDate, _ := time.Parse("2024-01-30", payload.EventDate.Format("2024-04-04"))
    endDate, _ := time.Parse("2024-01-30", payload.EventEndAt.Format("2024-04-04"))
    events[0].EventDate = &eDate
    events[0].EventEndAt = &endDate

    data, err = DB.Change(c.Param("id"), &events)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to update event", "exception": err.Error()})
        return
    }

    err = surrealdb.Unmarshal(data, &events)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to Unmarshal event", "exception": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Event was updated successfully", "event": events[0].Id})
}
