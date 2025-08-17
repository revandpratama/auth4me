package usecase

import (
	"github.com/revandpratama/auth4me/internal/auth/entity"
	"github.com/revandpratama/auth4me/internal/auth/repository"
)

type RBACUsecase interface {
	GetAllRolePermissions() ([]entity.RolePermission, error)
	GetRolePermissionsByRoleID(roleID uint) ([]entity.Permission, error)
	GetRolePermissionsByPermissionID(permissionID uint) ([]entity.Role, error)
	GetRolePermissionsByRoleName(roleName string) ([]entity.Permission, error)
	CreateRole(role *entity.Role) error
	UpdateRole(role *entity.Role) error
	DeleteRole(id uint) error
	GetAllRoles() ([]entity.Role, error)
	GetRoleByID(id uint) (*entity.Role, error)
	CreatePermission(permission *entity.Permission) error
	UpdatePermission(permission *entity.Permission) error
	DeletePermission(id uint) error
	GetAllPermissions() ([]entity.Permission, error)
	GetPermissionByID(id uint) (*entity.Permission, error)
	CreateRolePermission(rolePermission *entity.RolePermission) error
	UpdateRolePermission(rolePermission *entity.RolePermission) error
	DeleteRolePermission(id uint) error
}

type rbacUsecase struct {
	repository repository.RBACRepository
}

func NewRBACUsecase(repository repository.RBACRepository) RBACUsecase {
	return &rbacUsecase{repository: repository}
}

func (r *rbacUsecase) GetAllRolePermissions() ([]entity.RolePermission, error) {
	return r.repository.GetAllRolePermissions()
}

func (r *rbacUsecase) GetRolePermissionsByRoleID(roleID uint) ([]entity.Permission, error) {
	return r.repository.GetRolePermissionsByRoleID(roleID)
}

func (r *rbacUsecase) GetRolePermissionsByPermissionID(permissionID uint) ([]entity.Role, error) {
	return r.repository.GetRolePermissionsByPermissionID(permissionID)
}

func (r *rbacUsecase) GetRolePermissionsByRoleName(roleName string) ([]entity.Permission, error) {
	return r.repository.GetRolePermissionsByRoleName(roleName)
}

func (r *rbacUsecase) CreateRole(role *entity.Role) error {
	return r.repository.CreateRole(role)
}

func (r *rbacUsecase) UpdateRole(role *entity.Role) error {
	return r.repository.UpdateRole(role)
}

func (r *rbacUsecase) DeleteRole(id uint) error {
	return r.repository.DeleteRole(id)
}

func (r *rbacUsecase) GetAllRoles() ([]entity.Role, error) {
	return r.repository.GetAllRoles()
}

func (r *rbacUsecase) GetRoleByID(id uint) (*entity.Role, error) {
	return r.repository.GetRoleByID(id)
}

func (r *rbacUsecase) CreatePermission(permission *entity.Permission) error {
	return r.repository.CreatePermission(permission)
}

func (r *rbacUsecase) UpdatePermission(permission *entity.Permission) error {
	return r.repository.UpdatePermission(permission)
}

func (r *rbacUsecase) DeletePermission(id uint) error {
	return r.repository.DeletePermission(id)
}

func (r *rbacUsecase) GetAllPermissions() ([]entity.Permission, error) {
	return r.repository.GetAllPermissions()
}

func (r *rbacUsecase) GetPermissionByID(id uint) (*entity.Permission, error) {
	return r.repository.GetPermissionByID(id)
}

func (r *rbacUsecase) CreateRolePermission(rolePermission *entity.RolePermission) error {
	return r.repository.CreateRolePermission(rolePermission)
}

func (r *rbacUsecase) UpdateRolePermission(rolePermission *entity.RolePermission) error {
	return r.repository.UpdateRolePermission(rolePermission)
}

func (r *rbacUsecase) DeleteRolePermission(id uint) error {
	return r.repository.DeleteRolePermission(id)
}
