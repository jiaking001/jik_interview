package handler

import (
	"app/internal/service"
	"github.com/gin-gonic/gin"
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

func (h *QuestionBankHandler) GetQuestionBank(ctx *gin.Context) {

}
