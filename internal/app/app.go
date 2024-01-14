package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorobot-nz/test-task/pkg/middleware"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"io/fs"
	"net"
	"os"
	"os/signal"
	"strconv"
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

var (
	logLevel int

	adminUsername string
	adminEmail    string
	adminPassword string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	logLevel, err = strconv.Atoi(os.Getenv("LOG_LEVEL"))

	if err != nil {
		panic(err)
	}

	adminUsername = os.Getenv("ADMIN_USERNAME")
	adminEmail = os.Getenv("ADMIN_EMAIL")
	adminPassword = os.Getenv("ADMIN_PASSWORD")

}

type App struct {
	logger *zap.Logger

	s     *grpc.Server
	store *storage.Storage[userv1.User]
}

func NewApp() *App {
	logger := applogger.NewLogger(zapcore.Level(logLevel))
	store := storage.NewStorage[userv1.User]()

	repository := usersrepository.NewStorageRepository(logger.Named("UsersRepository"), store)
	service := usersservice.NewService(logger.Named("UsersService"), repository)
	handler := usershandler.NewHandler(logger.Named("UsersHandler"), service)

	server := grpc.NewServer(grpc.UnaryInterceptor(
		grpcmiddleware.ChainUnaryServer(
			grpczap.UnaryServerInterceptor(logger),
			middleware.AuthMiddleware(),
		),
	))

	userv1.RegisterUserServiceServer(server, handler)

	return &App{
		s:      server,
		logger: logger,
		store:  store,
	}
}

func (a *App) Run() {
	a.initStore()

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

	a.saveStore()
}

func (a *App) initStore() {
	file, err := os.Open("users_store.txt")
	if err != nil {
		if adminEmail == "" || adminUsername == "" || adminPassword == "" {
			a.logger.Fatal("No default admin params")
		}

		id := uuid.New().String()

		password, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 10)
		if err != nil {
			a.logger.Fatal("Failed to generate password", zap.Error(err))
		}

		a.store.Set(id, userv1.User{
			Id:       id,
			Email:    adminEmail,
			Username: adminUsername,
			Password: string(password),
			Admin:    true,
		})
	} else {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				a.logger.Error("Failed to close file", zap.Error(err))
			}
		}(file)

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			var user userv1.User
			userBytes := scanner.Bytes()
			err := json.Unmarshal(userBytes, &user)
			if err != nil {
				a.logger.Fatal("Failed to load user", zap.Error(err))
			}
			a.store.Set(user.GetId(), user)
		}

		if a.store.Size() == 0 {
			if adminEmail == "" || adminUsername == "" || adminPassword == "" {
				a.logger.Fatal("No default admin params")
			}

			id := uuid.New().String()

			password, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 10)
			if err != nil {
				a.logger.Fatal("Failed to generate password", zap.Error(err))
			}

			a.store.Set(id, userv1.User{
				Id:       id,
				Email:    adminEmail,
				Username: adminUsername,
				Password: string(password),
				Admin:    true,
			})
		}
	}
}

func (a *App) saveStore() {
	file, err := os.Create("users_store.txt")
	if err != nil {
		a.logger.Error("Failed to create file", zap.Error(err))
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			a.logger.Error("Failed to close file", zap.Error(err))
		}
	}(file)

	list := a.store.List()

	for _, item := range list {
		marshal, err := json.Marshal(&item)
		if err != nil {
			a.logger.Error("Failed to marshal item", zap.Error(err))
			return
		}
		err = os.WriteFile(file.Name(), marshal, fs.ModeAppend)
		if err != nil {
			a.logger.Error("Failed to write", zap.Error(err))
			return
		}
	}
}
