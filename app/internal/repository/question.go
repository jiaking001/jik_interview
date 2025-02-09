package repository

import (
	v1 "app/api/v1"
	"app/internal/model"
	"context"
	"errors"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	GetQuestion(ctx context.Context, req *v1.QuestionRequest) ([]model.Question, int, error)
	GetCount(ctx context.Context) (int, error)
	Create(ctx context.Context, question *model.Question) error
	GetByTitle(ctx context.Context, title string) (*model.Question, error)
	GetByID(ctx context.Context, id uint64) (*model.Question, error)
	DeleteById(ctx context.Context, question *model.Question, id uint64) error
	Update(ctx context.Context, question *model.Question) error
	GetQuestionByBankId(ctx context.Context, bankId uint64) ([]model.Question, int64, error)
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

func (r *questionRepository) GetQuestionByBankId(ctx context.Context, bankId uint64) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64
	if err := r.DB(ctx).Joins("INNER JOIN question_bank_question ON question.id = question_bank_question.question_id").
		Where("question_bank_question.question_bank_id = ?", bankId).
		Find(&questions).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return questions, total, nil
}

func (r *questionRepository) Update(ctx context.Context, question *model.Question) error {
	if err := r.DB(ctx).Save(question).Error; err != nil {
		return err
	}
	return nil
}

func (r *questionRepository) DeleteById(ctx context.Context, question *model.Question, id uint64) error {
	if err := r.DB(ctx).Where("id = ?", id).Delete(&question).Error; err != nil {
		return err
	}
	return nil
}

func (r *questionRepository) GetByID(ctx context.Context, id uint64) (*model.Question, error) {
	var question model.Question
	if err := r.DB(ctx).Where("id = ?", id).First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) GetByTitle(ctx context.Context, title string) (*model.Question, error) {
	var question model.Question
	if err := r.DB(ctx).Where("title = ?", title).First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) Create(ctx context.Context, question *model.Question) error {
	if err := r.DB(ctx).Create(question).Error; err != nil {
		return err
	}
	return nil
}

func (r *questionRepository) GetCount(ctx context.Context) (int, error) {
	var total int64
	var question model.Question
	if err := r.DB(ctx).Model(&question).Count(&total).Error; err != nil {
		return 0, err
	}
	return int(total), nil
}

func (r *questionRepository) GetQuestion(ctx context.Context, req *v1.QuestionRequest) ([]model.Question, int, error) {
	var questions []model.Question
	var total int64
	var s string
	if req.SortOrder != nil && req.SortField != nil {
		var sortOrder string
		var sortField string
		if *req.SortField == "createTime" {
			sortField = "question.create_time"
		} else {
			sortField = "question.update_time"
		}
		if *req.SortOrder == "ascend" {
			sortOrder = "asc"
		} else {
			sortOrder = "desc"
		}
		s = sortField + " " + sortOrder
	}
	var id, title, userId, questionBankID string
	if req.ID != nil {
		id = *req.ID
	}
	if req.Title != nil {
		title = *req.Title
	}
	if req.UserID != nil {
		userId = *req.UserID
	}
	if req.QuestionBankID != nil {
		questionBankID = *req.QuestionBankID
	}
	// TODO 未实现根据标签查询
	if err := r.DB(ctx).Joins("INNER JOIN question_bank_question "+
		"ON question.id = question_bank_question.question_id").Where(""+
		"question.id LIKE ? AND "+
		"question.title LIKE ? AND "+
		"question.user_id LIKE ? AND "+
		"question_bank_question.question_bank_id LIKE ?",
		"%"+id+"%",
		"%"+title+"%",
		"%"+userId+"%",
		"%"+questionBankID+"%",
	).Order(s).Distinct().Find(&questions).Group("question.id").Count(&total).Error; err != nil {
		return nil, 0, v1.ErrNotFound
	}
	return questions, int(total), nil
}
