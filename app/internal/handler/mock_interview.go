package handler

import (
	v1 "app/api/v1"
	"app/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (h *MockInterviewHandler) AddMockInterview(ctx *gin.Context) {
	session := sessions.Default(ctx)
	t := session.Get("user_login")
	if t == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	token := t.(string)

	var req v1.MockInterviewAddRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.mockInterviewService.AddMockInterview(ctx, &req, token)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, ok)
}

func (h *MockInterviewHandler) GetMockInterview(ctx *gin.Context) {
	var req v1.MockInterviewGetRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.mockInterviewService.GetMockInterview(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, ok)
}

func (h *MockInterviewHandler) MockInterview(ctx *gin.Context) {
	var req v1.MockInterviewEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.mockInterviewService.MockInterview(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, ok)
}
