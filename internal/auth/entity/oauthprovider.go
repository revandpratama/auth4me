package entity

import "time"

type OAuthProvider struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"type:uuid;index" json:"user_id"`
	Provider     string    `gorm:"size:100;index" json:"provider"`    
	ProviderID   string    `gorm:"size:255;index" json:"provider_id"`
	AccessToken  string    `gorm:"size:500" json:"-"`
	RefreshToken string    `gorm:"size:500" json:"-"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (OAuthProvider) TableName() string {
	return "auth4me.oauth_providers"
}