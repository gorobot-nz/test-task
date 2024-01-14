package users

import (
	"context"
	userv1 "github.com/gorobot-nz/test-task/gen/proto/user/v1"
	"go.uber.org/zap"
)

type Repository interface {
	NewUser(ctx context.Context, user *userv1.User) (string, error)
	GetUsers(ctx context.Context, page, limit int32) ([]*userv1.User, error)
	GetUserById(ctx context.Context, id string) (*userv1.User, error)
	GetUserByUsername(ctx context.Context, username string) (*userv1.User, error)
	UpdateUser(ctx context.Context, user *userv1.User) (*userv1.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type Service struct {
	repository Repository

	logger *zap.Logger
}

func NewService(logger *zap.Logger, repository Repository) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}

func (s *Service) NewUser(ctx context.Context, user *userv1.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetUsers(ctx context.Context, page, limit int32) ([]*userv1.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetUserById(ctx context.Context, id string) (*userv1.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*userv1.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) UpdateUser(ctx context.Context, user *userv1.User) (*userv1.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
