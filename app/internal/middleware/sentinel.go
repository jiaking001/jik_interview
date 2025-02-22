package middleware

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
)

func InitSentinel() {
	// 初始化 Sentinel
	err := api.InitDefault()
	if err != nil {
		log.Fatalf("Failed to initialize Sentinel: %v", err)
	}

	// 加载流量控制规则（限流）
	_, err = flow.LoadRules([]*flow.Rule{
		// 对查看题库列表限流
		{
			Resource:               "POST:/api/questionBank/list/page/vo",
			Threshold:              60, // 每秒最多允许60个请求
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
		// 对查看题目列表限流
		{
			Resource:               "POST:/api/question/list/page/vo",
			Threshold:              60, // 每秒最多允许60个请求
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
		// 对搜索题目列表限流
		{
			Resource:               "POST:/api/question/search/page/vo",
			Threshold:              60, // 每秒最多允许60个请求
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})

	if err != nil {
		log.Fatalf("Failed to load flow rules: %v", err)
	}

	// 配置熔断规则
	rule := circuitbreaker.Rule{
		Resource:         "POST:/api/questionBank/list/page/vo",
		Strategy:         circuitbreaker.ErrorRatio,
		Threshold:        0.2,  // 异常比例阈值
		MinRequestAmount: 10,   // 最小请求数
		StatIntervalMs:   1000, // 统计时间窗口
		RetryTimeoutMs:   5000, // 熔断持续时间
	}
	if _, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{&rule}); err != nil {
		log.Fatalf("Failed to load circuit breaker rules: %+v", err)
	}
}
