package repository

import (
	v1 "app/api/v1"
	"app/internal/model"
	"context"
)

type MockInterviewRepository interface {
	GetMockInterview(ctx context.Context, id uint64) (*model.MockInterview, error)
	AddMockInterview(ctx context.Context, interview v1.MockInterview) (uint64, error)
	UpdateMockInterview(ctx context.Context, interview *model.MockInterview) error
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

func (r *mockInterviewRepository) UpdateMockInterview(ctx context.Context, interview *model.MockInterview) error {
	if err := r.DB(ctx).Save(interview).Error; err != nil {
		return err
	}
	return nil
}

// AddMockInterview 添加面试
func (r *mockInterviewRepository) AddMockInterview(ctx context.Context, interview v1.MockInterview) (uint64, error) {
	mockInterview := &model.MockInterview{
		CreateTime:     interview.CreateTime,
		Difficulty:     interview.Difficulty,
		JobPosition:    interview.JobPosition,
		UpdateTime:     interview.UpdateTime,
		UserID:         interview.UserID,
		WorkExperience: interview.WorkExperience,
		Status:         interview.Status, // 待开始模拟面试
	}
	if err := r.DB(ctx).Create(mockInterview).Error; err != nil {
		return 0, err
	}
	return mockInterview.ID, nil
}

func (r *mockInterviewRepository) GetMockInterview(ctx context.Context, id uint64) (*model.MockInterview, error) {
	var mockInterview model.MockInterview
	if err := r.DB(ctx).Where("id = ?", id).First(&mockInterview).Error; err != nil {
		return nil, err
	}
	return &mockInterview, nil
}
