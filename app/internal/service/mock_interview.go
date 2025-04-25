package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/aiServer/ai"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	chatModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"strings"
	"time"
)

type MockInterviewService interface {
	MockInterview(ctx context.Context, req *v1.MockInterviewEventRequest) (string, error)
	AddMockInterview(ctx context.Context, req *v1.MockInterviewAddRequest, token string) (uint64, error)
	GetMockInterview(ctx *gin.Context, v *v1.MockInterviewGetRequest) (v1.MockInterview, error)
}

func NewMockInterviewService(
	service *Service,
	mockInterviewRepository repository.MockInterviewRepository,
) MockInterviewService {
	return &mockInterviewService{
		Service:                 service,
		mockInterviewRepository: mockInterviewRepository,
	}
}

type mockInterviewService struct {
	*Service
	mockInterviewRepository repository.MockInterviewRepository
}

func (m mockInterviewService) GetMockInterview(ctx *gin.Context, req *v1.MockInterviewGetRequest) (v1.MockInterview, error) {
	mockInterview, err := m.mockInterviewRepository.GetMockInterview(ctx, req.ID)
	if err != nil {
		return v1.MockInterview{}, err
	}

	return v1.MockInterview{
		CreateTime:     mockInterview.CreateTime,
		Difficulty:     mockInterview.Difficulty,
		ID:             mockInterview.ID,
		IsDelete:       mockInterview.IsDelete,
		JobPosition:    mockInterview.JobPosition,
		Messages:       mockInterview.Messages,
		Status:         mockInterview.Status,
		UpdateTime:     mockInterview.UpdateTime,
		UserID:         mockInterview.UserID,
		WorkExperience: mockInterview.WorkExperience,
	}, nil
}

// AddMockInterview 添加模拟面试
func (m mockInterviewService) AddMockInterview(ctx context.Context, req *v1.MockInterviewAddRequest, token string) (uint64, error) {
	// 解析 token
	claims, err := m.jwt.ParseToken(token)
	if err != nil {
		return 0, err
	}
	// 获取用户 ID
	userId := claims.User.ID
	// 创建 MockInterview
	interview := v1.MockInterview{
		CreateTime:     time.Time{},
		Difficulty:     req.Difficulty,
		JobPosition:    req.JobPosition,
		UpdateTime:     time.Time{},
		UserID:         userId,
		WorkExperience: req.WorkExperience,
		Status:         0, // 待开始模拟面试
	}
	// 将数据添加到数据库
	id, err := m.mockInterviewRepository.AddMockInterview(ctx, interview)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m mockInterviewService) MockInterview(ctx context.Context, req *v1.MockInterviewEventRequest) (string, error) {
	// 获取模拟模拟面试的信息
	mockInterview, err := m.mockInterviewRepository.GetMockInterview(ctx, req.ID)
	if err != nil {
		return "", err
	}
	// 定义系统 Prompt
	systemPrompt := fmt.Sprintf("你是一位严厉的程序员面试官，我是候选人，来应聘 %s 的 %s 岗位，面试难度为 %s。请你向我依次提出问题（最多 20 个问题），我也会依次回复。在这期间请完全保持真人面试官的口吻，比如适当引导学员、或者表达出你对学员回答的态度。\n"+
		"必须满足如下要求：\n"+
		"1. 当学员回复 “开始” 时，你要正式开始面试\n"+
		"2. 当学员表示希望 “结束面试” 时，你要结束面试\n"+
		"3. 此外，当你觉得这场面试可以结束时（比如候选人回答结果较差、不满足工作年限的招聘需求、或者候选人态度不礼貌），必须主动提出面试结束，不用继续询问更多问题了。并且要在回复中包含字符串【面试结束】\n"+
		"4. 面试结束后，应该给出候选人整场面试的表现和总结。\n"+
		"5. 使用纯文本回复", mockInterview.WorkExperience, mockInterview.JobPosition, mockInterview.Difficulty)

	switch req.Event {
	// 开始模拟面试
	case "start":
		// 定义用户 Prompt
		userPrompt := "开始"
		// 调用AI接口开始面试
		var chatMessages []*chatModel.ChatCompletionMessage
		// 添加系统预设
		chatMessages = append(chatMessages, &chatModel.ChatCompletionMessage{
			Role: chatModel.ChatMessageRoleSystem,
			Content: &chatModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(systemPrompt),
			},
		})
		// 添加用户预设
		chatMessages = append(chatMessages, &chatModel.ChatCompletionMessage{
			Role: chatModel.ChatMessageRoleUser,
			Content: &chatModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(userPrompt),
			},
		})
		result := ai.DoChat(chatMessages)
		// 将AI的回复序列化成json
		chatMessages = append(chatMessages, &chatModel.ChatCompletionMessage{
			Role: chatModel.ChatMessageRoleAssistant,
			Content: &chatModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(result),
			},
		})
		jsonStr, err := json.Marshal(chatMessages)
		if err != nil {
			return "", err
		}
		// 将AI的回复记录到数据库
		mockInterview.ID = req.ID
		mockInterview.Messages = string(jsonStr)
		mockInterview.Status = 1 // 进行中
		err = m.mockInterviewRepository.UpdateMockInterview(ctx, mockInterview)
		if err != nil {
			return "", err
		}
		return result, nil
	// 进行模拟面试
	case "chat":
		// 获取历史消息记录
		historyInterview, err := m.mockInterviewRepository.GetMockInterview(ctx, req.ID)
		if err != nil {
			return "", err
		}
		historyMessages := historyInterview.Messages
		// 将历史消息记录反序列化
		var historyChatMessages []model.MockInterviewMessage
		err = json.Unmarshal([]byte(historyMessages), &historyChatMessages)
		if err != nil {
			return "", err
		}
		var chatMessages []*chatModel.ChatCompletionMessage
		for _, message := range historyChatMessages {
			chatMessages = append(chatMessages, &chatModel.ChatCompletionMessage{
				Role: message.Role,
				Content: &chatModel.ChatCompletionMessageContent{
					StringValue: volcengine.String(message.Content),
				},
			})
		}
		// 定义用户 Prompt
		userPrompt := req.Message
		chatMessages = append(chatMessages, &chatModel.ChatCompletionMessage{
			Role: chatModel.ChatMessageRoleUser,
			Content: &chatModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(userPrompt),
			},
		})
		// 调用AI接口结束面试
		result := ai.DoChat(chatMessages)
		// 将AI的回复序列化成json
		chatMessages = append(chatMessages, &chatModel.ChatCompletionMessage{
			Role: chatModel.ChatMessageRoleAssistant,
			Content: &chatModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(result),
			},
		})
		jsonStr, err := json.Marshal(chatMessages)
		if err != nil {
			return "", err
		}
		// 将AI的回复记录到数据库
		if strings.Contains(result, "【面试结束】") {
			mockInterview.Status = 2 // 结束
		}
		mockInterview.ID = req.ID
		mockInterview.Messages = string(jsonStr)
		err = m.mockInterviewRepository.UpdateMockInterview(ctx, mockInterview)
		if err != nil {
			return "", err
		}
		return result, nil
	// 结束模拟面试
	case "end":
		// 获取历史消息记录
		historyInterview, err := m.mockInterviewRepository.GetMockInterview(ctx, req.ID)
		if err != nil {
			return "", err
		}
		historyMessages := historyInterview.Messages
		// 将历史消息记录反序列化
		var historyChatMessages []model.MockInterviewMessage
		err = json.Unmarshal([]byte(historyMessages), &historyChatMessages)
		if err != nil {
			return "", err
		}
		var chatMessages []*chatModel.ChatCompletionMessage
		for _, message := range historyChatMessages {
			chatMessages = append(chatMessages, &chatModel.ChatCompletionMessage{
				Role: message.Role,
				Content: &chatModel.ChatCompletionMessageContent{
					StringValue: volcengine.String(message.Content),
				},
			})
		}
		// 定义用户 Prompt
		userPrompt := "结束"
		chatMessages = append(chatMessages, &chatModel.ChatCompletionMessage{
			Role: chatModel.ChatMessageRoleUser,
			Content: &chatModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(userPrompt),
			},
		})
		// 调用AI接口结束面试
		result := ai.DoChat(chatMessages)
		// 将AI的回复序列化成json
		chatMessages = append(chatMessages, &chatModel.ChatCompletionMessage{
			Role: chatModel.ChatMessageRoleAssistant,
			Content: &chatModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(result),
			},
		})
		jsonStr, err := json.Marshal(chatMessages)
		if err != nil {
			return "", err
		}
		// 将AI的回复记录到数据库
		mockInterview.ID = req.ID
		mockInterview.Messages = string(jsonStr)
		mockInterview.Status = 2 // 结束
		err = m.mockInterviewRepository.UpdateMockInterview(ctx, mockInterview)
		if err != nil {
			return "", err
		}
		return result, nil
	}

	return "", nil
}
