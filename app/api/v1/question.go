package v1

import "time"

// 分页查询

type Question struct {
	Answer     *string    `json:"answer,omitempty"`     // 回答内容
	Content    *string    `json:"content,omitempty"`    // 问题内容
	CreateTime *time.Time `json:"createTime,omitempty"` // 创建时间
	EditTime   *time.Time `json:"editTime,omitempty"`   // 编辑时间
	ID         *string    `json:"id,omitempty"`         // 问题 ID
	IsDelete   *int8      `json:"isDelete,omitempty"`   // 是否删除
	Tags       *string    `json:"tags,omitempty"`       // 标签
	Title      *string    `json:"title,omitempty"`      // 问题标题
	UpdateTime *time.Time `json:"updateTime,omitempty"` // 更新时间
	UserID     *string    `json:"userId,omitempty"`     // 用户 ID
}
type QuestionRequest struct {
	Answer         *string  `json:"answer,omitempty"`         // 回答内容
	Content        *string  `json:"content,omitempty"`        // 问题内容
	Current        *int     `json:"current,omitempty"`        // 当前页码
	ID             *string  `json:"id,omitempty"`             // 问题 ID
	NotID          *int     `json:"notId,omitempty"`          // 排除的 ID
	PageSize       *int     `json:"pageSize,omitempty"`       // 每页大小
	QuestionBankID *string  `json:"questionBankId,omitempty"` // 题库 ID
	SearchText     *string  `json:"searchText,omitempty"`     // 搜索文本
	SortField      *string  `json:"sortField,omitempty"`      // 排序字段
	SortOrder      *string  `json:"sortOrder,omitempty"`      // 排序顺序
	Tags           []string `json:"tag,omitempty"`            // 标签列表
	Title          *string  `json:"title,omitempty"`          // 问题标题
	UserID         *string  `json:"userId,omitempty"`         // 用户 ID
}
type QuestionQueryResponseData[T any] struct {
	Records []T  `json:"records"` // 当前页的记录列表
	Total   *int `json:"total"`   // 总记录数
	Size    *int `json:"size"`    // 每页大小
	Current *int `json:"current"` // 当前页码
	Pages   *int `json:"pages"`   // 总页数
}

// 添加题目

type AddQuestionRequest struct {
	Answer  *string  `json:"answer,omitempty"`  // 回答内容
	Content *string  `json:"content,omitempty"` // 内容
	Tags    []string `json:"tags,omitempty"`    // 标签列表
	Title   *string  `json:"title,omitempty"`   // 标题
}

// 删除题目

type DeleteQuestionRequest struct {
	Id string `json:"id"`
}

// 更新题目

type UpdateQuestionRequest struct {
	Answer  *string  `json:"answer,omitempty"`  // 回答内容
	Content *string  `json:"content,omitempty"` // 内容
	ID      *string  `json:"id,omitempty"`      // ID
	Tags    []string `json:"tags,omitempty"`    // 标签列表
	Title   *string  `json:"title,omitempty"`   // 标题
}

// 获取题目详情

type GetQuestionRequest struct {
	ID *string `form:"id,omitempty"` // ID
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

type SearchResult struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source Question `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
