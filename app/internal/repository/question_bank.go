package repository

import (
	v1 "app/api/v1"
	"app/internal/model"
	"context"
)

type QuestionBankRepository interface {
	GetQuestionBank(ctx context.Context) ([]model.QuestionBank, error)
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

func (r *questionBankRepository) GetQuestionBank(ctx context.Context) ([]model.QuestionBank, error) {
	var questionBanks []model.QuestionBank
	if err := r.DB(ctx).Find(&questionBanks).Error; err != nil {
		return nil, v1.ErrNotFound
	}
	return questionBanks, nil
}
