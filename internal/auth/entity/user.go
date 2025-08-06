package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Password  string `gorm:"" json:"-"` // hashed password; can be empty for OAuth users
	FullName  string `gorm:"size:255" json:"full_name"`
	AvatarPath string `gorm:"size:500" json:"avatar_path"`

	Providers []OAuthProvider `gorm:"foreignKey:UserID" json:"providers"`

	EmailVerified      bool      `gorm:"default:false" json:"email_verified"`
	VerificationToken  string    `gorm:"size:255" json:"-"`
	VerificationSentAt time.Time `json:"-"`

	MFAEnabled bool   `gorm:"default:false" json:"mfa_enabled"`
	MFASecret  string `gorm:"size:255" json:"-"` 

	Role string `gorm:"default:'user'" json:"role"` 

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "auth4me.users"
}