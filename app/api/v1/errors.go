package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// user
	ErrIllegalRole           = newError(40000, "未设置用户权限")
	ErrAccountAlreadyUse     = newError(40000, "账号已被注册")
	ErrInconsistentPasswords = newError(40000, "两次密码输入不一致")
	ErrIllegalPassword       = newError(40000, "密码不规范")
	ErrIllegalAccount        = newError(40000, "账号不规范")
	ErrPassword              = newError(40000, "账号或密码错误")
	ParamsError              = newError(40000, "请求参数错误")
	NotLoginError            = newError(40100, "未登录")
	ErrBanRole               = newError(40000, "账号因违规行为已被封禁")

	// questionBank
	ErrTitleAlreadyUse = newError(40000, "题库或题目已存在")

	ErrSystemIsBusy = newError(40000, "系统繁忙,请稍后再试")

	ErrBotLogin = newError(40000, "爬虫用户，拒绝访问")
)
