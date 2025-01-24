package handler

import (
	"app/internal/service"
	"github.com/gin-gonic/gin"
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

func (h *QuestionHandler) GetQuestion(ctx *gin.Context) {

}
