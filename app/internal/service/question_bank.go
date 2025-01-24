package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
)

type QuestionBankService interface {
	GetQuestionBank(ctx context.Context, id int64) (*model.QuestionBank, error)
}

func NewQuestionBankService(
	service *Service,
	questionBankRepository repository.QuestionBankRepository,
) QuestionBankService {
	return &questionBankService{
		Service:                service,
		questionBankRepository: questionBankRepository,
	}
}

type questionBankService struct {
	*Service
	questionBankRepository repository.QuestionBankRepository
}

func (s *questionBankService) GetQuestionBank(ctx context.Context, id int64) (*model.QuestionBank, error) {
	return s.questionBankRepository.GetQuestionBank(ctx, id)
}
