package users

import (
	"context"
	userv1 "github.com/gorobot-nz/test-task/gen/proto/user/v1"
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
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) NewUser(context.Context, *userv1.NewUserRequest) (*userv1.NewUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewUser not implemented")
}
func (h *Handler) GetUsers(context.Context, *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}
func (h *Handler) GetUserById(context.Context, *userv1.GetUserByIdRequest) (*userv1.GetUserByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserById not implemented")
}
func (h *Handler) GetUserByUsername(context.Context, *userv1.GetUserByUsernameRequest) (*userv1.GetUserByUsernameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByUsername not implemented")
}
func (h *Handler) UpdateUser(context.Context, *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (h *Handler) DeleteUser(context.Context, *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
