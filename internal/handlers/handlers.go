package handlers

import (
	"bytes"
	"encoding/base64"
	"image"
	"log"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func GenerateImageSizesHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request", "error": err.Error()})
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding image", "error": err.Error()})
		return
	}

	img1 := imaging.Resize(img, 800, 600, imaging.Lanczos)
	img2 := imaging.Resize(img, 400, 300, imaging.Lanczos)

	buf1 := new(bytes.Buffer)
	_ = imaging.Encode(buf1, img1, imaging.JPEG)
	img1Base64 := base64.StdEncoding.EncodeToString(buf1.Bytes())

	buf2 := new(bytes.Buffer)
	_ = imaging.Encode(buf2, img2, imaging.JPEG)
	img2Base64 := base64.StdEncoding.EncodeToString(buf2.Bytes())

	c.JSON(http.StatusOK, gin.H{
		"images": []string{img1Base64, img2Base64},
	})
}
