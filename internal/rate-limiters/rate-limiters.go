package ratelimiters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Define the Message struct as before
type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

// Create a Gin compatible rate limiter middleware
func SetLimit(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, try again later.",
			})
			return
		}
		c.Next()
	}
}
