package repository

import (
	v1 "app/api/v1"
	"app/internal/model"
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type QuestionBankRepository interface {
	GetQuestionBank(ctx context.Context, req *v1.QuestionBankRequest) ([]model.QuestionBank, int, error)
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
	tx := r.DB(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(bank).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id = ?", id).Delete(bank).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("question_bank_id = ?", id).Delete(&model.QuestionBankQuestion{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
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

func (r *questionBankRepository) GetQuestionBank(ctx context.Context, req *v1.QuestionBankRequest) ([]model.QuestionBank, int, error) {
	var questionBanks []model.QuestionBank
	var total int64
	var s string
	if req.SortOrder != nil && req.SortField != nil {
		var sortOrder string
		var sortField string
		if *req.SortField == "createTime" {
			sortField = "create_time"
		} else {
			sortField = "update_time"
		}
		if *req.SortOrder == "ascend" {
			sortOrder = "asc"
		} else {
			sortOrder = "desc"
		}
		s = sortField + " " + sortOrder
	}
	// 拼接查询字符串
	// query 查询字符串
	var conditions []string
	var params []interface{}
	var query string
	var id, title, description string
	if req.ID != nil {
		id = *req.ID
		conditions = append(conditions, "id LIKE ?")
		params = append(params, "%"+id+"%")
	}
	if req.Title != nil {
		title = *req.Title
		conditions = append(conditions, "title LIKE ?")
		params = append(params, "%"+title+"%")
	}
	if req.Description != nil {
		description = *req.Description
		conditions = append(conditions, "description LIKE ?")
		params = append(params, "%"+description+"%")
	}

	// 构造完整的查询条件
	query = strings.Join(conditions, " AND ")

	if err := r.DB(ctx).Where(query, params...).Order(s).Find(&questionBanks).Count(&total).Error; err != nil {
		return nil, 0, v1.ErrNotFound
	}
	return questionBanks, int(total), nil
}
