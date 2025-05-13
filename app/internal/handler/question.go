package handler

import (
	v1 "app/api/v1"
	"app/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QuestionHandler struct {
	*Handler
	questionService service.QuestionService
}

func NewQuestionHandler(
	handler *Handler,
	questionService service.QuestionService,
) *QuestionHandler {
	return &QuestionHandler{
		Handler:         handler,
		questionService: questionService,
	}
}

func (h *QuestionHandler) ListPage(ctx *gin.Context) {
	// TODO 排序，查找特定值未实现
	var req v1.QuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	page, err := h.questionService.ListQuestionByPage(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.QuestionBankQueryResponseData[v1.Question]{
		Records: page.Records,
		Size:    page.Size,
		Total:   page.Total,
		Current: page.Current,
		Pages:   page.Pages,
	})
}

func (h *QuestionHandler) AddQuestion(ctx *gin.Context) {
	session := sessions.Default(ctx)
	t := session.Get("user_login")
	if t == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	token := t.(string)

	var req v1.AddQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	id, err := h.questionService.AddQuestion(ctx, &req, token)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, id)
}

func (h *QuestionHandler) DeleteQuestion(ctx *gin.Context) {
	var req v1.DeleteQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.questionService.DeleteQuestion(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, ok)
}

func (h *QuestionHandler) UpdateQuestion(ctx *gin.Context) {
	var req v1.UpdateQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.questionService.UpdateQuestion(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, ok)
}

func (h *QuestionHandler) GetQuestion(ctx *gin.Context) {
	var req v1.GetQuestionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 获取题目热度信息
	isHot, _ := ctx.Get("isHot")
	cacheKey, _ := ctx.Get("cacheKey")

	question, err := h.questionService.GetQuestionById(ctx, &req, isHot, cacheKey)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, question)
}

func (h *QuestionHandler) ListPageVo(ctx *gin.Context) {
	var req v1.QuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	question, err := h.questionService.ListQuestionVoByPage(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, question)
}

func (h *QuestionHandler) SearchPageVo(ctx *gin.Context) {
	var req v1.QuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	question, err := h.questionService.SearchQuestionVoByPage(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, question)
}

func (h *QuestionHandler) DeleteBatchQuestion(ctx *gin.Context) {
	var req v1.BatchDeleteQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.questionService.DeleteBatchQuestion(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, ok)
}

func (h *QuestionHandler) AiGenerateQuestion(ctx *gin.Context) {
	session := sessions.Default(ctx)
	t := session.Get("user_login")
	if t == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	token := t.(string)

	var req v1.AddQuestionByAIRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.questionService.AddQuestionByAI(ctx, &req, token)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, ok)
}
