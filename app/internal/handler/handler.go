package handler

import (
	"app/pkg/jwt"
	"app/pkg/log"
	"app/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(
	logger *log.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return utils.Uint64TOString(v.(*jwt.MyCustomClaims).User.ID)
}
