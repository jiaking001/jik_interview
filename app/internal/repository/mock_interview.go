package repository

import (
	"app/internal/model"
	"context"
)

type MockInterviewRepository interface {
	GetMockInterview(ctx context.Context, id int64) (*model.MockInterview, error)
}

func NewMockInterviewRepository(
	repository *Repository,
) MockInterviewRepository {
	return &mockInterviewRepository{
		Repository: repository,
	}
}

type mockInterviewRepository struct {
	*Repository
}

func (r *mockInterviewRepository) GetMockInterview(ctx context.Context, id int64) (*model.MockInterview, error) {
	var mockInterview model.MockInterview

	return &mockInterview, nil
}
