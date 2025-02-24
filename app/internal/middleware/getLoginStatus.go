package middleware

import (
	v1 "app/api/v1"
	"app/pkg/jwt"
	"app/pkg/utils"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

// GetLoginStatus 获取用户登录状态
func GetLoginStatus(j *jwt.JWT, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		t := session.Get("user_login")
		if t == nil {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
			ctx.Abort()
			return
		}
		token := t.(string)
		// 解析 token
		claims, err := j.ParseToken(token)
		if err != nil {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		}
		// 获取 User-Agent
		userAgent := ctx.GetHeader("User-Agent")

		// 判断是否已登录
		// 解析 User-Agent
		deviceType := utils.GetDeviceType(userAgent)
		// 检查当前设备类型是否已经登录
		idStr := utils.Uint64TOString(claims.User.ID)
		nowToken, err := rdb.HGet(ctx, "user_tokens:"+idStr, deviceType).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
			ctx.Abort()
			return
		}
		if token != nowToken {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
