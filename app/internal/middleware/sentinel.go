package middleware

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
)

func InitSentinel() {
	// 初始化 Sentinel
	err := api.InitDefault()
	if err != nil {
		log.Fatalf("Failed to initialize Sentinel: %v", err)
	}

	// 加载流量控制规则
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
}
