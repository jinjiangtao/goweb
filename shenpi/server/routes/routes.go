package routes

import (
	"shenpi/controllers"
	"shenpi/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.Cors())

	userController := controllers.NewUserController()
	templateController := controllers.NewTemplateController()
	approvalController := controllers.NewApprovalController()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		user := api.Group("/user")
		user.Use(middleware.Auth())
		{
			user.GET("/profile", userController.Profile)
			user.GET("/list", userController.List)
			user.GET("/:id", userController.GetUser)
		}

		template := api.Group("/template")
		template.Use(middleware.Auth())
		{
			template.POST("", templateController.Create)
			template.POST("/default", templateController.CreateDefault)
			template.GET("", templateController.List)
			template.GET("/:id", templateController.Get)
			template.PUT("/:id", templateController.Update)
			template.DELETE("/:id", templateController.Delete)
			template.GET("/:id/fields", templateController.GetFormFields)
		}

		approval := api.Group("/approval")
		approval.Use(middleware.Auth())
		{
			approval.POST("", approvalController.Submit)
			approval.GET("/my", approvalController.MyRequests)
			approval.GET("/todo", approvalController.MyApprovals)
			approval.GET("/stats", approvalController.Stats)
			approval.GET("/:id", approvalController.Get)
			approval.GET("/:id/records", approvalController.Records)
			approval.GET("/:id/progress", approvalController.GetProgress)
			approval.POST("/:id/approve", approvalController.Approve)
			approval.POST("/:id/reject", approvalController.Reject)
			approval.POST("/:id/revoke", approvalController.Revoke)
			approval.POST("/:id/addsign", approvalController.AddSign)
			approval.POST("/:id/transfer", approvalController.Transfer)
		}
	}
}
