package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	userv1 "github.com/gorobot-nz/test-task/gen/proto/user/v1"

	usershandler "github.com/gorobot-nz/test-task/internal/handler/grpc/users"
	usersrepository "github.com/gorobot-nz/test-task/internal/repository/users"
	usersservice "github.com/gorobot-nz/test-task/internal/service/users"
	applogger "github.com/gorobot-nz/test-task/pkg/logger"
	"github.com/gorobot-nz/test-task/pkg/storage"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type App struct {
	logger *zap.Logger

	s *grpc.Server
}

func NewApp() *App {
	logger := applogger.NewLogger(zapcore.DebugLevel)
	store := storage.NewStorage[userv1.User]()

	repository := usersrepository.NewStorageRepository(logger.Named("UsersRepository"), store)
	service := usersservice.NewService(logger.Named("UsersService"), repository)
	handler := usershandler.NewHandler(logger.Named("UsersHandler"), service)

	server := grpc.NewServer(grpc.UnaryInterceptor(
		grpcmiddleware.ChainUnaryServer(
			grpczap.UnaryServerInterceptor(logger),
		),
	))

	userv1.RegisterUserServiceServer(server, handler)

	return &App{
		s:      server,
		logger: logger,
	}
}

func (a *App) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", "8000"))

	if err != nil {
		a.logger.Fatal("Failed listen", zap.Error(err))
	}

	go func() {
		a.logger.Fatal("Failed to serve gRPC", zap.Error(a.s.Serve(l)))
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	a.s.GracefulStop()
}
