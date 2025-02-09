package repository

import (
	"app/internal/model"
	"context"
)

type QuestionBankQuestionRepository interface {
	GetQuestionBankQuestion(ctx context.Context, id uint64, flag int) ([]model.QuestionBankQuestion, error)
	AddQuestionBankQuestion(ctx context.Context, id uint64, id2 uint64) error
	GetQuestionBankQuestionId(ctx context.Context, id uint64, id2 uint64) (uint64, error)
	RemoveQuestionBankQuestion(ctx context.Context, id uint64, id2 uint64) (bool, error)
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

func (r *questionBankQuestionRepository) RemoveQuestionBankQuestion(ctx context.Context, id uint64, id2 uint64) (bool, error) {
	var questionBankQuestion model.QuestionBankQuestion
	if err := r.DB(ctx).Where("question_id = ? and question_bank_id = ?", id, id2).Delete(&questionBankQuestion).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *questionBankQuestionRepository) GetQuestionBankQuestionId(ctx context.Context, id uint64, id2 uint64) (uint64, error) {
	var questionBankQuestion model.QuestionBankQuestion
	if err := r.DB(ctx).Where("question_id = ? and question_bank_id = ?", id, id2).Find(&questionBankQuestion).Error; err != nil {
		return 0, err
	}
	return questionBankQuestion.ID, nil
}

func (r *questionBankQuestionRepository) AddQuestionBankQuestion(ctx context.Context, id uint64, id2 uint64) error {
	questionBankQuestion := model.QuestionBankQuestion{
		QuestionID:     id,
		QuestionBankID: id2,
	}
	if err := r.DB(ctx).Create(&questionBankQuestion).Error; err != nil {
		return err
	}
	return nil
}

func (r *questionBankQuestionRepository) GetQuestionBankQuestion(ctx context.Context, id uint64, flag int) ([]model.QuestionBankQuestion, error) {
	var questionBankQuestion []model.QuestionBankQuestion

	switch flag {
	case 0:
		if err := r.DB(ctx).Where("question_id = ?", id).Find(&questionBankQuestion).Error; err != nil {
			return nil, err
		}
	default:
		if err := r.DB(ctx).Where("question_bank_id = ?", id).Find(&questionBankQuestion).Error; err != nil {
			return nil, err
		}
	}

	return questionBankQuestion, nil
}
