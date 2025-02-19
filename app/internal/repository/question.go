package repository

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/pkg/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type QuestionRepository interface {
	AddDataToEs(ctx context.Context, data []model.QuestionEs) error
	GetAllQuestion(ctx context.Context, time time.Time) ([]model.Question, error)
	GetQuestion(ctx context.Context, req *v1.QuestionRequest) ([]model.Question, int, error)
	GetCount(ctx context.Context) (int, error)
	Create(ctx context.Context, question *model.Question) error
	GetByTitle(ctx context.Context, title string) (*model.Question, error)
	GetByID(ctx context.Context, id uint64) (*model.Question, error)
	DeleteById(ctx context.Context, question *model.Question, id uint64) error
	Update(ctx context.Context, question *model.Question) error
	GetQuestionByBankId(ctx context.Context, bankId uint64) ([]model.Question, int64, error)
	GetEsQuestion(ctx context.Context, req *v1.QuestionRequest) ([]v1.Question, int, error)
	DeleteBatchQuestion(ctx context.Context, questions []string) error
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

func (r *questionRepository) DeleteBatchQuestion(ctx context.Context, questions []string) error {
	tx := r.DB(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 每次删除1000条
	for i := 0; i < len(questions); i += 1000 {
		for j := i; j < i+1000 && j < len(questions); j++ {
			if err := r.DB(ctx).Where("id = ?", questions[j]).Delete(&model.Question{}).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := r.DB(ctx).Where("question_id = ?", questions[j]).Delete(&model.QuestionBankQuestion{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (r *questionRepository) GetEsQuestion(ctx context.Context, req *v1.QuestionRequest) ([]v1.Question, int, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		// 排序字段
		"sort": []map[string]interface{}{},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				// 匹配条件,满足其中之一即可
				"should": []map[string]interface{}{},
				// 过滤条件
				"filter": []map[string]interface{}{
					{"term": map[string]interface{}{"is_delete": "0"}},
				},
			},
		},
	}
	if req.SortOrder != nil && req.SortField != nil {
		if *req.SortField == "createTime" {
			if *req.SortOrder == "ascend" {
				query["sort"] = append(query["sort"].([]map[string]interface{}),
					map[string]interface{}{"create_time": map[string]interface{}{
						"order": "asc",
					}})
			} else {
				query["sort"] = append(query["sort"].([]map[string]interface{}),
					map[string]interface{}{"create_time": map[string]interface{}{
						"order": "desc",
					}})
			}
		} else {
			if *req.SortOrder == "ascend" {
				query["sort"] = append(query["sort"].([]map[string]interface{}),
					map[string]interface{}{"update_time": map[string]interface{}{
						"order": "asc",
					}})
			} else {
				query["sort"] = append(query["sort"].([]map[string]interface{}),
					map[string]interface{}{"update_time": map[string]interface{}{
						"order": "desc",
					}})
			}
		}
	}
	if req.SearchText != nil && *req.SearchText != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["minimum_should_match"] = 1
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]map[string]interface{}),
			map[string]interface{}{
				"match": map[string]interface{}{"title": req.SearchText},
			},
			map[string]interface{}{
				"match": map[string]interface{}{"content": req.SearchText},
			},
			map[string]interface{}{
				"match": map[string]interface{}{"answer": req.SearchText},
			},
		)
	}
	if req.Tags != nil && len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = append(
				query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]map[string]interface{}),
				map[string]interface{}{
					"term": map[string]interface{}{"tags": tag},
				},
			)
		}
	}
	if req.ID != nil && *req.ID != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{"id": *req.ID},
			},
		)
	}
	if req.UserID != nil && *req.UserID != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{"user_id": *req.UserID},
			},
		)
	}
	if req.QuestionBankID != nil && *req.QuestionBankID != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{"question_id": *req.QuestionBankID},
			},
		)
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}
	var size, current int
	size = 12   // 默认值
	current = 1 // 默认值
	// 分页字段
	if req.PageSize != nil && req.Current != nil {
		size = *req.PageSize
		current = *req.Current
	}
	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("question"),
		r.es.Search.WithBody(&buf),
		r.es.Search.WithTrackTotalHits(true),
		r.es.Search.WithPretty(),
		r.es.Search.WithSize(size),
		// 分页字段
		r.es.Search.WithFrom(size*(current-1)),
	)
	if err != nil || res.IsError() {
		return nil, 0, err
	}
	var rr map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&rr); err != nil {
		return nil, 0, err
	}
	var total int
	var questions []v1.Question
	total = int(rr["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	for _, hit := range rr["hits"].(map[string]interface{})["hits"].([]interface{}) {
		if _, v := hit.(map[string]interface{})["_source"]; v {
			body, err := json.Marshal(hit.(map[string]interface{})["_source"])
			if err != nil {
				fmt.Println(err)
				continue
			}
			var q model.QuestionEs
			if err := json.Unmarshal(body, &q); err != nil {
				fmt.Println(err)
				continue
			}

			id := utils.Int64TOString(q.Id)
			userId := utils.Int64TOString(q.UserId)
			questions = append(questions, v1.Question{
				Answer:     &q.Answer,
				Content:    &q.Content,
				CreateTime: &q.CreateTime,
				EditTime:   &q.EditTime,
				ID:         &id,
				IsDelete:   &q.IsDelete,
				Tags:       &q.Tags,
				Title:      &q.Title,
				UpdateTime: &q.UpdateTime,
				UserID:     &userId,
			})
		}
	}
	return questions, total, nil
}

func (r *questionRepository) AddDataToEs(ctx context.Context, data []model.QuestionEs) error {
	for _, question := range data {
		body := map[string]interface{}{
			"doc":           question,
			"doc_as_upsert": true, // 设置 upsert 行为
		}
		marshal, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req := esapi.UpdateRequest{
			Index:      "question",
			DocumentID: strconv.FormatInt(question.Id, 10),
			Body:       bytes.NewReader(marshal),
		}
		if _, err = req.Do(ctx, r.es); err != nil {
			return err
		}
	}
	return nil
}

func (r *questionRepository) GetAllQuestion(ctx context.Context, time time.Time) ([]model.Question, error) {
	var questions []model.Question
	if err := r.DB(ctx).Unscoped().Where("update_time >= ?", time).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
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
	tx := r.DB(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := r.DB(ctx).Where("id = ?", id).Delete(&question).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := r.DB(ctx).Where("question_id = ?", id).Delete(&model.QuestionBankQuestion{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
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
		// TODO 未实现根据题库id查询
		_ = questionBankID
	}
	// TODO 未实现根据标签查询
	if err := r.DB(ctx).Joins("LEFT JOIN question_bank_question "+
		"ON question.id = question_bank_question.question_id").Where(""+
		"question.id LIKE ? AND "+
		"question.title LIKE ? AND "+
		"question.user_id LIKE ?",
		//"AND question_bank_question.question_bank_id LIKE ?",
		"%"+id+"%",
		"%"+title+"%",
		"%"+userId+"%",
		//"%"+questionBankID+"%",
	).Order(s).Distinct().Find(&questions).Group("question.id").Count(&total).Error; err != nil {
		return nil, 0, v1.ErrNotFound
	}
	return questions, int(total), nil
}
