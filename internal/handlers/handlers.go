package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func GenerateImageSizesHandler(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
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

	encodeImage := func(img image.Image, format string) (string, error) {
		buf := new(bytes.Buffer)
		switch format {
		case "jpeg", "jpg":
			if err := jpeg.Encode(buf, img, nil); err != nil {
				return "", err
			}
		case "png":
			if err := png.Encode(buf, img); err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("unsupported format: %s", format)
		}
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	}

	formats := strings.Split(header.Filename, ".")
	imgFormat := formats[len(formats)-1]

	img1Base64, err := encodeImage(img1, imgFormat)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encoding image"})
		return
	}

	img2Base64, _ := encodeImage(img2, imgFormat)
	img3Base64, _ := encodeImage(img3, imgFormat)
	img4Base64, _ := encodeImage(img4, imgFormat)

	c.JSON(http.StatusOK, gin.H{
		"images": []string{img1Base64, img2Base64, img3Base64, img4Base64},
	})
}
