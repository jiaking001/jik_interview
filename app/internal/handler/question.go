package handler

import (
	v1 "app/api/v1"
	"app/internal/model"
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
	userInterface := session.Get("user_login")
	if userInterface == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	user := userInterface.(*model.User)
	if user == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	var req v1.AddQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	id, err := h.questionService.AddQuestion(ctx, &req, user.ID)
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
