package handler

import (
	"app/api/v1"
	"app/internal/model"
	"app/internal/service"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func init() {
	// 注册 model.User 类型
	gob.Register(&model.User{})
}

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

// Register godoc
// @Summary 用户注册
// @Description 用于实现用户的注册功能
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param  request body v1.RegisterRequest true "注册请求参数"
// @Success 200 {object} v1.Response
// @Router /register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	req := new(v1.RegisterRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary 用户登录
// @Description 用于实现用户的登录功能
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest  true "注册请求参数"
// @Success 200 {object} v1.LoginResponseData
// @Router /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 获取 User-Agent
	userAgent := ctx.GetHeader("User-Agent")
	token, user, err := h.userService.Login(ctx, &req, userAgent)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	// 设置 session
	session := sessions.Default(ctx)
	session.Set("user_login", token)
	err = session.Save()
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.LoginResponseData{
		Id:          user.ID,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	})
}

// Logout godoc
// @Summary 用户注销
// @Description 用于实现用户的注销功能
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /logout [post]
func (h *UserHandler) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	// 获取用户信息
	t := session.Get("user_login")
	if t == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	token := t.(string)
	// 删除 session
	session.Delete("user_login")
	err := session.Save()
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	// 获取 User-Agent
	userAgent := ctx.GetHeader("User-Agent")
	ok, err := h.userService.Logout(ctx, token, userAgent)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, ok)
}

// GetLoginUser godoc
// @Summary 获取用户登录状态
// @Description 用于获取全局用户的登录状态
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200 {object} v1.GetLoginUserResponseData
// @Router /get/login [post]
func (h *UserHandler) GetLoginUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	t := session.Get("user_login")
	if t == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	token := t.(string)
	// 获取 User-Agent
	userAgent := ctx.GetHeader("User-Agent")
	user, err := h.userService.GetLoginUser(ctx, token, userAgent)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.GetLoginUserResponseData{
		Id:          user.ID,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	})
}

// ListPage godoc
// @Summary 用户列表
// @Description 用于获取所有用户
// @Tags 用户模块（管理员）
// @Accept json
// @Produce json
// @Param  request body v1.UserQueryRequest  true "注册请求参数"
// @Success 200 {object} v1.UserQueryResponseData[v1.User]
// @Router /list/page [post]
func (h *UserHandler) ListPage(ctx *gin.Context) {
	var req v1.UserQueryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	page, err := h.userService.ListUserByPage(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.UserQueryResponseData[v1.User]{
		Records: page.Records,
		Size:    page.Size,
		Total:   page.Total,
		Current: page.Current,
		Pages:   page.Pages,
	})
}

// AddUser godoc
// @Summary 添加用户
// @Description 用于添加新用户
// @Tags 用户模块（管理员）
// @Accept json
// @Produce json
// @Param  request body v1.AddUserRequest  true "注册请求参数"
// @Success 200 {object} uint64
// @Router /add [post]
func (h *UserHandler) AddUser(ctx *gin.Context) {
	var req v1.AddUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	id, err := h.userService.AddUser(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, id)
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 用于删除用户
// @Tags 用户模块（管理员）
// @Accept json
// @Produce json
// @Param  request body v1.DeleteUserRequest  true "注册请求参数"
// @Success 200 {object} bool
// @Router /delete [post]
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	var req v1.DeleteUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.userService.DeleteUser(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, ok)
}

// UpdateUser godoc
// @Summary 修改用户
// @Description 用于修改用户信息
// @Tags 用户模块（管理员）
// @Accept json
// @Produce json
// @Param request body v1.UpdateUserRequest  true "注册请求参数"
// @Success 200 {object} bool
// @Router /update [post]
func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var req v1.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	ok, err := h.userService.UpdateUser(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, ok)
}

// AddUserSignIn godoc
// @Summary 用户签到
// @Description 用于用户签到
// @Tags 用户模块（管理员）
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /add/sign_in [post]
func (h *UserHandler) AddUserSignIn(ctx *gin.Context) {
	// 必须要登录才能签到
	session := sessions.Default(ctx)
	t := session.Get("user_login")
	if t == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	token := t.(string)
	if token == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}

	result, err := h.userService.AddUserSignIn(ctx, token)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, result)
}

// GetUserSignIn godoc
// @Summary 签到日历
// @Description 用于用户查看签到日历
// @Tags 用户模块（管理员）
// @Accept json
// @Produce json
// @Success 200 {object} []int
// @Router /get/sign_in [get]
func (h *UserHandler) GetUserSignIn(ctx *gin.Context) {
	// 必须要登录才能获取签到记录
	session := sessions.Default(ctx)
	t := session.Get("user_login")
	if t == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	token := t.(string)
	if token == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	var req v1.GetUserSignInRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	var year int
	if req.Year == nil {
		date := time.Now()
		year = date.Year()
	} else {
		year = *req.Year
	}
	var dayList []int
	dayList, err := h.userService.GetUserSignIn(ctx, token, year)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, dayList)
}
