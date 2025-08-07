package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/internal/auth/handler"
	"github.com/revandpratama/auth4me/internal/auth/repository"
	"github.com/revandpratama/auth4me/internal/auth/usecase"
	"github.com/revandpratama/auth4me/internal/middleware"
	"gorm.io/gorm"
)

func InitAuthHandler(db *gorm.DB) handler.AuthHandler {
	repo := repository.NewAuthRepository(db)
	usecase := usecase.NewAuthUsecase(repo)
	return handler.NewAuthHandler(usecase)
}
func InitAuthRoutes(api fiber.Router, handler handler.AuthHandler) {

	publicAuth := api.Group("/auth")

	publicAuth.Post("/login", handler.LoginHandler)
	publicAuth.Post("/register", handler.RegisterHandler)
	publicAuth.Post("/logout", handler.LogoutHandler)
	publicAuth.Post("/refresh-token", handler.RefreshTokenHandler)

	auth := api.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	auth.Get("/user", handler.GetUserHandler)

}
