package service

import (
	v1 "app/api/v1"
	"app/internal/repository"
	"context"
	"strconv"
)

type QuestionBankQuestionService interface {
	ListQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionQueryRequest) (v1.PageQuestionBankQuestionVO, error)
	AddQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionRequest) (string, error)
	RemoveQuestionBankQuestion(ctx context.Context, req *v1.QuestionBankQuestionRequest) (bool, error)
}

func NewQuestionBankQuestionService(
	service *Service,
	questionBankQuestionRepository repository.QuestionBankQuestionRepository,
) QuestionBankQuestionService {
	return &questionBankQuestionService{
		Service:                        service,
		questionBankQuestionRepository: questionBankQuestionRepository,
	}
}

type questionBankQuestionService struct {
	*Service
	questionBankQuestionRepository repository.QuestionBankQuestionRepository
}

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
