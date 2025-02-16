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
// @Schemes
// @Description 目前只支持邮箱登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
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
// @Summary 账号登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.LoginResponse
// @Router /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	_, user, err := h.userService.Login(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	session := sessions.Default(ctx)
	session.Set("user_login", user)
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
// @Summary 退出登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200
// @Router /logout [post]
func (h *UserHandler) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("user_login")
	err := session.Save()
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, true)
}

// GetLoginUser godoc
// @Summary 获取用户登录状态
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.GetLoginUserRequest true "params"
// @Success 200 {object} v1.GetLoginUserResponse
// @Router /get/login [post]
func (h *UserHandler) GetLoginUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userInterface := session.Get("user_login")
	if userInterface == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	user := userInterface.(*model.User)
	if user == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
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

func (h *UserHandler) AddUserSignIn(ctx *gin.Context) {
	// 必须要登录才能签到
	session := sessions.Default(ctx)
	userInterface := session.Get("user_login")
	if userInterface == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	user := userInterface.(*model.User)
	result, err := h.userService.AddUserSignIn(ctx, user.ID)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, result)
}

func (h *UserHandler) GetUserSignIn(ctx *gin.Context) {
	// 必须要登录才能获取签到记录
	session := sessions.Default(ctx)
	userInterface := session.Get("user_login")
	if userInterface == nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.NotLoginError, nil)
		return
	}
	user := userInterface.(*model.User)
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
	dayList, err := h.userService.GetUserSignIn(ctx, user.ID, year)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	v1.HandleSuccess(ctx, dayList)
}
