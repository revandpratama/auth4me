package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/internal/auth/handler"
	"github.com/revandpratama/auth4me/internal/auth/repository"
	"github.com/revandpratama/auth4me/internal/auth/usecase"
	"github.com/revandpratama/auth4me/internal/middleware"
	"golang.org/x/oauth2"
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

func InitRBACHandler(db *gorm.DB) handler.RBACHandler {
	repo := repository.NewRBACRepository(db)
	usecase := usecase.NewRBACUsecase(repo)
	return handler.NewRBACHandler(usecase)
}

func InitRBACRoutes(api fiber.Router, handler handler.RBACHandler) {

	rbac := api.Group("/rbac")

	rbac.Use(middleware.AuthMiddleware())

	rbac.Get("/roles", handler.GetAllRoles)
	rbac.Get("/roles/:id", handler.GetRoleByID)
	rbac.Post("/roles", handler.CreateRole)
	rbac.Put("/roles/:id", handler.UpdateRole)
	rbac.Delete("/roles/:id", handler.DeleteRole)

	rbac.Get("/permissions", handler.GetAllPermissions)
	rbac.Get("/permissions/:id", handler.GetPermissionByID)
	rbac.Post("/permissions", handler.CreatePermission)
	rbac.Put("/permissions/:id", handler.UpdatePermission)
	rbac.Delete("/permissions/:id", handler.DeletePermission)

	rbac.Get("/role-permissions", handler.GetAllRolePermissions)
	rbac.Get("/role-permissions/roles/:roleID", handler.GetRolePermissionsByRoleID)
	rbac.Get("/role-permissions/permissions/:permissionID", handler.GetRolePermissionsByPermissionID)
	rbac.Post("/role-permissions", handler.CreateRolePermission)
	rbac.Put("/role-permissions/:id", handler.UpdateRolePermission)
	rbac.Delete("/role-permissions/:id", handler.DeleteRolePermission)

}

func InitOauthHandler(db *gorm.DB, oauthCfg *oauth2.Config) handler.OAuthHandler {
	oauthRepo := repository.NewOAuthRepository(db)
	authRepo := repository.NewAuthRepository(db)
	usecase := usecase.NewOAuthUsecase(oauthCfg, authRepo, oauthRepo)
	return handler.NewOAuthHandler(usecase)
}

func InitOauthRoutes(api fiber.Router, handler handler.OAuthHandler) {

	oauth := api.Group("/oauth")

	oauth.Get("/google", handler.GoogleLogin)
	oauth.Get("/google/callback", handler.GoogleOAuthCallback)

}
