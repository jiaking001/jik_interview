package model

import (
	"app/internal/middleware"
	"gorm.io/gorm"
	"time"
)

// User 用户表
type User struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement;comment:'id'"`                                   // 主键ID
	UserAccount  string         `gorm:"type:varchar(256);not null;comment:'账号'"`                                 // 账号
	UserPassword string         `gorm:"type:varchar(512);not null;comment:'密码'"`                                 // 密码
	UnionId      *string        `gorm:"type:varchar(256);comment:'微信开放平台id';index:idx_unionId"`                  // 微信开放平台id
	MpOpenId     *string        `gorm:"type:varchar(256);comment:'公众号openId'"`                                   // 公众号openId
	UserName     *string        `gorm:"type:varchar(256);comment:'用户昵称'"`                                        // 用户昵称
	UserAvatar   *string        `gorm:"type:varchar(1024);comment:'用户头像'"`                                       // 用户头像
	UserProfile  *string        `gorm:"type:varchar(512);comment:'用户简介'"`                                        // 用户简介
	UserRole     string         `gorm:"type:varchar(256);default:'user';not null;comment:'用户角色：user/admin/ban'"` // 用户角色
	EditTime     time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:'编辑时间'"`                  // 编辑时间
	CreateTime   time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:'创建时间'"`                  // 创建时间
	UpdateTime   time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:'更新时间'"`   // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index;comment:'删除时间'"`
	IsDelete     int8           `gorm:"type:tinyint;default:0;not null;comment:'是否删除'"` // 是否删除

	// 新增字段
	VipExpireTime *time.Time `gorm:"type:datetime;comment:'会员过期时间'"`    // 会员过期时间
	VipCode       *string    `gorm:"type:varchar(128);comment:'会员兑换码'"` // 会员兑换码
	VipNumber     *uint64    `gorm:"type:bigint;comment:'会员编号'"`        // 会员编号

	// 最新新增字段
	ShareCode  *string `gorm:"type:varchar(20);default:null;comment:'分享码'"` // 分享码
	InviteUser *uint64 `gorm:"type:bigint;default:null;comment:'邀请用户id'"`   // 邀请用户id
}

func (u *User) TableName() string {
	return "users"
}

// BeforeCreate 在创建记录前生成雪花算法 ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uint64(middleware.SnowflakeNode.Generate().Int64())
	return nil
}
