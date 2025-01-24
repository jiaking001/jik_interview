package repository

import (
	"app/internal/model"
	"context"
)

type QuestionBankRepository interface {
	GetQuestionBank(ctx context.Context, id int64) (*model.QuestionBank, error)
}

func NewQuestionBankRepository(
	repository *Repository,
) QuestionBankRepository {
	return &questionBankRepository{
		Repository: repository,
	}
}

type questionBankRepository struct {
	*Repository
}

func (r *questionBankRepository) GetQuestionBank(ctx context.Context, id int64) (*model.QuestionBank, error) {
	var questionBank model.QuestionBank

	return &questionBank, nil
}
