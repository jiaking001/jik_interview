package repository

import (
	"app/internal/model"
	"context"
)

type QuestionRepository interface {
	GetQuestion(ctx context.Context, id int64) (*model.Question, error)
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

func (r *questionRepository) GetQuestion(ctx context.Context, id int64) (*model.Question, error) {
	var question model.Question

	return &question, nil
}
