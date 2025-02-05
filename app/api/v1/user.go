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
	Id          uint64
	UserName    *string
	UserAvatar  *string
	UserProfile *string
	UserRole    string
	CreateTime  time.Time
	UpdateTime  time.Time
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

// 获取全局用户状态

type GetLoginUserResponseData struct {
	Id          uint64
	UserName    *string
	UserAvatar  *string
	UserProfile *string
	UserRole    string
	CreateTime  time.Time
	UpdateTime  time.Time
}
type GetLoginUserResponse struct {
	Response
	Data GetLoginUserResponseData
}
