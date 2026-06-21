package middleware

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return param.TimeStamp.Format("2006-01-02 15:04:05") + " | " +
			strconv.Itoa(param.StatusCode) + " | " +
			param.Latency.String() + " | " +
			param.ClientIP + " | " +
			param.Method + " | " +
			param.Path + "\n"
	})
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		log.Printf("Panic recovered: %v", err)
		c.JSON(500, gin.H{
			"error":   "Internal Server Error",
			"message": "Something went wrong",
		})
		c.Abort()
	})
}
