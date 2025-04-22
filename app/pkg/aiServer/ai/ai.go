package ai

import (
	"context"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"log"
	"os"
)

// DEFAULT_MODEL 默认使用 v3
const DEFAULT_MODEL = "deepseek-v3-250324"

func DoChat(systemPrompt string, userPrompt string, opt ...string) string {
	// _ = godotenv.Load(".env")
	// 请确保您已将 API Key 存储在环境变量 ARK_API_KEY 中
	// 初始化Ark客户端，从环境变量中读取您的API Key
	client := arkruntime.NewClientWithApiKey(
		// 从环境变量中获取您的 API Key。此为默认方式，您可根据需要进行修改
		os.Getenv("ARK_API_KEY"),
		// 此为默认路径，您可根据业务所在地域进行配置
		arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
	)

	ctx := context.Background()

	m := DEFAULT_MODEL
	if opt != nil && opt[0] != "" {
		m = opt[0]
	}

	// fmt.Println("----- standard request -----")
	req := model.CreateChatCompletionRequest{
		// 指定您创建的方舟推理接入点 ID，此处已帮您修改为您的推理接入点 ID
		Model: m,
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(systemPrompt),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(userPrompt),
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		// fmt.Printf("standard chat error: %v\n", err)
		log.Printf("standard chat error: %v\n", err)
		return ""
	}
	// fmt.Println(*resp.Choices[0].Message.Content.StringValue)
	return *resp.Choices[0].Message.Content.StringValue
}
