package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/utils"
	"context"
	"strconv"
)

type QuestionService interface {
	ListQuestionByPage(ctx context.Context, req *v1.QuestionRequest) (v1.QuestionQueryResponseData[v1.Question], error)
	AddQuestion(ctx context.Context, req *v1.AddQuestionRequest, id uint64) (string, error)
	DeleteQuestion(ctx context.Context, req *v1.DeleteQuestionRequest) (bool, error)
	UpdateQuestion(ctx context.Context, req *v1.UpdateQuestionRequest) (bool, error)
	ListQuestionByBankId(ctx context.Context, bankId uint64) (v1.PageQuestionVO, error)
	GetQuestionById(ctx context.Context, req *v1.GetQuestionRequest) (v1.QuestionVO, error)
	ListQuestionVoByPage(ctx context.Context, req *v1.QuestionRequest) (v1.PageQuestionVO, error)
	SearchQuestionVoByPage(ctx context.Context, req *v1.QuestionRequest) (v1.PageQuestionVO, error)
	DeleteBatchQuestion(ctx context.Context, req *v1.BatchDeleteQuestionRequest) (bool, error)
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

func (s *questionService) DeleteBatchQuestion(ctx context.Context, req *v1.BatchDeleteQuestionRequest) (bool, error) {
	// 检验参数的合法性
	if len(req.QuestionIdList) == 0 {
		return false, v1.ParamsError
	}
	err := s.questionRepository.DeleteBatchQuestion(ctx, req.QuestionIdList)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *questionService) SearchQuestionVoByPage(ctx context.Context, req *v1.QuestionRequest) (v1.PageQuestionVO, error) {
	current := req.Current
	size := req.PageSize
	questions, total, err := s.questionRepository.GetEsQuestion(ctx, req)
	if err != nil {
		return v1.PageQuestionVO{}, err
	}
	var questionVOList []v1.QuestionVO

	for _, question := range questions {
		var id, userId string
		id = *question.ID
		userId = *question.UserID
		// 字符串转字符串数组
		tagList, err := utils.StringToStrings(*question.Tags)
		if err != nil {
			return v1.PageQuestionVO{}, err
		}

		q := v1.QuestionVO{
			Answer:     question.Answer,
			Content:    question.Content,
			CreateTime: question.CreateTime,
			ID:         &id,
			TagList:    tagList,
			Title:      question.Title,
			UpdateTime: question.UpdateTime,
			UserID:     &userId,
		}
		questionVOList = append(questionVOList, q)
	}
	pages := total / *size + 1
	return v1.PageQuestionVO{
		Records: questionVOList,
		Total:   &total,
		Pages:   &pages,
		Size:    size,
		Current: current,
	}, nil

}

func (s *questionService) ListQuestionVoByPage(ctx context.Context, req *v1.QuestionRequest) (v1.PageQuestionVO, error) {
	current := req.Current
	size := req.PageSize
	questions, total, err := s.questionRepository.GetQuestion(ctx, req)
	if err != nil {
		return v1.PageQuestionVO{}, err
	}
	var questionVOList []v1.QuestionVO

	for _, question := range questions {
		var id, userId string
		// uint64转string
		id = utils.Uint64TOString(question.ID)
		userId = utils.Uint64TOString(question.UserID)

		tagList, err := utils.StringToStrings(*question.Tags)
		if err != nil {
			return v1.PageQuestionVO{}, err
		}

		q := v1.QuestionVO{
			Answer:     question.Answer,
			Content:    question.Content,
			CreateTime: &question.CreateTime,
			ID:         &id,
			TagList:    tagList,
			Title:      question.Title,
			UpdateTime: &question.UpdateTime,
			UserID:     &userId,
		}
		questionVOList = append(questionVOList, q)
	}
	pages := total / *size + 1
	return v1.PageQuestionVO{
		Records: questionVOList,
		Total:   &total,
		Pages:   &pages,
		Size:    size,
		Current: current,
	}, nil
}

func (s *questionService) GetQuestionById(ctx context.Context, req *v1.GetQuestionRequest) (v1.QuestionVO, error) {
	if req.ID == nil || *req.ID == "" {
		return v1.QuestionVO{}, v1.ParamsError
	}

	id, err := strconv.ParseUint(*req.ID, 10, 64)
	if err != nil {
		return v1.QuestionVO{}, v1.ParamsError
	}
	question, err := s.questionRepository.GetByID(ctx, id)
	if err != nil {
		return v1.QuestionVO{}, err
	}

	tagList, err := utils.StringToStrings(*question.Tags)
	if err != nil {
		return v1.QuestionVO{}, err
	}

	userId := strconv.FormatUint(question.UserID, 10)
	return v1.QuestionVO{
		Answer:     question.Answer,
		Content:    question.Content,
		CreateTime: &question.CreateTime,
		ID:         req.ID,
		TagList:    tagList,
		Title:      question.Title,
		UpdateTime: &question.UpdateTime,
		UserID:     &userId,
	}, nil
}

func (s *questionService) ListQuestionByBankId(ctx context.Context, bankId uint64) (v1.PageQuestionVO, error) {
	questions, total, err := s.questionRepository.GetQuestionByBankId(ctx, bankId)
	if err != nil {
		return v1.PageQuestionVO{}, err
	}

	var questionList []v1.QuestionVO
	for _, question := range questions {
		var id, userId string
		id = strconv.Itoa(int(question.ID))
		userId = strconv.Itoa(int(question.UserID))

		tagList, err := utils.StringToStrings(*question.Tags)
		if err != nil {
			return v1.PageQuestionVO{}, err
		}

		q := v1.QuestionVO{
			Answer:     question.Answer,
			Content:    question.Content,
			CreateTime: &question.CreateTime,
			ID:         &id,
			TagList:    tagList,
			Title:      question.Title,
			UpdateTime: &question.UpdateTime,
			UserID:     &userId,
		}
		questionList = append(questionList, q)
	}
	t := int(total)
	return v1.PageQuestionVO{
		Records: questionList,
		Total:   &t,
	}, nil
}

func (s *questionService) UpdateQuestion(ctx context.Context, req *v1.UpdateQuestionRequest) (bool, error) {
	if req == nil || *req.ID == "" {
		return false, v1.ParamsError
	}

	id, err := utils.StringToUint64(*req.ID)
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
		tags := utils.StringsToString(req.Tags)
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

func (s *questionService) DeleteQuestion(ctx context.Context, req *v1.DeleteQuestionRequest) (bool, error) {
	if req.Id <= "0" {
		return false, v1.ParamsError
	}
	id, err := utils.StringToUint64(req.Id)
	if err != nil {
		return false, v1.ParamsError
	}
	bank, err := s.questionRepository.GetByID(ctx, id)
	if err != nil {
		return false, err
	}

	// 删除
	bank.IsDelete = 1
	err = s.questionRepository.DeleteById(ctx, bank, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *questionService) AddQuestion(ctx context.Context, req *v1.AddQuestionRequest, id uint64) (string, error) {
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
	tags := utils.StringsToString(req.Tags)

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
	questions, total, err := s.questionRepository.GetQuestion(ctx, req)
	if err != nil {
		return v1.QuestionQueryResponseData[v1.Question]{}, err
	}
	var questionList []v1.Question
	for _, question := range questions {
		var id, userId string
		id = utils.Uint64TOString(question.ID)
		userId = utils.Uint64TOString(question.UserID)
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
	pages := total / *size + 1
	return v1.QuestionQueryResponseData[v1.Question]{
		Records: questionList,
		Total:   &total,
		Pages:   &pages,
		Size:    size,
		Current: current,
	}, nil
}
