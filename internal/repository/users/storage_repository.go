package users

import (
	"context"
	"errors"

	userv1 "github.com/gorobot-nz/test-task/gen/proto/user/v1"

	"github.com/gorobot-nz/test-task/pkg/storage"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type StorageRepository struct {
	storage *storage.Storage[userv1.User]

	logger *zap.Logger
}

func NewStorageRepository(logger *zap.Logger) *StorageRepository {
	return &StorageRepository{
		logger:  logger,
		storage: storage.NewStorage[userv1.User](),
	}
}

func (s *StorageRepository) Create(ctx context.Context, user *userv1.User) (string, error) {
	_ = s.logger.Named("Create")

	user.Id = uuid.New().String()

	s.storage.Set(user.GetId(), userv1.User{
		Id:       user.GetId(),
		Email:    user.GetEmail(),
		Username: user.GetEmail(),
		Password: user.GetPassword(),
		Admin:    user.GetAdmin(),
	})

	return user.GetId(), nil
}

func (s *StorageRepository) List(ctx context.Context, page, limit int32) ([]*userv1.User, error) {
	_ = s.logger.Named("List")

	list := s.storage.List()

	resultList := make([]*userv1.User, len(list))

	for index := range resultList {
		resultList[index] = &list[index]
	}

	if page < 0 && limit < 0 {
		return resultList, nil
	}

	offset := (page - 1) * limit
	border := offset + limit

	if offset >= int32(len(resultList)) {
		return nil, errors.New("last page")
	}

	if border >= int32(len(resultList)) {
		border = int32(len(resultList))
	}

	resultList = resultList[offset : offset+limit]

	return resultList, nil
}

func (s *StorageRepository) GetById(ctx context.Context, id string) (*userv1.User, error) {
	_ = s.logger.Named("GetById")

	get, ok := s.storage.Get(id)

	if !ok {
		return nil, errors.New("no such user")
	}

	return &get, nil
}

func (s *StorageRepository) Update(ctx context.Context, user *userv1.User) (*userv1.User, error) {
	_ = s.logger.Named("Update")

	_, ok := s.storage.Get(user.GetId())

	if !ok {
		return nil, errors.New("no such user")
	}

	s.storage.Set(user.GetId(), userv1.User{
		Id:       user.GetId(),
		Email:    user.GetEmail(),
		Username: user.GetEmail(),
		Password: user.GetPassword(),
		Admin:    user.GetAdmin(),
	})

	return user, nil
}

func (s *StorageRepository) Delete(ctx context.Context, id string) error {
	_ = s.logger.Named("Delete")

	_, ok := s.storage.Get(id)

	if !ok {
		return errors.New("no such user")
	}
	s.storage.Delete(id)

	return nil
}
