package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"context"
	"strconv"
)

type QuestionBankService interface {
	ListBankByPage(ctx context.Context, req *v1.QuestionBankRequest) (v1.QuestionBankQueryResponseData[v1.QuestionBank], error)
	AddQuestionBank(ctx context.Context, req *v1.AddQuestionBankRequest) (string, error)
	DeleteUser(ctx context.Context, req *v1.DeleteQuestionBankRequest) (bool, error)
	UpdateQuestionBank(ctx context.Context, req *v1.UpdateQuestionBankRequest) (bool, error)
	GetQuestionBankById(ctx context.Context, req *v1.GetQuestionBankRequest) (v1.GetQuestionBankResponse, error)
	ListBankByVOPage(ctx context.Context, req *v1.QuestionBankRequest) (v1.QuestionBankVO, error)
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

func (s *questionBankService) ListBankByVOPage(ctx context.Context, req *v1.QuestionBankRequest) (v1.QuestionBankVO, error) {
	if req.ID == nil || *req.ID == "" {
		return v1.QuestionBankVO{}, v1.ParamsError
	}

	id, err := strconv.ParseUint(*req.ID, 10, 64)
	if err != nil {
		return v1.QuestionBankVO{}, v1.ParamsError
	}
	bank, err := s.questionBankRepository.GetByID(ctx, id)
	if err != nil {
		return v1.QuestionBankVO{}, v1.ParamsError
	}

	return v1.QuestionBankVO{
		CreateTime:  &bank.CreateTime,
		Description: bank.Description,
		ID:          req.ID,
		Picture:     bank.Picture,
		Title:       bank.Title,
		UpdateTime:  &bank.UpdateTime,
		UserID:      req.UserID,
	}, nil
}

func (s *questionBankService) GetQuestionBankById(ctx context.Context, req *v1.GetQuestionBankRequest) (v1.GetQuestionBankResponse, error) {
	if req == nil || *req.ID == "" {
		return v1.GetQuestionBankResponse{}, v1.ParamsError
	}

	id, err := strconv.ParseUint(*req.ID, 10, 64)
	if err != nil {
		return v1.GetQuestionBankResponse{}, v1.ParamsError
	}
	bank, err := s.questionBankRepository.GetByID(ctx, id)
	if err != nil {
		return v1.GetQuestionBankResponse{}, err
	}

	return v1.GetQuestionBankResponse{
		CreateTime:  &bank.CreateTime,
		Description: bank.Description,
		ID:          req.ID,
		Picture:     bank.Picture,
		Title:       bank.Title,
		UpdateTime:  &bank.UpdateTime,
		UserID:      req.UserID,
	}, nil
}

func (s *questionBankService) UpdateQuestionBank(ctx context.Context, req *v1.UpdateQuestionBankRequest) (bool, error) {
	if req == nil || *req.ID == "" {
		return false, v1.ParamsError
	}

	id, err := strconv.ParseUint(*req.ID, 10, 64)
	if err != nil {
		return false, v1.ParamsError
	}
	bank, err := s.questionBankRepository.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	if req.Picture != nil && *req.Picture != "" {
		bank.Picture = req.Picture
	}
	if req.Title != nil && *req.Title != "" {
		bank.Title = req.Title
	}
	if req.Description != nil && *req.Description != "" {
		bank.Description = req.Description
	}

	err = s.questionBankRepository.Update(ctx, bank)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *questionBankService) DeleteUser(ctx context.Context, req *v1.DeleteQuestionBankRequest) (bool, error) {
	if req.Id <= "0" {
		return false, v1.ParamsError
	}
	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return false, v1.ParamsError
	}
	bank, err := s.questionBankRepository.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	err = s.questionBankRepository.DeleteById(ctx, bank, id)
	if err != nil {
		return false, err
	}
	bank.IsDelete = 1
	err = s.questionBankRepository.Update(ctx, bank)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *questionBankService) AddQuestionBank(ctx context.Context, req *v1.AddQuestionBankRequest) (string, error) {
	if *req.Title == "" {
		return "", v1.ErrIllegalAccount
	}
	questionBank, err := s.questionBankRepository.GetByTitle(ctx, *req.Title)
	if err != nil {
		return "", v1.ErrInternalServerError
	}
	if questionBank != nil {
		return "", v1.ErrTitleAlreadyUse
	}

	questionBank = &model.QuestionBank{
		Description: req.Description,
		Picture:     req.Picture,
		Title:       req.Title,
	}
	err = s.questionBankRepository.Create(ctx, questionBank)
	if err != nil {
		return "", err
	}
	var q *model.QuestionBank
	q, err = s.questionBankRepository.GetByTitle(ctx, *req.Title)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(q.ID, 10), nil
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
	total, err := s.questionBankRepository.GetCount(ctx)
	if err != nil {
		return v1.QuestionBankQueryResponseData[v1.QuestionBank]{}, err
	}
	pages := total / *size + 1
	return v1.QuestionBankQueryResponseData[v1.QuestionBank]{
		Records: questionBankList,
		Total:   &total,
		Pages:   &pages,
		Size:    size,
		Current: current,
	}, nil
}
