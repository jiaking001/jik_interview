package handler

import (
	"app/internal/service"
	"github.com/gin-gonic/gin"
)

type MockInterviewHandler struct {
	*Handler
	mockInterviewService service.MockInterviewService
}

func NewMockInterviewHandler(
	handler *Handler,
	mockInterviewService service.MockInterviewService,
) *MockInterviewHandler {
	return &MockInterviewHandler{
		Handler:              handler,
		mockInterviewService: mockInterviewService,
	}
}

func (h *MockInterviewHandler) GetMockInterview(ctx *gin.Context) {

}
