package service

import (
	v1 "app/api/v1"
	"app/internal/repository"
	"context"
	"strconv"
)

type QuestionBankService interface {
	ListBankByPage(ctx context.Context, req *v1.QuestionBankRequest) (v1.QuestionBankQueryResponseData[v1.QuestionBank], error)
}

func NewQuestionBankService(
	service *Service,
	questionBankRepository repository.QuestionBankRepository,
) QuestionBankService {
	return &questionBankService{
		Service:                service,
		questionBankRepository: questionBankRepository,
	}
}

type questionBankService struct {
	*Service
	questionBankRepository repository.QuestionBankRepository
}

func (s *questionBankService) ListBankByPage(ctx context.Context, req *v1.QuestionBankRequest) (v1.QuestionBankQueryResponseData[v1.QuestionBank], error) {
	current := req.Current
	size := req.PageSize
	questionBanks, err := s.questionBankRepository.GetQuestionBank(ctx)
	if err != nil {
		return v1.QuestionBankQueryResponseData[v1.QuestionBank]{}, err
	}
	var questionBankList []v1.QuestionBank
	for _, questionBank := range questionBanks {
		var id, userId string
		id = strconv.Itoa(int(questionBank.ID))
		userId = strconv.Itoa(int(questionBank.UserID))
		q := v1.QuestionBank{
			CreateTime:  &questionBank.CreateTime,
			Description: questionBank.Description,
			EditTime:    &questionBank.EditTime,
			ID:          &id,
			IsDelete:    &questionBank.IsDelete,
			Picture:     questionBank.Picture,
			Title:       questionBank.Title,
			UpdateTime:  &questionBank.UpdateTime,
			UserID:      &userId,
		}
		questionBankList = append(questionBankList, q)
	}
	total := 10
	pages := 10
	return v1.QuestionBankQueryResponseData[v1.QuestionBank]{
		Records: questionBankList,
		Total:   &total,
		Pages:   &pages,
		Size:    size,
		Current: current,
	}, nil
}
