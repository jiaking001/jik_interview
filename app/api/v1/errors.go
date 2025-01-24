package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrAccountAlreadyUse     = newError(1001, "The account is already in use.")
	ErrInconsistentPasswords = newError(1002, "The password is inconsistent.")
)
