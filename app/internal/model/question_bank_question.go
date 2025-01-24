package model

import (
	"time"
)

// QuestionBankQuestion 题库题目关系表
type QuestionBankQuestion struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement;comment:'id'"`                                          // 主键ID
	QuestionBankID uint64    `gorm:"type:bigint;not null;comment:'题库 id';uniqueIndex:question_bank_question_unique"` // 题库ID
	QuestionID     uint64    `gorm:"type:bigint;not null;comment:'题目 id';uniqueIndex:question_bank_question_unique"` // 题目ID
	UserID         uint64    `gorm:"type:bigint;not null;comment:'创建用户 id'"`                                         // 创建用户ID
	CreateTime     time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:'创建时间'"`                         // 创建时间
	UpdateTime     time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:'更新时间'"`          // 更新时间

	// 新增字段
	QuestionOrder int `gorm:"type:int;default:0;not null;comment:'题目顺序（题号）'"` // 题目顺序（题号）
}

func (m *QuestionBankQuestion) TableName() string {
	return "question_bank_question"
}
