package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/utils"
	"context"
	"strconv"
)

// QuestionService 定义了一个问题服务的接口
type QuestionService interface {
	// 根据页码获取问题列表
	ListQuestionByPage(ctx context.Context, req *v1.QuestionRequest) (v1.QuestionQueryResponseData[v1.Question], error)
	// 添加问题
	AddQuestion(ctx context.Context, req *v1.AddQuestionRequest, token string) (string, error)
	// 删除问题
	DeleteQuestion(ctx context.Context, req *v1.DeleteQuestionRequest) (bool, error)
	// 更新问题
	UpdateQuestion(ctx context.Context, req *v1.UpdateQuestionRequest) (bool, error)
	// 根据题库ID获取问题列表
	ListQuestionByBankId(ctx context.Context, bankId uint64) (v1.PageQuestionVO, error)
	// 根据问题ID获取问题
	GetQuestionById(ctx context.Context, req *v1.GetQuestionRequest) (v1.QuestionVO, error)
	// 根据页码获取问题列表（VO）
	ListQuestionVoByPage(ctx context.Context, req *v1.QuestionRequest) (v1.PageQuestionVO, error)
	// 根据页码和关键词搜索问题列表（VO）
	SearchQuestionVoByPage(ctx context.Context, req *v1.QuestionRequest) (v1.PageQuestionVO, error)
	// 批量删除问题
	DeleteBatchQuestion(ctx context.Context, req *v1.BatchDeleteQuestionRequest) (bool, error)
	// 通过 AI 生成题目
	AddQuestionByAI(ctx context.Context, req *v1.AddQuestionRequest, token string) (string, error)
}

// NewQuestionService 创建一个新的问题服务实例
func NewQuestionService(
	service *Service,
	questionRepository repository.QuestionRepository,
) QuestionService {
	return &questionService{
		Service:            service,
		questionRepository: questionRepository,
	}
}

// questionService 实现了QuestionService接口
type questionService struct {
	*Service
	questionRepository repository.QuestionRepository
}

func (s *questionService) AddQuestionByAI(ctx context.Context, req *v1.AddQuestionRequest, token string) (string, error) {
	// 将生成的题目添加到数据库
	questionId, err := s.AddQuestion(ctx, req, token)
	if err != nil {
		return "", err
	}
	return questionId, nil
}

// DeleteBatchQuestion 批量删除问题
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

// SearchQuestionVoByPage 根据页码和关键词搜索问题列表（VO）
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

// ListQuestionVoByPage 根据页码获取问题列表（VO）
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

// GetQuestionById 根据问题ID获取问题
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

// ListQuestionByBankId 根据题库ID获取问题列表
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

// UpdateQuestion 更新问题
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

// DeleteQuestion 删除问题
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

// AddQuestion 添加问题
func (s *questionService) AddQuestion(ctx context.Context, req *v1.AddQuestionRequest, token string) (string, error) {
	// 解析 token
	claims, err := s.jwt.ParseToken(token)
	if err != nil {
		return "", err
	}

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
		UserID:  claims.User.ID,
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

// ListQuestionByPage 根据页码获取问题列表
func (s *questionService) ListQuestionByPage(ctx context.Context, req *v1.QuestionRequest) (v1.QuestionQueryResponseData[v1.Question], error) {
	current := req.Current
	size := req.PageSize
	if req.PageSize == nil {
		return v1.QuestionQueryResponseData[v1.Question]{}, v1.ParamsError
	}
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
