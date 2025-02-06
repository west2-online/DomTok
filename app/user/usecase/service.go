package usecase

import (
	"context"

	"github.com/west2-online/DomTok/app/user/entities"
)

// PersistencePort 表示持久化存储接口 (或者也可以叫做 DBPort)
type PersistencePort interface {
	IsUserExist(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, entity *entities.User) error
}

// CachePort 表示缓存接口
type CachePort interface {
}

type UseCase struct {
	DB PersistencePort
}

func NewUserCase(db PersistencePort) *UseCase {
	return &UseCase{DB: db}
}
