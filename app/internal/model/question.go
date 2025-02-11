package model

import (
	"gorm.io/gorm"
	"time"
)

// Question 题目表
type Question struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement;comment:'id'"`                                 // 主键ID
	Title      *string        `gorm:"type:varchar(256);comment:'标题';index:idx_title"`                        // 标题
	Content    *string        `gorm:"type:text;comment:'内容'"`                                                // 内容
	Tags       *string        `gorm:"type:varchar(1024);comment:'标签列表（json 数组）'"`                            // 标签列表（JSON数组）
	Answer     *string        `gorm:"type:text;comment:'推荐答案'"`                                              // 推荐答案
	UserID     uint64         `gorm:"type:bigint;not null;comment:'创建用户id';index:idx_userId"`                // 创建用户ID
	EditTime   time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:'编辑时间'"`                // 编辑时间
	CreateTime time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:'创建时间'"`                // 创建时间
	UpdateTime time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:'更新时间'"` // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"index;comment:'删除时间'"`
	IsDelete   int8           `gorm:"type:tinyint;default:0;not null;comment:'是否删除'"` // 是否删除

	// 审核相关字段
	ReviewStatus  int        `gorm:"type:int;default:0;not null;comment:'状态：0-待审核, 1-通过, 2-拒绝'"` // 审核状态
	ReviewMessage *string    `gorm:"type:varchar(512);comment:'审核信息'"`                           // 审核信息
	ReviewerID    *uint64    `gorm:"type:bigint;comment:'审核人id'"`                                // 审核人ID
	ReviewTime    *time.Time `gorm:"type:datetime;comment:'审核时间'"`                               // 审核时间

	// 互动相关字段
	ViewNum   int `gorm:"type:int;default:0;not null;comment:'浏览量'"` // 浏览量
	ThumbNum  int `gorm:"type:int;default:0;not null;comment:'点赞数'"` // 点赞数
	FavourNum int `gorm:"type:int;default:0;not null;comment:'收藏数'"` // 收藏数

	// 排序相关字段
	Priority int `gorm:"type:int;default:0;not null;comment:'优先级'"` // 优先级

	// 新增字段
	Source  *string `gorm:"type:varchar(512);comment:'题目来源'"`                           // 题目来源
	NeedVip int8    `gorm:"type:tinyint;default:0;not null;comment:'仅会员可见（1 表示仅会员可见）'"` // 仅会员可见
}

func (m *Question) TableName() string {
	return "question"
}

type QuestionEs struct {
	Id         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Tags       []string  `json:"tags"`
	Answer     string    `json:"answer"`
	UserId     int64     `json:"user_id"`
	EditTime   time.Time `json:"edit_time"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	IsDelete   int8      `json:"is_delete"`
}
