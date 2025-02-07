package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type QuestionService interface {
	ListQuestionByPage(ctx context.Context, req *v1.QuestionRequest) (v1.QuestionQueryResponseData[v1.Question], error)
	AddQuestion(ctx *gin.Context, req *v1.AddQuestionRequest, id uint64) (string, error)
	DeleteQuestion(ctx *gin.Context, req *v1.DeleteQuestionRequest) (bool, error)
	UpdateQuestion(ctx *gin.Context, req *v1.UpdateQuestionRequest) (bool, error)
}

func NewQuestionService(
	service *Service,
	questionRepository repository.QuestionRepository,
) QuestionService {
	return &questionService{
		Service:            service,
		questionRepository: questionRepository,
	}
}

type questionService struct {
	*Service
	questionRepository repository.QuestionRepository
}

func (s *questionService) UpdateQuestion(ctx *gin.Context, req *v1.UpdateQuestionRequest) (bool, error) {
	if req == nil || *req.ID == "" {
		return false, v1.ParamsError
	}

	id, err := strconv.ParseUint(*req.ID, 10, 64)
	if err != nil {
		return false, v1.ParamsError
	}
	question, err := s.questionRepository.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	if req.Title != nil && *req.Title != "" {
		question.Title = req.Title
	}
	if req.Answer != nil && *req.Answer != "" {
		question.Answer = req.Answer
	}

	if req.Tags != nil {
		// 将字符串数组转化为字符串
		quotedTags := make([]string, len(req.Tags))
		for i, tag := range req.Tags {
			quotedTags[i] = strconv.Quote(tag) // 使用 strconv.Quote 添加双引号
		}
		tags := strings.Join(quotedTags, ",")
		tags = "[" + tags + "]"

		question.Tags = &tags
	}

	if req.Content != nil && *req.Content != "" {
		question.Content = req.Content
	}

	err = s.questionRepository.Update(ctx, question)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *questionService) DeleteQuestion(ctx *gin.Context, req *v1.DeleteQuestionRequest) (bool, error) {
	if req.Id <= "0" {
		return false, v1.ParamsError
	}
	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return false, v1.ParamsError
	}
	bank, err := s.questionRepository.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	err = s.questionRepository.DeleteById(ctx, bank, id)
	if err != nil {
		return false, err
	}
	bank.IsDelete = 1
	err = s.questionRepository.Update(ctx, bank)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *questionService) AddQuestion(ctx *gin.Context, req *v1.AddQuestionRequest, id uint64) (string, error) {
	if *req.Title == "" {
		return "", v1.ErrIllegalAccount
	}
	questionBank, err := s.questionRepository.GetByTitle(ctx, *req.Title)
	if err != nil {
		return "", v1.ErrInternalServerError
	}
	if questionBank != nil {
		return "", v1.ErrTitleAlreadyUse
	}

	// 将字符串数组转化为字符串
	quotedTags := make([]string, len(req.Tags))
	for i, tag := range req.Tags {
		quotedTags[i] = strconv.Quote(tag) // 使用 strconv.Quote 添加双引号
	}
	tags := strings.Join(quotedTags, ",")
	tags = "[" + tags + "]"

	questionBank = &model.Question{
		Answer:  req.Answer,
		Content: req.Content,
		Tags:    &tags,
		Title:   req.Title,
		UserID:  id,
	}
	err = s.questionRepository.Create(ctx, questionBank)
	if err != nil {
		return "", err
	}
	var q *model.Question
	q, err = s.questionRepository.GetByTitle(ctx, *req.Title)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(q.ID, 10), nil
}

func (s *questionService) ListQuestionByPage(ctx context.Context, req *v1.QuestionRequest) (v1.QuestionQueryResponseData[v1.Question], error) {
	current := req.Current
	size := req.PageSize
	questions, err := s.questionRepository.GetQuestion(ctx)
	if err != nil {
		return v1.QuestionQueryResponseData[v1.Question]{}, err
	}
	var questionList []v1.Question
	for _, question := range questions {
		var id, userId string
		id = strconv.Itoa(int(question.ID))
		userId = strconv.Itoa(int(question.UserID))
		q := v1.Question{
			Answer:     question.Answer,
			Content:    question.Content,
			CreateTime: &question.CreateTime,
			EditTime:   &question.EditTime,
			ID:         &id,
			IsDelete:   &question.IsDelete,
			Tags:       question.Tags,
			Title:      question.Title,
			UpdateTime: &question.UpdateTime,
			UserID:     &userId,
		}
		questionList = append(questionList, q)
	}
	total, err := s.questionRepository.GetCount(ctx)
	if err != nil {
		return v1.QuestionQueryResponseData[v1.Question]{}, err
	}
	pages := total / *size + 1
	return v1.QuestionQueryResponseData[v1.Question]{
		Records: questionList,
		Total:   &total,
		Pages:   &pages,
		Size:    size,
		Current: current,
	}, nil
}
