package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/internal/auth/dto"
	"github.com/revandpratama/auth4me/internal/auth/usecase"
)

type OAuthHandler interface {
	GoogleLogin(c *fiber.Ctx) error
	GoogleOAuthCallback(c *fiber.Ctx) error
}

type oauthHandler struct {
	usecase usecase.OAuthUsecase
}

func NewOAuthHandler(usecase usecase.OAuthUsecase) OAuthHandler {
	return &oauthHandler{
		usecase: usecase,
	}
}

func (h *oauthHandler) GoogleLogin(c *fiber.Ctx) error {

	url, state := h.usecase.GetOAuthURL()
	if url == "" {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error: failed to get oauth url",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "google oauth login success",
		Data: map[string]string{
			"url":   url,
			"state": state,
		},
	})
}

func (h *oauthHandler) GoogleOAuthCallback(c *fiber.Ctx) error {

	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request: code is required",
		})
	}

	stateFromQuery := c.Query("state")
	if stateFromQuery == "" {
		return c.Status(http.StatusBadRequest).JSON(&Response{
			Code:    http.StatusBadRequest,
			Message: "bad request: state is required",
		})
	}

	stateFromHeader := c.Get("X-OAuth-State") 

	if stateFromHeader == "" || stateFromHeader != stateFromQuery {
        return c.Status(http.StatusBadRequest).JSON(&Response{
            Code: http.StatusBadRequest,
            Message: "bad request: invalid state token",
        })
    }

	refreshToken, accessToken, err := h.usecase.GoogleOAuthCallback(code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(&Response{
		Code:    http.StatusOK,
		Message: "google oauth login success",
		Data: dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})

}
