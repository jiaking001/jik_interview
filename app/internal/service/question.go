package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

type QuestionService interface {
	ListQuestionByPage(ctx context.Context, req *v1.QuestionRequest) (v1.QuestionQueryResponseData[v1.Question], error)
	AddQuestion(ctx context.Context, req *v1.AddQuestionRequest, id uint64) (string, error)
	DeleteQuestion(ctx context.Context, req *v1.DeleteQuestionRequest) (bool, error)
	UpdateQuestion(ctx context.Context, req *v1.UpdateQuestionRequest) (bool, error)
	ListQuestionByBankId(ctx context.Context, id uint64) (v1.PageQuestionVO, error)
	GetQuestionById(ctx context.Context, req *v1.GetQuestionRequest) (v1.QuestionVO, error)
	ListQuestionVoByPage(ctx context.Context, req *v1.QuestionRequest) (v1.PageQuestionVO, error)
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

func (s *questionService) ListQuestionVoByPage(ctx context.Context, req *v1.QuestionRequest) (v1.PageQuestionVO, error) {
	current := req.Current
	size := req.PageSize
	questions, err := s.questionRepository.GetQuestion(ctx)
	if err != nil {
		return v1.PageQuestionVO{}, err
	}
	var questionVOList []v1.QuestionVO

	for _, question := range questions {
		var id, userId string
		id = strconv.Itoa(int(question.ID))
		userId = strconv.Itoa(int(question.UserID))

		// 原始的 JSON 格式的字符串
		jsonStr := question.Tags

		// 定义一个字符串数组变量来存储解析后的结果
		var tagList []string

		// 使用 json.Unmarshal 解析 JSON 字符串
		err = json.Unmarshal([]byte(*jsonStr), &tagList)
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
	total, err := s.questionRepository.GetCount(ctx)
	if err != nil {
		return v1.PageQuestionVO{}, err
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

	// 原始的 JSON 格式的字符串
	jsonStr := question.Tags

	// 定义一个字符串数组变量来存储解析后的结果
	var tagList []string

	// 使用 json.Unmarshal 解析 JSON 字符串
	err = json.Unmarshal([]byte(*jsonStr), &tagList)
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

func (s *questionService) ListQuestionByBankId(ctx context.Context, id uint64) (v1.PageQuestionVO, error) {
	//TODO 未实现根据题库id查询题库
	questions, err := s.questionRepository.GetQuestion(ctx)
	if err != nil {
		return v1.PageQuestionVO{}, err
	}

	var questionList []v1.QuestionVO
	for _, question := range questions {
		var id, userId string
		id = strconv.Itoa(int(question.ID))
		userId = strconv.Itoa(int(question.UserID))

		// 原始的 JSON 格式的字符串
		jsonStr := question.Tags

		// 定义一个字符串数组变量来存储解析后的结果
		var tagList []string

		// 使用 json.Unmarshal 解析 JSON 字符串
		err = json.Unmarshal([]byte(*jsonStr), &tagList)
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

	total, err := s.questionRepository.GetCount(ctx)
	if err != nil {
		return v1.PageQuestionVO{}, err
	}

	return v1.PageQuestionVO{
		Records: questionList,
		Total:   &total,
	}, nil
}

func (s *questionService) UpdateQuestion(ctx context.Context, req *v1.UpdateQuestionRequest) (bool, error) {
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

func (s *questionService) DeleteQuestion(ctx context.Context, req *v1.DeleteQuestionRequest) (bool, error) {
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
