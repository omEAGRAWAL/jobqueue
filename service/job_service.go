// package service
//
// import (
//
//	"net/http"
//	"strconv"
//
//	"github.com/gin-gonic/gin"
//	"github.com/jmoiron/sqlx"
//	"github.com/sirupsen/logrus"
//	"jobqueue/model"
//
// )
//
// type Job = model.Job
//
//	type JobService struct {
//		DB  *sqlx.DB
//		Log *logrus.Logger
//	}
//
//	func NewJobService(db *sqlx.DB, log *logrus.Logger) *JobService {
//		return &JobService{DB: db, Log: log}
//	}
//
//	func (s *JobService) SubmitJob(c *gin.Context) {
//		var job Job
//		err := s.DB.QueryRowx(
//			`INSERT INTO jobs (payload) VALUES ($1) RETURNING id`,
//			job.Payload,
//		).Scan(&job.ID)
//
//		if err != nil {
//			s.Log.Errorf("DB Insert error: %v", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert job"})
//			return
//		}
//
//		c.JSON(http.StatusOK, job)
//	}
//
//	func (s *JobService) GetJob(c *gin.Context) {
//		id, _ := strconv.Atoi(c.Param("id"))
//		var job Job
//		err := s.DB.Get(&job, "SELECT * FROM jobs WHERE id=$1", id)
//		if err != nil {
//			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
//			return
//		}
//		c.JSON(http.StatusOK, job)
//	}
//
//	func (s *JobService) ListJobs(c *gin.Context) {
//		var jobs []Job
//		err := s.DB.Select(&jobs, "SELECT * FROM jobs ORDER BY created_at DESC LIMIT 50")
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list jobs"})
//			return
//		}
//		c.JSON(http.StatusOK, jobs)
//	}
//
//	func (s *JobService) EnqueuePendingJobs(jobChan chan Job) {
//		for {
//			var jobs []Job
//			_ = s.DB.Select(&jobs, "SELECT * FROM jobs WHERE status='pending'")
//			for _, job := range jobs {
//				jobChan <- job
//			}
//		}
//	}
package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"jobqueue/model"
)

// Constants for all log messages
const (
	LogJobInsertFail   = "DB Insert error"
	LogJobFetchFail    = "DB Get error"
	LogJobsListFail    = "DB Select error in ListJobs"
	LogPendingFetchErr = "DB Select error in EnqueuePendingJobs"

	ErrInsertJob  = "Failed to insert job"
	ErrFetchJob   = "Job not found"
	ErrListJobs   = "Failed to list jobs"
	ErrPendingJob = "Failed to fetch pending jobs"
)

type Job = model.Job

type JobService struct {
	DB  *sqlx.DB
	Log *logrus.Logger
}

func NewJobService(db *sqlx.DB, log *logrus.Logger) *JobService {
	return &JobService{DB: db, Log: log}
}

func (s *JobService) SubmitJob(c *gin.Context) {
	var job Job
	if err := c.ShouldBindJSON(&job); err != nil {
		s.Log.WithError(err).Error("Invalid job payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job payload"})
		return
	}

	err := s.DB.QueryRowx(
		`INSERT INTO jobs (payload) VALUES ($1) RETURNING id`,
		job.Payload,
	).Scan(&job.ID)

	if err != nil {
		s.Log.WithError(err).Error(LogJobInsertFail)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInsertJob})
		return
	}

	s.Log.WithField("job_id", job.ID).Info("Job submitted")
	c.JSON(http.StatusOK, job)
}

func (s *JobService) GetJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.Log.WithError(err).Error("Invalid job ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var job Job
	err = s.DB.Get(&job, "SELECT * FROM jobs WHERE id=$1", id)
	if err != nil {
		s.Log.WithError(err).Error(LogJobFetchFail)
		c.JSON(http.StatusNotFound, gin.H{"error": ErrFetchJob})
		return
	}

	s.Log.WithField("job_id", job.ID).Info("Fetched job")
	c.JSON(http.StatusOK, job)
}

func (s *JobService) ListJobs(c *gin.Context) {
	var jobs []Job
	err := s.DB.Select(&jobs, "SELECT * FROM jobs ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		s.Log.WithError(err).Error(LogJobsListFail)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrListJobs})
		return
	}

	s.Log.Infof("Fetched %d jobs", len(jobs))
	c.JSON(http.StatusOK, jobs)
}

func (s *JobService) EnqueuePendingJobs(jobChan chan Job) {
	for {
		var jobs []Job
		err := s.DB.Select(&jobs, "SELECT * FROM jobs WHERE status='pending'")
		if err != nil {
			s.Log.WithError(err).Error(LogPendingFetchErr)
			continue
		}

		for _, job := range jobs {
			s.Log.WithField("job_id", job.ID).Info("Enqueuing pending job")
			jobChan <- job
		}
	}
}
