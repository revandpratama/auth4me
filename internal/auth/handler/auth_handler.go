package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/internal/auth/dto"
	"github.com/revandpratama/auth4me/internal/auth/usecase"
)

type AuthHandler interface {
	LoginHandler(c *fiber.Ctx) error
	GetUserHandler(c *fiber.Ctx) error
	RegisterHandler(c *fiber.Ctx) error
	LogoutHandler(c *fiber.Ctx) error
}

type authHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *authHandler {
	return &authHandler{
		authUsecase: authUsecase,
	}
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (h *authHandler) LoginHandler(c *fiber.Ctx) error {

	var loginRequest dto.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return err
	}

	token, err := h.authUsecase.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "login success",
		Data: dto.LoginResponse{
			Token: token,
		},
	})
}

func (h *authHandler) RegisterHandler(c *fiber.Ctx) error {

	var registerRequest dto.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return err
	}

	if err := h.authUsecase.Register(&registerRequest); err != nil {
		return err
	}

	// TODO: Send verification email

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "register success",
	})
}

func (h *authHandler) LogoutHandler(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "logout success",
	})
}

func (h *authHandler) GetUserHandler(c *fiber.Ctx) error {

	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(&Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized, user id is nil",
		})
	}

	user, err := h.authUsecase.GetUserByID(userID.(string))
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "get user success",
		Data:    user,
	})
}
