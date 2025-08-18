package app

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/config"
	"github.com/revandpratama/auth4me/internal/auth"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func WithRESTServer() Option {
	return func(app *App) error {

		fiberApp := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		fiberApp.Use(func(c *fiber.Ctx) error {
			c.Set("Content-Type", "application/json")
			return c.Next()
		})

		fiberApp.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		})

		api := fiberApp.Group("/api")

		api.Get("/test-700ms", func(c *fiber.Ctx) error {
			time.Sleep(500 * time.Millisecond)
			return c.SendString("Hello. 700ms delay!")
		})

		authHandler := auth.InitAuthHandler(app.DB)
		auth.InitAuthRoutes(api, authHandler)

		rbacHandler := auth.InitRBACHandler(app.DB)
		auth.InitRBACRoutes(api, rbacHandler)

		oauthConfig := &oauth2.Config{
			ClientID:     config.ENV.GOOGLE_CLIENT_ID,
			ClientSecret: config.ENV.GOOGLE_CLIENT_SECRET,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
			RedirectURL: config.ENV.GOOGLE_REDIRECT_URL,
			// RedirectURL: "http://localhost:3000/auth/google/callback",
		}
		oauthHandler := auth.InitOauthHandler(app.DB, oauthConfig)
		auth.InitOauthRoutes(api, oauthHandler)

		app.fiberApp = fiberApp

		go func() {
			if err := fiberApp.Listen(fmt.Sprintf(":%s", config.ENV.REST_PORT)); err != nil {
				log.Error().Err(err).Msg("failed to start server")
				os.Exit(1)
			}
		}()

		log.Info().Msgf("REST server started on PORT %s", config.ENV.REST_PORT)

		return nil
	}
}
