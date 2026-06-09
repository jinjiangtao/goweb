package routes

import (
	"neitui/handlers"
	"neitui/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/user/login", handlers.Login)

		public := api.Group("/public")
		{
			public.GET("/jobs", handlers.GetPublicJobs)
			public.GET("/jobs/:id", handlers.GetPublicJobDetail)
			public.POST("/submit", handlers.SubmitApplication)
			public.GET("/submissions", handlers.GetMySubmissions)
		}

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.POST("/jobs", handlers.CreateJob)
			auth.GET("/jobs/my", handlers.GetMyJobs)
			auth.PUT("/jobs/:id/status", handlers.UpdateJobStatus)

			auth.POST("/referrals", handlers.CreateReferral)
			auth.GET("/referrals/my", handlers.GetMyReferrals)
			auth.GET("/referrals/my/stats", handlers.GetMyReferralStats)

			hr := auth.Group("")
			hr.Use(middleware.RoleMiddleware("hr", "admin"))
			{
				hr.GET("/admin/referrals", handlers.GetAllReferrals)
				hr.PUT("/admin/referrals/:id/status", handlers.UpdateReferralStatus)
				hr.GET("/admin/stats", handlers.GetStats)
				hr.GET("/admin/jobs", handlers.GetAllJobs)
				hr.GET("/admin/referrals/export", handlers.ExportReferrals)

				hr.GET("/admin/submissions", handlers.GetAllSubmissions)
				hr.PUT("/admin/submissions/:id/status", handlers.UpdateSubmissionStatus)
				hr.POST("/admin/submissions/:id/convert", handlers.ConvertSubmissionToReferral)
			}

			admin := auth.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("/admin/users", handlers.GetUsers)
				admin.POST("/admin/users", handlers.CreateUser)
				admin.PUT("/admin/users/:id/status", handlers.UpdateUserStatus)
				admin.PUT("/admin/users/:id/reset-password", handlers.ResetPassword)
			}
		}
	}
}
