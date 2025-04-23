//go:build wireinject
// +build wireinject

package wire

import (
	"app/internal/handler"
	"app/internal/job"
	"app/internal/repository"
	"app/internal/server"
	"app/internal/service"
	"app/pkg/app"
	"app/pkg/jwt"
	"app/pkg/log"
	"app/pkg/server/http"
	"app/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewElasticsearch,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewQuestionBankRepository,
	repository.NewQuestionRepository,
	repository.NewQuestionBankQuestionRepository,
	repository.NewMockInterviewRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewQuestionBankService,
	service.NewQuestionService,
	service.NewQuestionBankQuestionService,
	service.NewMockInterviewService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewQuestionBankHandler,
	handler.NewQuestionHandler,
	handler.NewQuestionBankQuestionHandler,
	handler.NewMockInterviewHandler,
)

var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
	job.NewQuestionJob,
)
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJobServer,
)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		jobSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
