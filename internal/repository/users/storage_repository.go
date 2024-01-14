package users

import (
	"context"
	userv1 "github.com/gorobot-nz/test-task/gen/proto/user/v1"
	"github.com/gorobot-nz/test-task/pkg/storage"
	"go.uber.org/zap"
)

type StorageRepository struct {
	storage *storage.Storage[*userv1.User]

	logger *zap.Logger
}

func NewStorageRepository(logger *zap.Logger) *StorageRepository {
	return &StorageRepository{
		logger:  logger,
		storage: storage.NewStorage[*userv1.User](),
	}
}

func (s StorageRepository) Create(ctx context.Context, user *userv1.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s StorageRepository) List(ctx context.Context, page, limit int32) ([]*userv1.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s StorageRepository) GetById(ctx context.Context, id string) (*userv1.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s StorageRepository) Update(ctx context.Context, user *userv1.User) (*userv1.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s StorageRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
