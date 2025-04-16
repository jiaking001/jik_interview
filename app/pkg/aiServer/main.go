package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"os"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func main() {
	// 初始化环境变量
	_ = godotenv.Load()
	// 请确保您已将 API Key 存储在环境变量 ARK_API_KEY 中
	// 初始化Ark客户端，从环境变量中读取您的API Key
	client := arkruntime.NewClientWithApiKey(
		// 从环境变量中获取您的 API Key。此为默认方式，您可根据需要进行修改
		os.Getenv("ARK_API_KEY"),
		// 此为默认路径，您可根据业务所在地域进行配置
		arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
	)

	ctx := context.Background()

	fmt.Println("----- standard request -----")
	req := model.CreateChatCompletionRequest{
		// 指定您创建的方舟推理接入点 ID，此处已帮您修改为您的推理接入点 ID
		Model: "deepseek-v3-250324",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是一名面试官."),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("总结一下10个最常见的后端面试题"),
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return
	}
	fmt.Println(*resp.Choices[0].Message.Content.StringValue)
}
