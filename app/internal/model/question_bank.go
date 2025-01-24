package model

import "time"

// QuestionBank 题库表
type QuestionBank struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement;comment:'id'"`                                 // 主键ID
	Title       *string   `gorm:"type:varchar(256);comment:'标题';index:idx_title"`                        // 标题
	Description *string   `gorm:"type:text;comment:'描述'"`                                                // 描述
	Picture     *string   `gorm:"type:varchar(2048);comment:'图片'"`                                       // 图片
	UserID      uint64    `gorm:"type:bigint;not null;comment:'创建用户id'"`                                 // 创建用户ID
	EditTime    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:'编辑时间'"`                // 编辑时间
	CreateTime  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:'创建时间'"`                // 创建时间
	UpdateTime  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:'更新时间'"` // 更新时间
	IsDelete    int8      `gorm:"type:tinyint;default:0;not null;comment:'是否删除'"`                        // 是否删除

	// 审核相关字段
	ReviewStatus  int        `gorm:"type:int;default:0;not null;comment:'状态：0-待审核, 1-通过, 2-拒绝'"` // 审核状态
	ReviewMessage *string    `gorm:"type:varchar(512);comment:'审核信息'"`                           // 审核信息
	ReviewerID    *uint64    `gorm:"type:bigint;comment:'审核人id'"`                                // 审核人ID
	ReviewTime    *time.Time `gorm:"type:datetime;comment:'审核时间'"`                               // 审核时间

	// 其他字段
	Priority int `gorm:"type:int;default:0;not null;comment:'优先级'"` // 优先级
	ViewNum  int `gorm:"type:int;default:0;not null;comment:'浏览量'"` // 浏览量
}

func (m *QuestionBank) TableName() string {
	return "question_bank"
}
