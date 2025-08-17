package repository

import (
	"github.com/revandpratama/auth4me/internal/auth/entity"
	"gorm.io/gorm"
)

type RBACRepository interface {
	GetAllRolePermissions() ([]entity.RolePermission, error)
	GetRolePermissionsByRoleID(roleID uint) ([]entity.Permission, error)
	GetRolePermissionsByPermissionID(permissionID uint) ([]entity.Role, error)
	GetRolePermissionsByRoleName(roleName string) ([]entity.Permission, error)
	CreateRolePermission(rolePermission *entity.RolePermission) error
	UpdateRolePermission(rolePermission *entity.RolePermission) error
	DeleteRolePermission(id uint) error
	GetAllRoles() ([]entity.Role, error)
	GetRoleByID(id uint) (*entity.Role, error)
	CreateRole(role *entity.Role) error
	UpdateRole(role *entity.Role) error
	DeleteRole(id uint) error
	GetAllPermissions() ([]entity.Permission, error)
	GetPermissionByID(id uint) (*entity.Permission, error)
	CreatePermission(permission *entity.Permission) error
	UpdatePermission(permission *entity.Permission) error
	DeletePermission(id uint) error
}

type rbacRepository struct {
	db *gorm.DB
}

func NewRBACRepository(db *gorm.DB) RBACRepository {
	return &rbacRepository{
		db: db,
	}
}

func (r *rbacRepository) GetAllRolePermissions() ([]entity.RolePermission, error) {
	var rolePermissions []entity.RolePermission
	err := r.db.Find(&rolePermissions).Error
	if err != nil {
		return nil, err
	}
	return rolePermissions, nil
}

func (r *rbacRepository) GetRolePermissionsByRoleID(roleID uint) ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.db.Joins("JOIN auth4me.role_permissions ON auth4me.role_permissions.permission_id = auth4me.permissions.id").Where("auth4me.role_permissions.role_id = ?", roleID).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *rbacRepository) GetRolePermissionsByPermissionID(permissionID uint) ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Joins("JOIN auth4me.role_permissions ON auth4me.role_permissions.role_id = auth4me.roles.id").Where("auth4me.role_permissions.permission_id = ?", permissionID).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *rbacRepository) GetRolePermissionsByRoleName(roleName string) ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.db.Joins("JOIN auth4me.role_permissions ON auth4me.role_permissions.permission_id = auth4me.permissions.id").Where("auth4me.roles.name = ?", roleName).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *rbacRepository) CreateRolePermission(rolePermission *entity.RolePermission) error {
	return r.db.Create(rolePermission).Error
}

func (r *rbacRepository) UpdateRolePermission(rolePermission *entity.RolePermission) error {
	err := r.db.Updates(rolePermission).Error
	return err
}

func (r *rbacRepository) DeleteRolePermission(id uint) error {
	err := r.db.Delete(&entity.RolePermission{}, "id = ?", id).Error
	return err
}

func (r *rbacRepository) GetAllRoles() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Model(&entity.Role{}).Preload("Permissions").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *rbacRepository) GetRoleByID(id uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *rbacRepository) CreateRole(role *entity.Role) error {
	return r.db.Create(role).Error
}

func (r *rbacRepository) UpdateRole(role *entity.Role) error {
	err := r.db.Updates(role).Error
	return err
}

func (r *rbacRepository) DeleteRole(id uint) error {
	err := r.db.Delete(&entity.Role{}, "id = ?", id).Error
	return err
}

func (r *rbacRepository) GetAllPermissions() ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.db.Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *rbacRepository) GetPermissionByID(id uint) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.First(&permission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *rbacRepository) CreatePermission(permission *entity.Permission) error {
	return r.db.Create(permission).Error
}

func (r *rbacRepository) UpdatePermission(permission *entity.Permission) error {
	err := r.db.Updates(permission).Error
	return err
}

func (r *rbacRepository) DeletePermission(id uint) error {
	err := r.db.Delete(&entity.Permission{}, "id = ?", id).Error
	return err
}
