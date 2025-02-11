package server

import (
	"app/internal/job"
	"app/pkg/log"
	"context"
)

type JobServer struct {
	log      *log.Logger
	userJob  job.UserJob
	question job.QuestionJob
}

func NewJobServer(
	log *log.Logger,
	userJob job.UserJob,
	question job.QuestionJob,
) *JobServer {
	return &JobServer{
		log:      log,
		userJob:  userJob,
		question: question,
	}
}

func (j *JobServer) Start(ctx context.Context) error {
	// Tips: If you want job to start as a separate process, just refer to the task implementation and adjust the code accordingly.

	// eg: kafka consumer
	err := j.question.DataToElasticsearch(ctx)
	return err
}
func (j *JobServer) Stop(ctx context.Context) error {
	return nil
}
