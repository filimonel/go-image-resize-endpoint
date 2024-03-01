package main

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/filimonel/go-image-resize-endpoint/internal/handlers"
	"github.com/filimonel/go-image-resize-endpoint/internal/middleware"
	ratelimiters "github.com/filimonel/go-image-resize-endpoint/internal/rate-limiters"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// Create instance of the Gin Engine
	r := gin.Default()

	// Set Cors
	r.Use(cors.New(middleware.SetCors()))

	// Set memory limit for multipart forms
	r.MaxMultipartMemory = 2 << 15

	// Define the rate limit rules
	limiter := rate.NewLimiter(rate.Every(time.Hour), 5)

	// Routes & Handlers
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello, world!",
		})
	})

	r.POST("/upload", ratelimiters.SetLimit(limiter), handlers.GenerateImageSizesHandler)

	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
