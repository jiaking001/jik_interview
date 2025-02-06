package server

import (
	"app/docs"
	"app/internal/handler"
	"app/internal/middleware"
	"app/pkg/jwt"
	"app/pkg/log"
	"app/pkg/server/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	// questionHandler *handler.QuestionHandler,
	// questionBankHandler *handler.QuestionBankHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	// 创建基于 Cookie 的存储引擎，"secret" 是用于加密的密钥
	store := cookie.NewStore([]byte("jik"))
	s.Use(
		sessions.Sessions("user", store),
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	v1 := s.Group("/api")
	{
		// No route group has permission
		noAuthRouter := v1.Group("/")
		{
			// 用户模块
			user := noAuthRouter.Group("/user")
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)
			user.GET("/get/login", userHandler.GetLoginUser)
			user.POST("/logout", userHandler.Logout)
			user.POST("/list/page", userHandler.ListPage)
			user.POST("/add", userHandler.AddUser)
			user.POST("/delete", userHandler.DeleteUser)
			user.POST("/update", userHandler.UpdateUser)
		}
		// Non-strict permission routing group
		//noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
		//{
		//	//noStrictAuthRouter.GET("/user", userHandler.GetProfile)
		//}

		// Strict permission routing group
		//strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
		//{
		//	//strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
		//}
	}

	return s
}
