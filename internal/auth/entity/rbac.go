package entity

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Permissions []Permission `gorm:"many2many:auth4me.role_permissions;" json:"permissions"`
}

func (Role) TableName() string {
	return "auth4me.roles"
}

type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Permission) TableName() string {
	return "auth4me.permissions"
}

type RolePermission struct {
	RoleID       uint `gorm:"type:uuid;index" json:"role_id"`
	PermissionID uint `gorm:"type:uuid;index" json:"permission_id"`
}

func (RolePermission) TableName() string {
	return "auth4me.role_permissions"
}
