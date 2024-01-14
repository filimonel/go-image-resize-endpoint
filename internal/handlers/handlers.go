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
	// Store the image file(binary data), header and err
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	defer file.Close()

	// Decode the image binary, store the decoded Image and err
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding image"})
		return
	}

	// Generate image sizes and store them into memory to be encoded later
	img1 := imaging.Resize(img, 150, 150, imaging.Lanczos)
	img2 := imaging.Resize(img, 800, 600, imaging.Lanczos)
	img3 := imaging.Resize(img, 1280, 720, imaging.Lanczos)
	img4 := imaging.Resize(img, 1920, 1080, imaging.Lanczos)

	// Grab the filename, so that we can extract the image format(e.g. .jpg, .png, etc)
	// The image format will be passed into the encodeImage() function
	formats := strings.Split(header.Filename, ".")
	imgFormat := formats[len(formats)-1]

	// Encode and store the base64 generated images in memory
	img1Base64, err := encodeImage(img1, imgFormat)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encoding image"})
		return
	}

	img2Base64, _ := encodeImage(img2, imgFormat)
	img3Base64, _ := encodeImage(img3, imgFormat)
	img4Base64, _ := encodeImage(img4, imgFormat)

	// Send Response back to the client
	c.JSON(http.StatusOK, gin.H{
		"images": []string{img1Base64, img2Base64, img3Base64, img4Base64},
	})
}

func encodeImage(img image.Image, format string) (string, error) {
	// Create buffer to hold the []byte copy of the image
	buf := new(bytes.Buffer)

	// Check Image format and run the correct image encoding
	switch format {
	case "jpeg", "jpg":
		// Directly encode the image that's been passed into the function
		// Create a []byte representation(copy) of the image
		// Write the []byte copy of the image to the buffer
		err := jpeg.Encode(buf, img, nil)

		if err != nil {
			return "", err
		}
	case "png":
		err := png.Encode(buf, img)

		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	// NOTE: buf.Bytes() is a reference to buffer we opened at the top of the function.
	// Directly convert the image data([]byte) inside the buffer to a base64 string
	// Return the Base64 converted image.
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
