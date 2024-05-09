package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)

func RegisterImageRoutes() {
    ApiRouter.GET("events/:eventId/images", GetImages)
    ApiRouter.POST("events/:eventId/images", UploadImages)
    ApiRouter.DELETE("events/:eventId/images/:imageId", DeleteImage)
}

type ImageFilter struct {
    PageNo int `json:"page,default=1" form:"page,default=1"`
    Limit int `json:"limit,default=25" form:"limit,default=25"`
    Start int `json:"start,default=0" form:"start,default=0"`
    EventId string `json:"event_id" form:"event_id"`
}

type Image struct {
    Id string `json:"id"`
    ImageUri string `json:"image_uri"`
    Event string `json:"event"`
    Status string `json:"status"`
    Created time.Time `json:"created"`
}

func GetImages(c *gin.Context) {
    filter := new(ImageFilter)

    if err := c.ShouldBind(&filter); err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    c.ShouldBind(&filter)

    filter.Start = (filter.PageNo - 1) * filter.Limit
    filter.EventId = c.Param("eventId")

    sql := `SELECT * from image where event = $event_id LIMIT $limit START $start`
    data, err := DB.Query(sql, &filter)
    fmt.Println(filter)

    if err != nil {
        c.JSON(412, gin.H{"message": "Unable to parse request", "exception": err.Error()})
        return
    }

    type res struct {
        Result []Image `json:"result"`
        Status string `json:"status"`
        Time string `json:"time"`
    }

    results := make([]res, 1)

    surrealdb.Unmarshal(data, &results)

    c.JSON(200, gin.H{"images": results[0].Result})
}

func UploadImages(c *gin.Context) {
    eventId := c.Param("eventId")

    form, err := c.MultipartForm()

    if err != nil {
        c.JSON(405, gin.H{"message": "Invalid request"})
        return
    }

    images := form.File["images[]"]

    count := 0

    for _, image := range images {
        uri := GetEventImageDir(&eventId)
        path, _ := filepath.Abs(uri)

        if err = c.SaveUploadedFile(image, path); err != nil {
            continue
        }

        sql := `CREATE image SET
            image_uri = $uri,
            status="pending",
            event=$event;
        `
        payload := map[string]string {
            "uri": uri,
            "event": eventId,
        }

        if res, err := DB.Query(sql, &payload); err == nil {
            fmt.Println(res)
            count++
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to creat image"})
            return
        }

    }

    if count == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"message": "No images were valid and uploaded"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d image(s) were uploaded", count)})
}

func DeleteImage(c *gin.Context) {
    payload := map[string]string {
        "eventId": c.Param("eventId"),
        "imageId": c.Param("imageId"),
    }

    data, err := DB.Select(payload["imageId"])

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Image not found"})
        return
    }

    image := make(map[string]string)

    if err := surrealdb.Unmarshal(data, &image); err != nil || image[ "event" ] != payload["eventId"] {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Image not found"})
        return
    }

    if err = os.Remove(image["image_uri"]); err != nil  {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to delete the image"})
        return
    }

    sql := `DELETE $image where event = $event`

    if _, err := DB.Query(sql, &payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Image not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Image was deleted"})
}

func GetEventImageDir(eventId *string) string {
    return "data/images/" + *eventId
}
