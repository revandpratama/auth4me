package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/internal/auth/dto"
	"github.com/revandpratama/auth4me/internal/auth/usecase"
)

type AuthHandler interface {
	LoginHandler(c *fiber.Ctx) error
	GetUserHandler(c *fiber.Ctx) error
	RegisterHandler(c *fiber.Ctx) error
	LogoutHandler(c *fiber.Ctx) error
	RefreshTokenHandler(c *fiber.Ctx) error
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
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request body",
		})
	}

	refreshToken, accessToken, err := h.authUsecase.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(&Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "login success",
		Data: dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
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

func (h *authHandler) RefreshTokenHandler(c *fiber.Ctx) error {

	refreshToken := c.Get("X-Refresh-Token")
	if refreshToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(&Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized, refresh token empty",
		})
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(&Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized, access token empty",
		})
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return c.Status(http.StatusUnauthorized).JSON(&Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized, malformed access token",
		})
	}

	accessToken := strings.TrimPrefix(authHeader, bearerPrefix)

	newRefreshToken, newAccessToken, err := h.authUsecase.RefreshToken(refreshToken, accessToken)
	if err != nil {
		// LOG THE REAL ERROR
		log.Printf("CRITICAL: Token refresh failed. Raw Error: %v", err)

		// Return a clean message to the client
		return c.Status(http.StatusUnauthorized).JSON(&Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized, invalid token",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "refresh token success",
		Data: dto.TokenResponse{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		},
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
