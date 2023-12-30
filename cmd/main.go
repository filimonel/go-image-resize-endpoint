package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/filimonel/go-image-resize-endpoint/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// Create instance of the Gin Engine
	r := gin.Default()

	// Set Cors
	r.Use(cors.New(middleware.SetCors()))

	// Routes & Handlers
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello, world!",
		})
	})

	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
