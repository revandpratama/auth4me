package repository

import (
	"github.com/revandpratama/auth4me/internal/auth/entity"
	"gorm.io/gorm"
)

type OAuthRepository interface {
	GetProvider(userID string, providerName string) (*entity.OAuthProvider, error)
	CreateProvider(provider *entity.OAuthProvider) error
	UpdateProvider(provider *entity.OAuthProvider) error
}

type oauthRepository struct {
	db *gorm.DB
}

func NewOAuthRepository(db *gorm.DB) OAuthRepository {
	return &oauthRepository{
		db: db,
	}
}

func (r *oauthRepository) GetProvider(userID string, providerName string) (*entity.OAuthProvider, error) {
	var provider entity.OAuthProvider
	err := r.db.Where("user_id = ? AND provider = ?", userID, providerName).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *oauthRepository) CreateProvider(provider *entity.OAuthProvider) error {
	return r.db.Create(provider).Error
}

func (r *oauthRepository) UpdateProvider(provider *entity.OAuthProvider) error {
	return r.db.Model(&entity.OAuthProvider{}).Where("user_id = ?", provider.UserID).Updates(provider).Error
}
