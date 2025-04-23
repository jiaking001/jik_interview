package v1

import (
	"time"
)

// MockInterviewEventRequest 模拟面试事件请求
type MockInterviewEventRequest struct {
	Event   string `json:"event,omitempty"`   // 事件类型
	ID      uint64 `json:"id,omitempty"`      // 事件关联的ID
	Message string `json:"message,omitempty"` // 事件消息内容
}

// MockInterview 模拟面试信息
type MockInterview struct {
	CreateTime     time.Time `json:"createTime,omitempty"`     // 创建时间
	Difficulty     string    `json:"difficulty,omitempty"`     // 面试难度
	ID             uint64    `json:"id,omitempty"`             // 面试 ID
	IsDelete       int8      `json:"isDelete,omitempty"`       // 是否删除
	JobPosition    string    `json:"jobPosition,omitempty"`    // 工作岗位
	Messages       string    `json:"messages,omitempty"`       // 消息列表（JSON 对象数组字段，同时包括了总结）
	Status         int       `json:"status,omitempty"`         // 状态（0-待开始、1-进行中、2-已结束）
	UpdateTime     time.Time `json:"updateTime,omitempty"`     // 更新时间
	UserID         uint64    `json:"userId,omitempty"`         // 创建人（用户 ID）
	WorkExperience string    `json:"workExperience,omitempty"` // 工作年限
}

// MockInterviewAddRequest 模拟面试添加请求
type MockInterviewAddRequest struct {
	Difficulty     string `json:"difficulty,omitempty"`     // 面试难度
	JobPosition    string `json:"jobPosition,omitempty"`    // 工作岗位
	WorkExperience string `json:"workExperience,omitempty"` // 工作年限
}

// MockInterviewGetRequest 获取模拟面试信息
type MockInterviewGetRequest struct {
	ID uint64 `json:"id,omitempty"` // 面试 ID
}
