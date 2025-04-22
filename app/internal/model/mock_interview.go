package model

import (
	"time"
)

// MockInterview 模拟面试表
type MockInterview struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement;comment:'id'"`                                 // 主键ID
	WorkExperience string    `gorm:"type:varchar(256);not null;comment:'工作年限'"`                             // 工作年限
	JobPosition    string    `gorm:"type:varchar(256);not null;comment:'工作岗位'"`                             // 工作岗位
	Difficulty     string    `gorm:"type:varchar(50);not null;comment:'面试难度'"`                              // 面试难度
	Messages       *string   `gorm:"type:mediumtext;comment:'消息列表（JSON 对象数组字段，同时包括了总结）'"`                   // 消息列表
	Status         int       `gorm:"type:int;default:0;not null;comment:'状态（0-待开始、1-进行中、2-已结束）'"`           // 状态
	UserID         uint64    `gorm:"type:bigint;not null;comment:'创建人（用户 id）';index:idx_userId"`            // 创建人用户ID
	CreateTime     time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:'创建时间'"`                // 创建时间
	UpdateTime     time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:'更新时间'"` // 更新时间
	IsDelete       int8      `gorm:"type:tinyint;default:0;not null;comment:'是否删除（逻辑删除）'"`                  // 是否删除
}

func (m *MockInterview) TableName() string {
	return "mock_interview"
}
