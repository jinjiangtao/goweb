package router

import (
	"shuiti/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/dashboard", handlers.GetDashboard)

		questions := api.Group("/questions")
		{
			questions.GET("", handlers.GetQuestions)
			questions.GET("/all", handlers.GetQuestionsNoPaging)
			questions.GET("/categories", handlers.GetCategories)
			questions.GET("/:id", handlers.GetQuestion)
			questions.POST("", handlers.CreateQuestion)
			questions.POST("/batch", handlers.BatchCreateQuestions)
			questions.PUT("/:id", handlers.UpdateQuestion)
			questions.DELETE("/:id", handlers.DeleteQuestion)
			questions.POST("/delete-batch", handlers.BatchDeleteQuestions)
		}

		exam := api.Group("/exam")
		{
			exam.POST("/generate", handlers.GenerateExam)
			exam.POST("/submit", handlers.SubmitExam)
			exam.GET("/records", handlers.GetExamRecords)
			exam.GET("/records/:id", handlers.GetExamRecord)
		}

		wrong := api.Group("/wrong")
		{
			wrong.GET("", handlers.GetWrongQuestions)
			wrong.GET("/practice", handlers.GetWrongQuestionsForPractice)
			wrong.DELETE("/:id", handlers.DeleteWrongQuestion)
			wrong.DELETE("", handlers.ClearWrongQuestions)
			wrong.POST("/remove-correct", handlers.RemoveWrongAfterCorrect)
		}

		stats := api.Group("/stats")
		{
			stats.GET("", handlers.GetStatistics)
		}
	}

	return r
}
