package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
)

type QuestionService interface {
	GetQuestion(ctx context.Context, id int64) (*model.Question, error)
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

func (s *questionService) GetQuestion(ctx context.Context, id int64) (*model.Question, error) {
	return s.questionRepository.GetQuestion(ctx, id)
}
