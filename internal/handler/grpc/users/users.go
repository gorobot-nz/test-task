package users

import (
	"context"
	"github.com/gorobot-nz/test-task/pkg/middleware"
	"golang.org/x/crypto/bcrypt"

	userv1 "github.com/gorobot-nz/test-task/gen/proto/user/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	NewUser(ctx context.Context, user *userv1.User) (string, error)
	GetUsers(ctx context.Context, page, limit int32) ([]*userv1.User, error)
	GetUserById(ctx context.Context, id string) (*userv1.User, error)
	GetUserByUsername(ctx context.Context, username string) (*userv1.User, error)
	UpdateUser(ctx context.Context, user *userv1.User) (*userv1.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type Handler struct {
	userv1.UnimplementedUserServiceServer

	service Service

	logger *zap.Logger
}

func NewHandler(logger *zap.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) NewUser(ctx context.Context, req *userv1.NewUserRequest) (*userv1.NewUserResponse, error) {
	log := h.logger.Named("NewUser")

	log.Debug("Request received", zap.Any("req", req))

	username := ctx.Value(middleware.Username).(string)
	password := ctx.Value(middleware.Password).(string)

	u, err := h.service.GetUserByUsername(ctx, username)
	if err != nil {
		log.Error("Failed to get admin user", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.GetPassword()), []byte(password))
	if err != nil {
		log.Error("Failed to verify password admin user", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	if !u.GetAdmin() {
		log.Error("Failed to verify admin status", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	user := &userv1.User{
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Admin:    req.GetAdmin(),
	}

	id, err := h.service.NewUser(ctx, user)
	if err != nil {
		log.Error("Failed to add new user", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Wrong request params")
	}

	log.Debug("New user's id", zap.String("id", id))

	return &userv1.NewUserResponse{Id: id}, nil
}

func (h *Handler) GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	log := h.logger.Named("GetUsers")

	log.Debug("Request received", zap.Any("req", req))

	var page, limit int32 = -1, -1

	if req.Page != nil {
		page = req.GetPage()
	}

	if req.Limit != nil {
		limit = req.GetLimit()
	}

	users, err := h.service.GetUsers(ctx, page, limit)
	if err != nil {
		log.Error("Failed to get users", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get users")
	}

	return &userv1.GetUsersResponse{Users: users}, nil
}

func (h *Handler) GetUserById(ctx context.Context, req *userv1.GetUserByIdRequest) (*userv1.GetUserByIdResponse, error) {
	log := h.logger.Named("GetUserById")

	log.Debug("Request received", zap.Any("req", req))

	user, err := h.service.GetUserById(ctx, req.GetId())
	if err != nil {
		log.Error("Failed to get user", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Failed to get user")
	}

	return &userv1.GetUserByIdResponse{User: user}, nil
}

func (h *Handler) GetUserByUsername(ctx context.Context, req *userv1.GetUserByUsernameRequest) (*userv1.GetUserByUsernameResponse, error) {
	log := h.logger.Named("GetUserByUsername")

	log.Debug("Request received", zap.Any("req", req))

	user, err := h.service.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		log.Error("Failed to get user", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Failed to get user")
	}

	return &userv1.GetUserByUsernameResponse{User: user}, nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	log := h.logger.Named("UpdateUser")

	log.Debug("Request received", zap.Any("req", req))

	username := ctx.Value(middleware.Username).(string)
	password := ctx.Value(middleware.Password).(string)

	u, err := h.service.GetUserByUsername(ctx, username)
	if err != nil {
		log.Error("Failed to get admin user", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.GetPassword()), []byte(password))
	if err != nil {
		log.Error("Failed to verify password admin user", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	if !u.GetAdmin() {
		log.Error("Failed to verify admin status", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	user := &userv1.User{
		Id:       req.GetId(),
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Admin:    req.GetAdmin(),
	}

	user, err = h.service.UpdateUser(ctx, user)
	if err != nil {
		log.Error("Failed to update user", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Failed to update user")
	}

	return &userv1.UpdateUserResponse{User: user}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	log := h.logger.Named("UpdateUser")

	log.Debug("Request received", zap.Any("req", req))

	username := ctx.Value(middleware.Username).(string)
	password := ctx.Value(middleware.Password).(string)

	u, err := h.service.GetUserByUsername(ctx, username)
	if err != nil {
		log.Error("Failed to get admin user", zap.Error(err))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.GetPassword()), []byte(password))
	if err != nil {
		log.Error("Failed to verify password admin user", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	if !u.GetAdmin() {
		log.Error("Failed to verify admin status", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	if u.GetId() == req.GetId() {
		log.Error("Cannot delete yourself", zap.Error(err))
		return nil, status.Error(codes.Aborted, "You can't delete yourself")
	}

	err = h.service.DeleteUser(ctx, req.GetId())
	if err != nil {
		log.Error("Failed to delete user", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Failed to delete user")
	}

	return &userv1.DeleteUserResponse{}, nil
}
