package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)

func RegisterEventRoutes(route *gin.Engine) {
    ApiRouter.GET("events", GetEvents)
    ApiRouter.POST("events", CreateEvent)
    route.GET("api/events/:eventId", GetEvent)
    ApiRouter.POST("events/:eventId", UpdateEvent)
    ApiRouter.DELETE("events/:eventId", DeleteEvent)
    ApiRouter.PUT("events/:eventId/publish", PublishEvent)
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

func GetEvent(c *gin.Context) {
    data, err := DB.Select(c.Param("eventId"))

    if err != nil {
        c.AbortWithStatusJSON(404, gin.H{"message": "Event was not found"})
        return
    }

    var results Event

    surrealdb.Unmarshal(data, &results)

    c.JSON(200, gin.H{"event": results})
}

func CreateEvent(c *gin.Context) {
    var payload CreateEventPayload;

    if err := c.ShouldBind(&payload); err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    payload.Status = "pending"

    authUser, exists := c.Get("tenant")
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

    data, err := DB.Query(`CREATE event SET
        title = $title,
        name=$name,
        event_date = <datetime>$event_date,
        event_end_at=<datetime>$event_end_at,
        status=$status,
        tenant=$tenant;`, &payload)

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
    eventId := result[0].Result[0].Id
    DB.Let("event", eventId)

    handleCoverPhoto(c, eventId)

    c.JSON(200, gin.H{"message": "Event was created successfully", "event": eventId})
}

func UpdateEvent(c *gin.Context) {
    var payload CreateEventPayload;

    if err := c.ShouldBind(&payload); err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    eventId := c.Param("eventId")

    DB.Let("event", eventId)

    data, err := DB.Select(eventId)

    if err != nil  || data == nil {
        c.JSON(412, gin.H{"message": "Unable to select event"})
        return
    }

    event := new(Event)

    if err = surrealdb.Unmarshal(data, &event); err != nil {
        c.JSON(412, gin.H{"message": "Unable to find event", "exception": err.Error()})
        return
    }

    if tenant, ok := c.Get("tenant"); ok {
        payload.Tenant = tenant.(*Tenant).Id
    }

    sql := `UPDATE $event set
        title = $title,
        name = $name,
        event_date = <datetime>$event_date,
        event_end_at = <datetime>$event_end_at,
        tenant = $tenant;`

    if data, err = DB.Query(sql, &payload); err != nil {
        c.JSON(412, gin.H{"message": "Unable to update event -> " + event.Id, "exception": err.Error()})
        return
    }

    handleCoverPhoto(c, eventId)

    c.JSON(200, gin.H{"message": "Event was updated successfully", "event": event.Id})
}

func DeleteEvent(c *gin.Context) {
    eventId := c.Param("eventId")

    payload := map[string]string{
        "eventId": eventId,
    }

    if _, err := DB.Query(`UPDATE $eventId SET status="deleted"`, &payload); err != nil {
        c.JSON(412, gin.H{"message": "Unable to delete event"})
    }

    c.JSON(200, gin.H{"message": "Event was deleted"})
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", msg, err)
  }
}

func PublishEvent(c *gin.Context) {
    eventId := c.Param("eventId")

    payload := map[string]string{
        "eventId": eventId,
    }

    if _, err := DB.Query(`UPDATE $eventId SET status="publish"`, &payload); err != nil {
        c.JSON(412, gin.H{"message": "Unable to publish event"})
    }

    go PublishEventMessage(eventId)

    c.JSON(200, gin.H{"message": "Event was published"})
}

func handleCoverPhoto(c *gin.Context, eventId string) {
    file, err := c.FormFile("cover_photo")

    if err == nil {
        relativePath := GetEventImageDir(&eventId) + "/cover_" + file.Filename
        path, err := filepath.Abs(relativePath)
        if err == nil {
            c.SaveUploadedFile(file, path)
        }

        DB.Query("UPDATE $event SET cover_photo=$path", map[string]string{
            "path": relativePath,
        })
    }
}
