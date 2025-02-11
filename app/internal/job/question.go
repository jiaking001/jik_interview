package job

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"time"
)

type QuestionJob interface {
	DataToElasticsearch(ctx context.Context) error
}

func NewQuestionJob(
	job *Job,
	questionRepo repository.QuestionRepository,
) QuestionJob {
	return &questionJob{
		questionRepo: questionRepo,
		Job:          job,
	}
}

type questionJob struct {
	questionRepo repository.QuestionRepository
	*Job
}

func (t questionJob) DataToElasticsearch(ctx context.Context) error {
	// 每过一分钟执行一次
	//ticker := time.NewTicker(1 * time.Minute)
	//defer ticker.Stop()
	//
	//for {
	//	select {
	//	case <-ticker.C:
	err := syncDataToElasticsearch(ctx, t)
	if err != nil {
		t.logger.Info("Failed to sync data: %v")
	}
	t.logger.Info("Synced data to elasticsearch")
	//}
	//}
	return nil
}

func syncDataToElasticsearch(ctx context.Context, t questionJob) error {
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	question, err := t.questionRepo.GetAllQuestion(ctx, fiveMinutesAgo)
	if err != nil {
		return err
	}
	var data []model.QuestionEs
	for _, q := range question {
		es := model.QuestionEs{
			Id:         int64(q.ID),
			Title:      *q.Title,
			Content:    *q.Content,
			Tags:       nil,
			Answer:     *q.Answer,
			UserId:     int64(q.UserID),
			EditTime:   q.EditTime,
			CreateTime: q.CreateTime,
			UpdateTime: q.UpdateTime,
			IsDelete:   q.IsDelete,
		}
		data = append(data, es)
	}
	err = t.questionRepo.AddDataToEs(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
