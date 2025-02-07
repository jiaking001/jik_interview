package handler

import (
	v1 "app/api/v1"
	"app/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QuestionBankHandler struct {
	*Handler
	questionBankService service.QuestionBankService
}

func NewQuestionBankHandler(
	handler *Handler,
	questionBankService service.QuestionBankService,
) *QuestionBankHandler {
	return &QuestionBankHandler{
		Handler:             handler,
		questionBankService: questionBankService,
	}
}

func (h *QuestionBankHandler) ListPage(ctx *gin.Context) {
	// TODO 排序，查找特定值未实现
	var req v1.QuestionBankRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	page, err := h.questionBankService.ListBankByPage(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.QuestionBankQueryResponseData[v1.QuestionBank]{
		Records: page.Records,
		Size:    page.Size,
		Total:   page.Total,
		Current: page.Current,
		Pages:   page.Pages,
	})
}
