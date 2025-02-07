package handler

import (
	v1 "app/api/v1"
	"app/internal/service"
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
