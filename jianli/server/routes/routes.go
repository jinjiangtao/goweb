package routes

import (
	"jianli-server/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	templateCtrl := controllers.NewTemplateController()
	resumeCtrl := controllers.NewResumeController()
	versionCtrl := controllers.NewVersionController()
	pdfCtrl := controllers.NewPDFController()

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})

		templates := api.Group("/templates")
		{
			templates.GET("", templateCtrl.GetTemplates)
			templates.GET("/:id", templateCtrl.GetTemplate)
		}

		resumes := api.Group("/resumes")
		{
			resumes.GET("", resumeCtrl.GetResumes)
			resumes.GET("/:id", resumeCtrl.GetResume)
			resumes.POST("", resumeCtrl.CreateResume)
			resumes.PUT("/:id", resumeCtrl.UpdateResume)
			resumes.DELETE("/:id", resumeCtrl.DeleteResume)

			resumes.GET("/:id/versions", versionCtrl.GetVersions)
		}

		versions := api.Group("/versions")
		{
			versions.GET("/:id", versionCtrl.GetVersion)
			versions.POST("/:id/restore", versionCtrl.RestoreVersion)
		}

		pdf := api.Group("/pdf")
		{
			pdf.POST("/export", pdfCtrl.ExportPDF)
		}
	}
}
