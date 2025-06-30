package handler

import (
	//"net/http"

	"github.com/gin-gonic/gin"
	"jobqueue/service"
)

func RegisterRoutes(r *gin.Engine, svc *service.JobService) {
	r.POST("/jobs", svc.SubmitJob)
	r.GET("/jobs/:id", svc.GetJob)
	r.GET("/jobs", svc.ListJobs)
	r.StaticFile("/", "index.html")
}
