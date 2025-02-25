package middleware

import (
	v1 "app/api/v1"
	"app/pkg/jwt"
	"app/pkg/utils"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// AntiCrawling 函数用于防止爬虫
func AntiCrawling(j *jwt.JWT, rdb *redis.Client, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获得 token
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
			v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
			ctx.Abort()
			return
		}
		// 获取 User-Agent
		userAgent := ctx.GetHeader("User-Agent")
		// 解析 User-Agent
		deviceType := utils.GetDeviceType(userAgent)

		// 获取用户id
		id := claims.User.ID
		// 获取当前的小时和分钟
		now := time.Now()
		hour := now.Hour()
		minute := now.Minute()
		// 计算从午夜开始到现在的总分钟数
		totalMinutes := hour*60 + minute

		// 定义 Redis 键
		countKey := fmt.Sprintf("user:%s:time:%s", utils.Uint64TOString(id), strconv.Itoa(totalMinutes))

		// 使用 Lua 脚本保证原子性
		script := `
			local count = redis.call('INCR', KEYS[1])
			if count == 1 then
				redis.call('EXPIRE', KEYS[1], ARGV[1])
			end
			return count
		`
		count, err := rdb.Eval(ctx, script, []string{countKey}, 180).Int()
		if err != nil {
			v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
			ctx.Abort()
			return
		}

		// 检查访问频率
		if count > 20 {
			// 封禁用户
			if err = db.Table("users").Where("id = ?", id).Update("user_role", "ban").Error; err != nil {
				v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
			}
			// 删除用户登录端的Token
			idStr := utils.Uint64TOString(id)
			if err = rdb.HDel(ctx, "user_tokens:"+idStr, deviceType).Err(); err != nil {
				v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
			}
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
