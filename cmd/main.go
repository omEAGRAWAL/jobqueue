package main

import (
	"jobqueue/config"
	"jobqueue/handler"
	"jobqueue/service"
	"jobqueue/utils"
	"jobqueue/worker"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitLogger()
	db := config.InitDB()
	defer db.Close()

	jobService := service.NewJobService(db, utils.Logger)
	pool := worker.WorkerPool{Jobs: make(chan service.Job, 100), DB: db, Log: utils.Logger}
	pool.Start(5)

	router := gin.Default()
	handler.RegisterRoutes(router, jobService)

	go jobService.EnqueuePendingJobs(pool.Jobs)

	//router.Run(":8081")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // fallback for local dev
	}
	router.Run(":" + port)

}
