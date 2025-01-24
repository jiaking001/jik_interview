package repository

import (
	"app/internal/model"
	"context"
)

type QuestionBankQuestionRepository interface {
	GetQuestionBankQuestion(ctx context.Context, id int64) (*model.QuestionBankQuestion, error)
}

func NewQuestionBankQuestionRepository(
	repository *Repository,
) QuestionBankQuestionRepository {
	return &questionBankQuestionRepository{
		Repository: repository,
	}
}

type questionBankQuestionRepository struct {
	*Repository
}

func (r *questionBankQuestionRepository) GetQuestionBankQuestion(ctx context.Context, id int64) (*model.QuestionBankQuestion, error) {
	var questionBankQuestion model.QuestionBankQuestion

	return &questionBankQuestion, nil
}
