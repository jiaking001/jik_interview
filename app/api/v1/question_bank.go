package v1

import "time"

// 分页查询

type QuestionBank struct {
	CreateTime  *time.Time `json:"createTime,omitempty"`
	Description *string    `json:"description,omitempty"`
	EditTime    *time.Time `json:"editTime,omitempty"`
	ID          *string    `json:"id,omitempty"`
	IsDelete    *int8      `json:"isDelete,omitempty"`
	Picture     *string    `json:"picture,omitempty"`
	Title       *string    `json:"title,omitempty"`
	UpdateTime  *time.Time `json:"updateTime,omitempty"`
	UserID      *string    `json:"userId,omitempty"`
}
type QuestionBankRequest struct {
	Current               *int    `json:"current,omitempty"`               // 当前页码
	Description           *string `json:"description,omitempty"`           // 描述
	ID                    *string `json:"id,omitempty"`                    // ID
	NeedQueryQuestionList *bool   `json:"needQueryQuestionList,omitempty"` // 是否需要查询问题列表
	NotID                 *int    `json:"notId,omitempty"`                 // 排除的 ID
	PageSize              *int    `json:"pageSize,omitempty"`              // 每页大小
	Picture               *string `json:"picture,omitempty"`               // 图片链接
	SearchText            *string `json:"searchText,omitempty"`            // 搜索文本
	SortField             *string `json:"sortField,omitempty"`             // 排序字段
	SortOrder             *string `json:"sortOrder,omitempty"`             // 排序顺序
	Title                 *string `json:"title,omitempty"`                 // 标题
	UserID                *string `json:"userId,omitempty"`                // 用户 ID
}
type QuestionBankQueryResponseData[T any] struct {
	Records []T  `json:"records"` // 当前页的记录列表
	Total   *int `json:"total"`   // 总记录数
	Size    *int `json:"size"`    // 每页大小
	Current *int `json:"current"` // 当前页码
	Pages   *int `json:"pages"`   // 总页数
}

// 添加题库

type AddQuestionBankRequest struct {
	Description *string `json:"description,omitempty"` // 描述
	Picture     *string `json:"picture,omitempty"`     // 图片链接
	Title       *string `json:"title,omitempty"`       // 标题
}

// 删除题库

type DeleteQuestionBankRequest struct {
	Id string `json:"id"`
}

// 更新题库

type UpdateQuestionBankRequest struct {
	Description *string `json:"description,omitempty"` // 描述
	ID          *string `json:"id,omitempty"`          // ID
	Picture     *string `json:"picture,omitempty"`     // 图片链接
	Title       *string `json:"title,omitempty"`       // 标题
}

// 获取题库详情

type GetQuestionBankRequest struct {
	Current               *int    `form:"current,omitempty"`               // 当前页码
	Description           *string `form:"description,omitempty"`           // 描述
	ID                    *string `form:"id,omitempty"`                    // ID
	NeedQueryQuestionList *bool   `form:"needQueryQuestionList,omitempty"` // 是否需要查询问题列表
	NotID                 *string `form:"notId,omitempty"`                 // 排除的 ID
	PageSize              *int    `form:"pageSize,omitempty"`              // 每页大小
	Picture               *string `form:"picture,omitempty"`               // 图片链接
	SearchText            *string `form:"searchText,omitempty"`            // 搜索文本
	SortField             *string `form:"sortField,omitempty"`             // 排序字段
	SortOrder             *string `form:"sortOrder,omitempty"`             // 排序顺序
	Title                 *string `form:"title,omitempty"`                 // 标题
	UserID                *string `form:"userId,omitempty"`                // 用户 ID
}
type GetQuestionBankResponse struct {
	CreateTime   *time.Time      `json:"createTime,omitempty"`   // 创建时间
	Description  *string         `json:"description,omitempty"`  // 描述
	ID           *string         `json:"id,omitempty"`           // ID
	Picture      *string         `json:"picture,omitempty"`      // 图片链接
	QuestionPage *PageQuestionVO `json:"questionPage,omitempty"` // 问题分页信息
	Title        *string         `json:"title,omitempty"`        // 标题
	UpdateTime   *time.Time      `json:"updateTime,omitempty"`   // 更新时间
	User         *UserVO         `json:"user,omitempty"`         // 用户信息
	UserID       *string         `json:"userId,omitempty"`       // 用户 ID
}
type UserVO struct {
	CreateTime  *time.Time `json:"createTime,omitempty"`  // 创建时间
	ID          *int       `json:"id,omitempty"`          // 用户 ID
	UserAvatar  *string    `json:"userAvatar,omitempty"`  // 用户头像
	UserName    *string    `json:"userName,omitempty"`    // 用户名
	UserProfile *string    `json:"userProfile,omitempty"` // 用户简介
	UserRole    *string    `json:"userRole,omitempty"`    // 用户角色
}
type PageQuestionVO struct {
	CountId          *string      `json:"countId,omitempty"`          // 计数 ID
	Current          *int         `json:"current,omitempty"`          // 当前页码
	MaxLimit         *int         `json:"maxLimit,omitempty"`         // 最大限制
	OptimizeCountSql *bool        `json:"optimizeCountSql,omitempty"` // 是否优化计数 SQL
	Orders           []OrderItem  `json:"orders,omitempty"`           // 排序项
	Pages            *int         `json:"pages,omitempty"`            // 总页数
	Records          []QuestionVO `json:"records,omitempty"`          // 问题记录
	SearchCount      *bool        `json:"searchCount,omitempty"`      // 是否搜索计数
	Size             *int         `json:"size,omitempty"`             // 每页大小
	Total            *int         `json:"total,omitempty"`            // 总记录数
}
type OrderItem struct {
	Asc    *bool   `json:"asc,omitempty"`    // 是否升序
	Column *string `json:"column,omitempty"` // 排序列
}
type QuestionVO struct {
	Answer     *string    `json:"answer,omitempty"`     // 回答内容
	Content    *string    `json:"content,omitempty"`    // 问题内容
	CreateTime *time.Time `json:"createTime,omitempty"` // 创建时间
	ID         *string    `json:"id,omitempty"`         // 问题 ID
	TagList    []string   `json:"tagList,omitempty"`    // 标签列表
	Title      *string    `json:"title,omitempty"`      // 问题标题
	UpdateTime *time.Time `json:"updateTime,omitempty"` // 更新时间
	User       *UserVO    `json:"user,omitempty"`       // 用户信息
	UserID     *string    `json:"userId,omitempty"`     // 用户 ID
}
