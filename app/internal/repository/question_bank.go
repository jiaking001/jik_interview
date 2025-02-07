package repository

import (
	v1 "app/api/v1"
	"app/internal/model"
	"context"
	"errors"
	"gorm.io/gorm"
)

type QuestionBankRepository interface {
	GetQuestionBank(ctx context.Context) ([]model.QuestionBank, error)
	GetByTitle(ctx context.Context, title string) (*model.QuestionBank, error)
	Create(ctx context.Context, bank *model.QuestionBank) error
	GetCount(ctx context.Context) (int, error)
	GetByID(ctx context.Context, id uint64) (*model.QuestionBank, error)
	DeleteById(ctx context.Context, bank *model.QuestionBank, id uint64) error
	Update(ctx context.Context, bank *model.QuestionBank) error
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

func (r *questionBankRepository) Update(ctx context.Context, bank *model.QuestionBank) error {
	if err := r.DB(ctx).Save(bank).Error; err != nil {
		return err
	}
	return nil
}

func (r *questionBankRepository) DeleteById(ctx context.Context, bank *model.QuestionBank, id uint64) error {
	if err := r.DB(ctx).Where("id = ?", id).Delete(&bank).Error; err != nil {
		return err
	}
	return nil
}

func (r *questionBankRepository) GetByID(ctx context.Context, id uint64) (*model.QuestionBank, error) {
	var bank model.QuestionBank
	if err := r.DB(ctx).Where("id = ?", id).First(&bank).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}
	return &bank, nil
}

func (r *questionBankRepository) GetCount(ctx context.Context) (int, error) {
	var total int64
	var questionBank model.QuestionBank
	if err := r.DB(ctx).Model(&questionBank).Count(&total).Error; err != nil {
		return 0, err
	}
	return int(total), nil
}

func (r *questionBankRepository) Create(ctx context.Context, bank *model.QuestionBank) error {
	if err := r.DB(ctx).Create(bank).Error; err != nil {
		return err
	}
	return nil
}

func (r *questionBankRepository) GetByTitle(ctx context.Context, title string) (*model.QuestionBank, error) {
	var questionBank model.QuestionBank
	if err := r.DB(ctx).Where("title = ?", title).First(&questionBank).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &questionBank, nil
}

func (r *questionBankRepository) GetQuestionBank(ctx context.Context) ([]model.QuestionBank, error) {
	var questionBanks []model.QuestionBank
	if err := r.DB(ctx).Find(&questionBanks).Error; err != nil {
		return nil, v1.ErrNotFound
	}
	return questionBanks, nil
}
