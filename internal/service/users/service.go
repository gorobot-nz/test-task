package users

import (
	"context"
	"errors"
	userv1 "github.com/gorobot-nz/test-task/gen/proto/user/v1"
	"github.com/gorobot-nz/test-task/pkg/validation"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(ctx context.Context, user *userv1.User) (string, error)
	List(ctx context.Context, page, limit int32) ([]*userv1.User, error)
	GetById(ctx context.Context, id string) (*userv1.User, error)
	Update(ctx context.Context, user *userv1.User) (*userv1.User, error)
	Delete(ctx context.Context, id string) error
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
	log := s.logger.Named("NewUser")

	if !validation.IsValidEmail(user.GetEmail()) {
		return "", errors.New("failure mail validation")
	}

	if !validation.IsValidPassword(user.GetPassword()) {
		return "", errors.New("failure password validation")
	}

	if !validation.IsValidUsername(user.GetUsername()) {
		return "", errors.New("failure username validation")
	}

	list, err := s.repository.List(ctx, -1, -1)
	if err != nil {
		log.Error("Failed to list users", zap.Error(err))
		return "", err
	}

	for _, val := range list {
		if val.GetUsername() == user.GetUsername() || val.GetEmail() == user.GetEmail() {
			return "", errors.New("not unique email or username")
		}
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.GetPassword()), 10)
	if err != nil {
		log.Error("Failed to generate password", zap.Error(err))
		return "", err
	}

	user.Password = string(password)

	id, err := s.repository.Create(ctx, user)
	if err != nil {
		log.Error("Failed to create user", zap.Error(err))
		return "", err
	}

	return id, nil
}

func (s *Service) GetUsers(ctx context.Context, page, limit int32) ([]*userv1.User, error) {
	log := s.logger.Named("GetUsers")

	list, err := s.repository.List(ctx, page, limit)
	if err != nil {
		log.Error("Failed to list users", zap.Error(err))
		return nil, err
	}

	return list, nil
}

func (s *Service) GetUserById(ctx context.Context, id string) (*userv1.User, error) {
	log := s.logger.Named("GetUserById")

	user, err := s.repository.GetById(ctx, id)
	if err != nil {
		log.Error("Failed to get user", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*userv1.User, error) {
	log := s.logger.Named("GetUserByUsername")

	list, err := s.repository.List(ctx, -1, -1)
	if err != nil {
		log.Error("Failed to list users", zap.Error(err))
		return nil, err
	}

	for _, val := range list {
		if val.GetUsername() == username {
			return val, nil
		}
	}

	return nil, errors.New("no such user")
}

func (s *Service) UpdateUser(ctx context.Context, user *userv1.User) (*userv1.User, error) {
	log := s.logger.Named("UpdateUser")

	if !validation.IsValidEmail(user.GetEmail()) {
		return nil, errors.New("failure mail validation")
	}

	if !validation.IsValidPassword(user.GetPassword()) {
		return nil, errors.New("failure password validation")
	}

	if !validation.IsValidUsername(user.GetUsername()) {
		return nil, errors.New("failure username validation")
	}

	updatedUser, err := s.repository.Update(ctx, user)
	if err != nil {
		log.Error("Failed to update user", zap.Error(err))
		return nil, err
	}

	return updatedUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	log := s.logger.Named("DeleteUser")

	err := s.repository.Delete(ctx, id)
	if err != nil {
		log.Error("Failed to delete user", zap.Error(err))
		return err
	}

	return nil
}
