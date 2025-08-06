package repository

import (
	"errors"

	"github.com/revandpratama/auth4me/internal/auth/entity"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(email string, password string) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	IsEmailExists(email string) (bool, error)
	CreateUser(user *entity.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Login(email string, password string) error {
	return nil
}

func (r *authRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err // return actual DB error
	}
	return &user, nil
}
func (r *authRepository) GetUserByID(id string) (*entity.User, error) {
	return nil, nil
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
