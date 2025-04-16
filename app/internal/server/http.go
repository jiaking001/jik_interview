package server

import (
	v1 "app/api/v1"
	"app/docs"
	"app/internal/handler"
	"app/internal/middleware"
	"app/pkg/jwt"
	"app/pkg/log"
	"app/pkg/server/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	sentinelPlugin "github.com/sentinel-group/sentinel-go-adapters/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	hp "net/http"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	rdb *redis.Client,
	db *gorm.DB,
	userHandler *handler.UserHandler,
	questionHandler *handler.QuestionHandler,
	questionBankHandler *handler.QuestionBankHandler,
	questionBankQuestionHandler *handler.QuestionBankQuestionHandler,
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
		// 集成 Sentinel
		sentinelPlugin.SentinelMiddleware(
			// 自定义降级逻辑
			sentinelPlugin.WithBlockFallback(func(ctx *gin.Context) {
				v1.HandleError(ctx, hp.StatusUnauthorized, v1.ErrSystemIsBusy, nil)
				ctx.Abort()
			}),
		),
		// 黑名单
		middleware.BlacklistMiddleware(),
		// 判断用户是否登录
		// middleware.GetLoginStatus(jwt, rdb),
	)

	r := s.Group("/api")
	{
		// No route group has permission
		noAuthRouter := r.Group("/")
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
			user.POST("/add/sign_in", middleware.GetLoginStatus(jwt, rdb), middleware.AntiCrawling(jwt, rdb, db), userHandler.AddUserSignIn)
			user.GET("/get/sign_in", middleware.GetLoginStatus(jwt, rdb), userHandler.GetUserSignIn)

			// 题库模块
			questionBank := noAuthRouter.Group("/questionBank")
			questionBank.POST("/list/page", questionBankHandler.ListPage)
			questionBank.POST("/list/page/vo", questionBankHandler.ListPage)
			questionBank.POST("/add", questionBankHandler.AddQuestionBank)
			questionBank.POST("/delete", questionBankHandler.DeleteQuestionBank)
			questionBank.POST("/update", questionBankHandler.UpdateQuestionBank)
			questionBank.GET("/get/vo", questionBankHandler.GetQuestionBank)

			// 题目模块
			question := noAuthRouter.Group("/question")
			question.POST("/list/page", questionHandler.ListPage)
			question.POST("/list/page/vo", questionHandler.ListPageVo)
			question.POST("/add", questionHandler.AddQuestion)
			question.POST("/delete", questionHandler.DeleteQuestion)
			question.POST("/update", questionHandler.UpdateQuestion)
			question.GET("/get/vo", questionHandler.GetQuestion)
			question.POST("/search/page/vo", questionHandler.SearchPageVo)
			question.POST("/delete/batch", questionHandler.DeleteBatchQuestion)
			question.POST("/ai/generate/question", questionHandler.AiGenerateQuestion)

			// 题目题库模块
			questionBankQuestion := noAuthRouter.Group("/questionBankQuestion")
			questionBankQuestion.POST("/list/page/vo", questionBankQuestionHandler.GetQuestionBankQuestion)
			questionBankQuestion.POST("/add", questionBankQuestionHandler.AddQuestionBankQuestion)
			questionBankQuestion.POST("/remove", questionBankQuestionHandler.RemoveQuestionBankQuestion)
			questionBankQuestion.POST("/add/batch", questionBankQuestionHandler.BatchAddQuestionBankQuestion)
			questionBankQuestion.POST("/remove/batch", questionBankQuestionHandler.BatchRemoveQuestionBankQuestion)
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
