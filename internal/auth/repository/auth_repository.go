package repository

import (
	"errors"

	"github.com/revandpratama/auth4me/internal/auth/entity"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	IsEmailExists(email string) (bool, error)
	CreateUser(user *entity.User) error
	GetUserPermissionsByRoleID(id uint) ([]entity.Permission, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := r.db.Model(&entity.User{}).Preload("Role.Permissions").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err // return actual DB error
	}
	return &user, nil
}
func (r *authRepository) GetUserByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Model(&entity.User{}).Preload("Role.Permissions").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) IsEmailExists(email string) (bool, error) {
	var user entity.User
	err := r.db.Select("id").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err // return actual DB error
	}
	return true, nil
}

func (r *authRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) GetUserPermissionsByRoleID(roleID uint) ([]entity.Permission, error) {
	var permissions []entity.Permission

	err := r.db.Joins("JOIN auth4me.role_permissions ON auth4me.role_permissions.permission_id = auth4me.permissions.id").Where("auth4me.role_permissions.role_id = ?", roleID).Find(&permissions).Error
	if err != nil {
		return nil, err // return actual DB error
	}
	return permissions, nil
}
