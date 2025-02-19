package v1

import "time"

// 根据题库id查找题目

type QuestionBankQuestionQueryRequest struct {
	Current        *int    `json:"current,omitempty"`
	ID             *string `json:"id,omitempty"`
	NotID          *int    `json:"notId,omitempty"`
	PageSize       *int    `json:"pageSize,omitempty"`
	QuestionBankID *string `json:"questionBankId,omitempty"`
	QuestionID     *string `json:"questionId,omitempty"`
	SortField      *string `json:"sortField,omitempty"`
	SortOrder      *string `json:"sortOrder,omitempty"`
	UserID         *string `json:"userId,omitempty"`
}
type PageQuestionBankQuestionVO struct {
	CountID          *string                  `json:"countId,omitempty"`
	Current          *int                     `json:"current,omitempty"`
	MaxLimit         *int                     `json:"maxLimit,omitempty"`
	OptimizeCountSql *bool                    `json:"optimizeCountSql,omitempty"`
	Orders           []OrderItem              `json:"orders,omitempty"`
	Pages            *int                     `json:"pages,omitempty"`
	Records          []QuestionBankQuestionVO `json:"records,omitempty"`
	SearchCount      *bool                    `json:"searchCount,omitempty"`
	Size             *int                     `json:"size,omitempty"`
	Total            *int                     `json:"total,omitempty"`
}
type QuestionBankQuestionVO struct {
	CreateTime     *time.Time `json:"createTime,omitempty"`
	ID             *string    `json:"id,omitempty"`
	QuestionBankID *string    `json:"questionBankId,omitempty"`
	QuestionID     *string    `json:"questionId,omitempty"`
	TagList        []string   `json:"tagList,omitempty"`
	UpdateTime     *time.Time `json:"updateTime,omitempty"`
	User           *UserVO    `json:"user,omitempty"`
	UserID         *string    `json:"userId,omitempty"`
}

// 添加移除题目题库关系

type QuestionBankQuestionRequest struct {
	QuestionBankID *string `json:"questionBankId,omitempty"`
	QuestionID     *string `json:"questionId,omitempty"`
}

// 批量添加题目题库关系

type QuestionBankQuestionBatchRequest struct {
	QuestionBankID *string  `json:"questionBankId,omitempty"`
	QuestionIDList []string `json:"questionIdList,omitempty"`
}
