package service

import (
	v1 "app/api/v1"
	"app/internal/repository"
	"context"
	"strconv"
)

type QuestionService interface {
	ListQuestionByPage(ctx context.Context, req *v1.QuestionRequest) (v1.QuestionQueryResponseData[v1.Question], error)
}

func NewQuestionService(
	service *Service,
	questionRepository repository.QuestionRepository,
) QuestionService {
	return &questionService{
		Service:            service,
		questionRepository: questionRepository,
	}
}

type questionService struct {
	*Service
	questionRepository repository.QuestionRepository
}

func (s *questionService) ListQuestionByPage(ctx context.Context, req *v1.QuestionRequest) (v1.QuestionQueryResponseData[v1.Question], error) {
	current := req.Current
	size := req.PageSize
	questions, err := s.questionRepository.GetQuestion(ctx)
	if err != nil {
		return v1.QuestionQueryResponseData[v1.Question]{}, err
	}
	var questionList []v1.Question
	for _, question := range questions {
		var id, userId string
		id = strconv.Itoa(int(question.ID))
		userId = strconv.Itoa(int(question.UserID))
		q := v1.Question{
			Answer:     question.Answer,
			Content:    question.Content,
			CreateTime: &question.CreateTime,
			EditTime:   &question.EditTime,
			ID:         &id,
			IsDelete:   &question.IsDelete,
			Tags:       question.Tags,
			Title:      question.Title,
			UpdateTime: &question.UpdateTime,
			UserID:     &userId,
		}
		questionList = append(questionList, q)
	}
	total := 10
	pages := 10
	return v1.QuestionQueryResponseData[v1.Question]{
		Records: questionList,
		Total:   &total,
		Pages:   &pages,
		Size:    size,
		Current: current,
	}, nil
}
