package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/user/entities"
	"github.com/west2-online/DomTok/pkg/errno"
)

// DBAdapter impl PersistencePort defined in use case package
type DBAdapter struct {
	client *gorm.DB
}

func NewDBAdapter(client *gorm.DB) *DBAdapter {
	return &DBAdapter{client: client}
}

func (d *DBAdapter) CreateUser(ctx context.Context, entity *entities.User) error {
	// 将 entity 转换成 mysql 这边的 model
	model := User{
		UserName: entity.UserName,
		Password: entity.Password,
		Email:    entity.Email,
	}
	if err := d.client.WithContext(ctx).Create(model).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create user: %v", err)
	}
	return nil
}

func (d *DBAdapter) IsUserExist(ctx context.Context, username string) (bool, error) {
	var user entities.User
	err := d.client.WithContext(ctx).Where("user_name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query user: %v", err)
	}
	return true, nil
}
