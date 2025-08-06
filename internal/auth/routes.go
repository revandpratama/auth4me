package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/internal/auth/handler"
	"github.com/revandpratama/auth4me/internal/auth/repository"
	"github.com/revandpratama/auth4me/internal/auth/usecase"
	"gorm.io/gorm"
)

func InitAuthHandler(db *gorm.DB) handler.AuthHandler {
	repo := repository.NewAuthRepository(db)
	usecase := usecase.NewAuthUsecase(repo)
	return handler.NewAuthHandler(usecase)
}
func InitAuthRoutes(api fiber.Router, handler handler.AuthHandler) {

	auth := api.Group("/auth")

	auth.Post("/login", handler.LoginHandler)
	auth.Post("/register", handler.RegisterHandler)
	auth.Post("/logout", handler.LogoutHandler)
	
	// auth.Use()
	// auth.Get("/user", handler.GetUserHandler)

}
