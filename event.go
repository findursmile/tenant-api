package main

import (
	"fmt"
	"path/filepath"
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
        fmt.Println("Unable to parse date")
        return err
    }
    *j = JsonDate(t)
    return nil
}

func (j JsonDate) MarshalJSON() ([]byte, error) {
    t := time.Time(j)
    var stamp = fmt.Sprintf(`"%s"`, t.Format("2006-01-01"))
    return []byte(stamp), nil
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
    Name string `form:"name" json:"name" binding:"required"`
    Title string `form:"title" json:"title" binding:"required"`
    // CoverPhoto string `form:"cover_photo" json:"cover_photo"`
    EventDate *JsonDate `form:"event_date" json:"event_date" binding:"required"`
    EventEndAt *JsonDate `form:"event_end_at" json:"event_end_at" binding:"required"`
    Tenant string `form:"tenant" json:"tenant"`
    Status string `form:"status" json:"status"`
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

    if err := c.ShouldBind(&payload); err != nil {
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

    file, err := c.FormFile("cover_photo")

    if err == nil {
        path, err := filepath.Abs("data/images/" + file.Filename)
        if err == nil {
            c.SaveUploadedFile(file, path)
        }
    }

    // data, err := DB.Create("event", &payload)
    data, err := DB.Query(`CREATE event SET
        title = $title,
        name=$name,
        event_date = <datetime>$event_date,
        event_end_at=<datetime>$event_end_at,
        status=$status,
        tenant=$tenant;`, payload)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to create event", "exception": err.Error()})
        return
    }

    type res struct {
        Result []Event `json:"result"`
        Status string `json:"status"`
        Time string `json:"time"`
    }

    result := make([]res, 1)

    err = surrealdb.Unmarshal(data, &result)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to Unmarshal event", "exception": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Event was created successfully", "event": result[0].Result[0].Id})
}

func UpdateEvent(c *gin.Context) {
    var payload CreateEventPayload;

    if err := c.ShouldBind(&payload); err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    data, err := DB.Select(c.Param("id"))

    if err != nil  || data == nil {
        c.JSON(412, gin.H{"message": "Unable to select event"})
        return
    }

    event := new(Event)

    if err = surrealdb.Unmarshal(data, &event); err != nil {
        c.JSON(412, gin.H{"message": "Unable to find event", "exception": err.Error()})
        return
    }

    payload.Tenant = GetTenant().Id
    payload.Status = event.Status

    if data, err = DB.Update(event.Id, &payload); err != nil {
        c.JSON(412, gin.H{"message": "Unable to update event -> " + event.Id, "exception": err.Error()})
        return
    }

    if err = surrealdb.Unmarshal(data, &event); err != nil {
        c.JSON(412, gin.H{"message": "Unable to Unmarshal event", "exception": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "Event was updated successfully", "event": event.Id})
}
