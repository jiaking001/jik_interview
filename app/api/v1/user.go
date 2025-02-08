package v1

import "time"

// 注册

type RegisterRequest struct {
	UserAccount   string `json:"userAccount" binding:"required" example:"1234456"`
	UserPassword  string `json:"userPassword" binding:"required" example:"123456"`
	CheckPassword string `json:"checkPassword" binding:"required" example:"123456"`
}

// 登录

type LoginRequest struct {
	UserAccount  string `json:"userAccount" binding:"required" example:"123456"`
	UserPassword string `json:"userPassword" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	Id          uint64    `json:"id"`
	UserName    *string   `json:"userName"`
	UserAvatar  *string   `json:"userAvatar"`
	UserProfile *string   `json:"userProfile"`
	UserRole    string    `json:"userRole"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

// 获取全局用户状态

type GetLoginUserResponseData struct {
	Id          uint64    `json:"id"`
	UserName    *string   `json:"userName"`
	UserAvatar  *string   `json:"userAvatar"`
	UserProfile *string   `json:"userProfile"`
	UserRole    string    `json:"userRole"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

// 分页查询

type User struct {
	CreateTime   *time.Time `json:"createTime,omitempty"`
	EditTime     *time.Time `json:"editTime,omitempty"`
	ID           *string    `json:"id,omitempty"`
	IsDelete     *int8      `json:"isDelete,omitempty"`
	MpOpenID     *string    `json:"mpOpenId,omitempty"`
	UnionID      *string    `json:"unionId,omitempty"`
	UpdateTime   *time.Time `json:"updateTime,omitempty"`
	UserAccount  *string    `json:"userAccount,omitempty"`
	UserAvatar   *string    `json:"userAvatar,omitempty"`
	UserName     *string    `json:"userName,omitempty"`
	UserPassword *string    `json:"userPassword,omitempty"`
	UserProfile  *string    `json:"userProfile,omitempty"`
	UserRole     *string    `json:"userRole,omitempty"`
}
type UserQueryRequest struct {
	Current     *int    `json:"current"`
	ID          *string `json:"id"`
	MpOpenID    *string `json:"mpOpenId"`
	PageSize    *int    `json:"pageSize"`
	SortField   *string `json:"sortField"`
	SortOrder   *string `json:"sortOrder"`
	UnionID     *string `json:"unionId"`
	UserName    *string `json:"userName"`
	UserAccount *string `json:"userAccount"`
	UserProfile *string `json:"userProfile"`
	UserRole    *string `json:"userRole"`
}
type PageResult[T any] struct {
	Records []T  `json:"records"` // 当前页的记录列表
	Total   *int `json:"total"`   // 总记录数
	Size    *int `json:"size"`    // 每页大小
	Current *int `json:"current"` // 当前页码
	Pages   *int `json:"pages"`   // 总页数
}
type UserQueryResponseData[T any] struct {
	Records []T  `json:"records"` // 当前页的记录列表
	Total   *int `json:"total"`   // 总记录数
	Size    *int `json:"size"`    // 每页大小
	Current *int `json:"current"` // 当前页码
	Pages   *int `json:"pages"`   // 总页数
}

// 添加用户

type AddUserRequest struct {
	UserAccount *string `json:"userAccount"`
	UserAvatar  *string `json:"userAvatar"`
	UserName    *string `json:"userName"`
	UserProfile *string `json:"userProfile"`
	UserRole    *string `json:"userRole"`
}

// 删除用户

type DeleteUserRequest struct {
	Id string `json:"id"`
}

// 更新用户

type UpdateUserRequest struct {
	Id          string  `json:"id"`
	UserAccount *string `json:"userAccount"`
	UserAvatar  *string `json:"userAvatar"`
	UserName    *string `json:"userName"`
	UserRole    *string `json:"userRole"`
	UserProfile *string `json:"userProfile"`
}

type UserVO struct {
	CreateTime  *time.Time `json:"createTime,omitempty"`  // 创建时间
	ID          *int       `json:"id,omitempty"`          // 用户 ID
	UserAvatar  *string    `json:"userAvatar,omitempty"`  // 用户头像
	UserName    *string    `json:"userName,omitempty"`    // 用户名
	UserProfile *string    `json:"userProfile,omitempty"` // 用户简介
	UserRole    *string    `json:"userRole,omitempty"`    // 用户角色
}
