package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/utils"
	"context"
	"strconv"
)

// QuestionBankQuestionService 定义了题目题库服务接口
type QuestionBankQuestionService interface {
	ListQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionQueryRequest) (v1.PageQuestionBankQuestionVO, error)
	AddQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionRequest) (string, error)
	RemoveQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionRequest) (bool, error)
	BatchAddQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionBatchRequest) (bool, error)
	BatchRemoveQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionBatchRemoveRequest) (bool, error)
}

// NewQuestionBankQuestionService 创建题目题库服务实例
func NewQuestionBankQuestionService(
	service *Service,
	questionBankQuestionRepository repository.QuestionBankQuestionRepository,
) QuestionBankQuestionService {
	return &questionBankQuestionService{
		Service:                        service,
		questionBankQuestionRepository: questionBankQuestionRepository,
	}
}

// questionBankQuestionService 题目题库服务结构体
type questionBankQuestionService struct {
	*Service
	questionBankQuestionRepository repository.QuestionBankQuestionRepository
}

// BatchRemoveQuestionBankQuestion 批量移除题目题库关系
func (s *questionBankQuestionService) BatchRemoveQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionBatchRemoveRequest) (bool, error) {
	// 判断参数是否合法
	if req.QuestionBankID == nil || req.QuestionIDList == nil || len(req.QuestionIDList) == 0 {
		return false, v1.ParamsError
	}
	questionBankID, err := utils.StringToUint64(*req.QuestionBankID)
	if err != nil {
		return false, err
	}
	// 批量移除题目题库关系
	var needRemoveQuestion []model.QuestionBankQuestion
	for _, questionID := range req.QuestionIDList {
		questionID, err := utils.StringToUint64(questionID)
		if err != nil {
			return false, err
		}
		needRemoveQuestion = append(needRemoveQuestion, model.QuestionBankQuestion{
			QuestionBankID: questionBankID,
			QuestionID:     questionID,
		})
	}
	if err = s.questionBankQuestionRepository.BatchRemoveQuestionBankQuestion(ctx, needRemoveQuestion); err != nil {
		return false, err
	}
	return true, nil
}

// BatchAddQuestionBankQuestion 批量添加题目题库关系
func (s *questionBankQuestionService) BatchAddQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionBatchRequest) (bool, error) {
	// 获取题库中已存在的题目
	var questionExistList = make(map[uint64]bool)
	if req.QuestionBankID == nil || len(req.QuestionIDList) == 0 {
		return false, v1.ParamsError
	}
	questionBankID, err := utils.StringToUint64(*req.QuestionBankID)
	if err != nil {
		return false, err
	}
	questions, err := s.questionBankQuestionRepository.GetQuestionBankQuestion(ctx, questionBankID, 1)
	if err != nil {
		return false, err
	}
	for _, question := range questions {
		questionExistList[question.QuestionID] = true
	}
	// 批量操作每次操作1000条
	for i := 0; i < len(req.QuestionIDList); i += 1000 {
		var needAddQuestion []model.QuestionBankQuestion
		for j := i; j < i+1000 && j < len(req.QuestionIDList); j++ {
			// 判断题库中是否已存在该题目
			questionID, err := utils.StringToUint64(req.QuestionIDList[j])
			if err != nil {
				return false, err
			}
			if _, ok := questionExistList[questionID]; ok {
				continue
			}
			needAddQuestion = append(needAddQuestion, model.QuestionBankQuestion{
				QuestionBankID: questionBankID,
				QuestionID:     questionID,
			})
		}
		err = s.questionBankQuestionRepository.BatchAddQuestionBankQuestion(ctx, needAddQuestion)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// RemoveQuestionBankQuestion 移除题目题库关系
func (s *questionBankQuestionService) RemoveQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionRequest) (bool, error) {
	questionID, err := strconv.ParseUint(*req.QuestionID, 10, 64)
	questionBankID, err := strconv.ParseUint(*req.QuestionBankID, 10, 64)
	if err != nil {
		return false, err
	}
	ok, err := s.questionBankQuestionRepository.RemoveQuestionBankQuestion(ctx, questionID, questionBankID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// AddQuestionBankQuestion 添加题目题库关系
func (s *questionBankQuestionService) AddQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionRequest) (string, error) {
	questionID, err := strconv.ParseUint(*req.QuestionID, 10, 64)
	questionBankID, err := strconv.ParseUint(*req.QuestionBankID, 10, 64)
	if err != nil {
		return "", err
	}
	err = s.questionBankQuestionRepository.AddQuestionBankQuestion(ctx, questionID, questionBankID)
	if err != nil {
		return "", err
	}
	id, err := s.questionBankQuestionRepository.GetQuestionBankQuestionId(ctx, questionID, questionBankID)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), err
}

// ListQuestionBankQuestion 获取题目题库列表
func (s *questionBankQuestionService) ListQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionQueryRequest) (v1.PageQuestionBankQuestionVO, error) {
	switch {
	// TODO
	case req.QuestionBankID != nil:
		return v1.PageQuestionBankQuestionVO{}, v1.ParamsError
	case req.QuestionID != nil:
		id, err := strconv.ParseUint(*req.QuestionID, 10, 64)
		questionBankQuestions, err := s.questionBankQuestionRepository.GetQuestionBankQuestion(ctx, id, 0)
		if err != nil {
			return v1.PageQuestionBankQuestionVO{}, err
		}
		var questionBankQuestionList []v1.QuestionBankQuestionVO
		for _, questionBankQuestion := range questionBankQuestions {
			questionBankId := strconv.FormatUint(questionBankQuestion.QuestionBankID, 10)
			questionId := strconv.FormatUint(questionBankQuestion.QuestionID, 10)
			q := v1.QuestionBankQuestionVO{
				CreateTime:     &questionBankQuestion.CreateTime,
				QuestionBankID: &questionBankId,
				QuestionID:     &questionId,
				UpdateTime:     &questionBankQuestion.UpdateTime,
			}
			questionBankQuestionList = append(questionBankQuestionList, q)
		}
		return v1.PageQuestionBankQuestionVO{
			Records: questionBankQuestionList,
			Size:    req.PageSize,
		}, nil
	default:
		return v1.PageQuestionBankQuestionVO{}, v1.ParamsError
	}
}
