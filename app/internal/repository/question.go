package repository

import (
	v1 "app/api/v1"
	"app/internal/model"
	"context"
)

type QuestionRepository interface {
	GetQuestion(ctx context.Context) ([]model.Question, error)
}

func NewQuestionRepository(
	repository *Repository,
) QuestionRepository {
	return &questionRepository{
		Repository: repository,
	}
}

type questionRepository struct {
	*Repository
}

func (r *questionRepository) GetQuestion(ctx context.Context) ([]model.Question, error) {
	var questions []model.Question
	if err := r.DB(ctx).Find(&questions).Error; err != nil {
		return nil, v1.ErrNotFound
	}
	return questions, nil
}
