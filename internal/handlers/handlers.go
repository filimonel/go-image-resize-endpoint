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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding image"})
		return
	}

	img1 := imaging.Resize(img, 150, 150, imaging.Lanczos)
	img2 := imaging.Resize(img, 800, 600, imaging.Lanczos)
	img3 := imaging.Resize(img, 1280, 720, imaging.Lanczos)
	img4 := imaging.Resize(img, 1920, 1080, imaging.Lanczos)

	buf1 := new(bytes.Buffer)
	_ = imaging.Encode(buf1, img1, imaging.JPEG)
	img1Base64 := base64.StdEncoding.EncodeToString(buf1.Bytes())

	buf2 := new(bytes.Buffer)
	_ = imaging.Encode(buf2, img2, imaging.JPEG)
	img2Base64 := base64.StdEncoding.EncodeToString(buf2.Bytes())

	buf3 := new(bytes.Buffer)
	_ = imaging.Encode(buf3, img3, imaging.JPEG)
	img3Base64 := base64.StdEncoding.EncodeToString(buf3.Bytes())

	buf4 := new(bytes.Buffer)
	_ = imaging.Encode(buf4, img4, imaging.JPEG)
	img4Base64 := base64.StdEncoding.EncodeToString(buf4.Bytes())

	c.JSON(http.StatusOK, gin.H{
		"images": []string{img1Base64, img2Base64, img3Base64, img4Base64},
	})
}
