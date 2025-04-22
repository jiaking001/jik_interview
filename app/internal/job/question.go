package job

import (
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/rabbmit"
	"context"
	"encoding/json"
	"log"
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
	// 消费消息
	ConsumeMessage(ctx, t)

	// 每过一分钟向消息队列发送一次消息
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := syncDataToElasticsearch(ctx, t)
			if err != nil {
				t.logger.Error("Failed to sync data: %v")
			}
			t.logger.Info("Synced data to elasticsearch")
		}
	}
}

func syncDataToElasticsearch(ctx context.Context, t questionJob) error {
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	question, err := t.questionRepo.GetAllQuestion(ctx, fiveMinutesAgo)
	if err != nil {
		return err
	}

	err = rabbmit.SendMessage(question)
	if err != nil {
		return err
	}

	return nil
}

func ConsumeMessage(ctx context.Context, t questionJob) {

	go func() {
		// 消费消息
		msgs, err := rabbmit.Ch.Consume(
			"my_queue", // 队列名称
			"",         // 消费者名称
			false,      // 手动确认
			false,      // 非排他性
			false,      // 不等待服务器确认
			false,      // 不阻塞
			nil,        // 无额外参数
		)

		if err != nil {
			log.Fatalf("Failed to register a consumer: %s", err)
		}

		// 用于阻塞协程
		forever := make(chan bool)

		go func() {
			for d := range msgs {
				var msg []model.Question
				// 反序列化 JSON 数据为结构体
				if err = json.Unmarshal(d.Body, &msg); err != nil {
					log.Printf("Failed to unmarshal message: %s", err)
					err = d.Nack(false, false)
					if err != nil {
						return
					}
					continue
				}

				var data []model.QuestionEs
				for _, q := range msg {
					es := model.QuestionEs{
						Id:         int64(q.ID),
						UserId:     int64(q.UserID),
						EditTime:   q.EditTime,
						CreateTime: q.CreateTime,
						UpdateTime: q.UpdateTime,
						IsDelete:   q.IsDelete,
					}
					if q.Title != nil {
						es.Title = *q.Title
					}
					if q.Content != nil {
						es.Content = *q.Content
					}
					if q.Tags != nil {
						es.Tags = *q.Tags
					}
					if q.Answer != nil {
						es.Answer = *q.Answer
					}
					data = append(data, es)
				}

				err = t.questionRepo.AddDataToEs(ctx, data)
				if err != nil {
					return
				}

				err = d.Ack(false)
				if err != nil {
					return
				} // 手动确认消息
			}
		}()

		<-forever
	}()
}
