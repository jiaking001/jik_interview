package v1

type RegisterRequest struct {
	UserAccount   string `json:"userAccount" binding:"required" example:"1234456"`
	UserPassword  string `json:"userPassword" binding:"required" example:"123456"`
	CheckPassword string `json:"checkPassword" binding:"required" example:"123456"`
}

type LoginRequest struct {
	UserAccount  string `json:"userAccount" binding:"required" example:"123456"`
	UserPassword string `json:"userPassword" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}
