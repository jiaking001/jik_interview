package handler

import (
	v1 "app/api/v1"
	"app/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QuestionBankQuestionHandler struct {
	*Handler
	questionBankQuestionService service.QuestionBankQuestionService
}

func NewQuestionBankQuestionHandler(
	handler *Handler,
	questionBankQuestionService service.QuestionBankQuestionService,
) *QuestionBankQuestionHandler {
	return &QuestionBankQuestionHandler{
		Handler:                     handler,
		questionBankQuestionService: questionBankQuestionService,
	}
}

func (h *QuestionBankQuestionHandler) GetQuestionBankQuestion(ctx *gin.Context) {
	var req v1.QuestionBankQuestionQueryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	questionBankQuestion, err := h.questionBankQuestionService.ListQuestionBankQuestion(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, questionBankQuestion)
}

func (h *QuestionBankQuestionHandler) AddQuestionBankQuestion(ctx *gin.Context) {
	var req v1.QuestionBankQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	questionBankQuestionId, err := h.questionBankQuestionService.AddQuestionBankQuestion(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, questionBankQuestionId)
}

func (h *QuestionBankQuestionHandler) RemoveQuestionBankQuestion(ctx *gin.Context) {
	var req v1.QuestionBankQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.questionBankQuestionService.RemoveQuestionBankQuestion(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, ok)
}
