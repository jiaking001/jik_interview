package middleware

import (
	v1 "app/api/v1"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

func CacheByRedis(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req v1.GetQuestionRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
			return
		}

		questionID := *req.ID
		cacheKey := fmt.Sprintf("question:cache:%s", questionID)
		countKey := fmt.Sprintf("question:count:%s", questionID)
		hotThreshold := 0

		// 访问次数加一
		rdb.Incr(ctx, countKey)

		// 判断是否是热点数据
		count, _ := rdb.Get(ctx, countKey).Int()
		isHot := count >= hotThreshold

		// 将热点判定结果存到context，方便handler使用
		ctx.Set("isHot", isHot)
		ctx.Set("cacheKey", cacheKey)

		if isHot {
			// 判断是否命中缓存
			cached, err := rdb.Get(ctx, cacheKey).Result()
			question := v1.QuestionVO{}
			err = json.Unmarshal([]byte(cached), &question)
			if err != nil {
				ctx.Next()
				return
			}
			if err == nil && cached != "" {
				v1.HandleSuccess(ctx, question)
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
