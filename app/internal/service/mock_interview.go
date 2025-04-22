package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
)

type MockInterviewService interface {
	GetMockInterview(ctx context.Context, id int64) (*model.MockInterview, error)
}

func NewMockInterviewService(
	service *Service,
	mockInterviewRepository repository.MockInterviewRepository,
) MockInterviewService {
	return &mockInterviewService{
		Service:                 service,
		mockInterviewRepository: mockInterviewRepository,
	}
}

type mockInterviewService struct {
	*Service
	mockInterviewRepository repository.MockInterviewRepository
}

func (s *mockInterviewService) GetMockInterview(ctx context.Context, id int64) (*model.MockInterview, error) {
	return s.mockInterviewRepository.GetMockInterview(ctx, id)
}
