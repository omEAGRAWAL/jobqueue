
package worker

import (
    "fmt"

    "github.com/jmoiron/sqlx"
    "github.com/sirupsen/logrus"
    "jobqueue/service"
)

type WorkerPool struct {
    Jobs chan service.Job
    DB   *sqlx.DB
    Log  *logrus.Logger
}

func (wp *WorkerPool) Start(n int) {
    for i := 0; i < n; i++ {
        go func(id int) {
            for job := range wp.Jobs {
                wp.Log.Infof("Worker %d processing job %d", id, job.ID)
                result := fmt.Sprintf("Processed: %s", job.Payload)
                _, err := wp.DB.Exec("UPDATE jobs SET status=$1, result=$2, updated_at=NOW() WHERE id=$3", "completed", result, job.ID)
                if err != nil {
                    wp.Log.Errorf("Job %d failed: %v", job.ID, err)
                }
            }
        }(i)
    }
}
