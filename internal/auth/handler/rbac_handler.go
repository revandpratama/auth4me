package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/internal/auth/entity"
	"github.com/revandpratama/auth4me/internal/auth/usecase"
)

type RBACHandler interface {
	GetAllRolePermissions(c *fiber.Ctx) error
	GetRolePermissionsByRoleID(c *fiber.Ctx) error
	GetRolePermissionsByPermissionID(c *fiber.Ctx) error
	GetRolePermissionsByRoleName(c *fiber.Ctx) error
	CreateRolePermission(c *fiber.Ctx) error
	DeleteRolePermission(c *fiber.Ctx) error
	UpdateRolePermission(c *fiber.Ctx) error

	GetAllRoles(c *fiber.Ctx) error
	GetRoleByID(c *fiber.Ctx) error
	CreateRole(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error

	GetAllPermissions(c *fiber.Ctx) error
	GetPermissionByID(c *fiber.Ctx) error
	CreatePermission(c *fiber.Ctx) error
	UpdatePermission(c *fiber.Ctx) error
	DeletePermission(c *fiber.Ctx) error
}

type rbacHandler struct {
	rbacUsecase usecase.RBACUsecase
}

func NewRBACHandler(rbacUsecase usecase.RBACUsecase) RBACHandler {
	return &rbacHandler{rbacUsecase}
}

func (h *rbacHandler) GetAllRolePermissions(c *fiber.Ctx) error {

	rolePermissions, err := h.rbacUsecase.GetAllRolePermissions()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "get all role permissions success",
		Data:    rolePermissions,
	})
}

func (h *rbacHandler) GetRolePermissionsByRoleID(c *fiber.Ctx) error {

	roleID, err := strconv.Atoi(c.Params("roleID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}
	rolePermissions, err := h.rbacUsecase.GetRolePermissionsByRoleID(uint(roleID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "get role permissions by role id success",
		Data:    rolePermissions,
	})
}

func (h *rbacHandler) GetRolePermissionsByPermissionID(c *fiber.Ctx) error {

	permissionID, err := strconv.Atoi(c.Params("permissionID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}
	rolePermissions, err := h.rbacUsecase.GetRolePermissionsByPermissionID(uint(permissionID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "get role permissions by permission id success",
		Data:    rolePermissions,
	})
}

func (h *rbacHandler) GetRolePermissionsByRoleName(c *fiber.Ctx) error {

	roleName := c.Params("roleName")
	rolePermissions, err := h.rbacUsecase.GetRolePermissionsByRoleName(roleName)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "get role permissions by role name success",
		Data:    rolePermissions,
	})
}

func (h *rbacHandler) CreateRolePermission(c *fiber.Ctx) error {

	var rolePermission entity.RolePermission
	if err := c.BodyParser(&rolePermission); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}
	err := h.rbacUsecase.CreateRolePermission(&rolePermission)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "create role permission success",
	})
}

func (h *rbacHandler) UpdateRolePermission(c *fiber.Ctx) error {

	var rolePermission entity.RolePermission
	if err := c.BodyParser(&rolePermission); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}
	err := h.rbacUsecase.UpdateRolePermission(&rolePermission)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "update role permission success",
	})
}

func (h *rbacHandler) DeleteRolePermission(c *fiber.Ctx) error {

	rolePermissionID, err := strconv.Atoi(c.Params("rolePermissionID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}
	err = h.rbacUsecase.DeleteRolePermission(uint(rolePermissionID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "delete role permission success",
	})
}

func (h *rbacHandler) GetAllRoles(c *fiber.Ctx) error {

	roles, err := h.rbacUsecase.GetAllRoles()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "get all roles success",
		Data:    roles,
	})
}

func (h *rbacHandler) GetRoleByID(c *fiber.Ctx) error {

	roleID, err := strconv.Atoi(c.Params("roleID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}
	role, err := h.rbacUsecase.GetRoleByID(uint(roleID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "get role by id success",
		Data:    role,
	})
}

func (h *rbacHandler) CreateRole(c *fiber.Ctx) error {

	var role entity.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	if err := h.rbacUsecase.CreateRole(&role); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "create role success",
	})
}

func (h *rbacHandler) UpdateRole(c *fiber.Ctx) error {

	var role entity.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	if err := h.rbacUsecase.UpdateRole(&role); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "update role success",
	})
}

func (h *rbacHandler) DeleteRole(c *fiber.Ctx) error {

	roleID, err := strconv.Atoi(c.Params("roleID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	if err := h.rbacUsecase.DeleteRole(uint(roleID)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "delete role success",
	})
}

func (h *rbacHandler) GetAllPermissions(c *fiber.Ctx) error {

	permissions, err := h.rbacUsecase.GetAllPermissions()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "get all permissions success",
		Data:    permissions,
	})
}

func (h *rbacHandler) GetPermissionByID(c *fiber.Ctx) error {

	permissionID, err := strconv.Atoi(c.Params("permissionID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}
	permission, err := h.rbacUsecase.GetPermissionByID(uint(permissionID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "get permission by id success",
		Data:    permission,
	})
}

func (h *rbacHandler) CreatePermission(c *fiber.Ctx) error {

	var permission entity.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	if err := h.rbacUsecase.CreatePermission(&permission); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "create permission success",
	})
}

func (h *rbacHandler) UpdatePermission(c *fiber.Ctx) error {

	var permission entity.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	if err := h.rbacUsecase.UpdatePermission(&permission); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "update permission success",
	})
}

func (h *rbacHandler) DeletePermission(c *fiber.Ctx) error {

	permissionID, err := strconv.Atoi(c.Params("permissionID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	if err := h.rbacUsecase.DeletePermission(uint(permissionID)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "delete permission success",
	})
}
