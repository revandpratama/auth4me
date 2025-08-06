package app

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/config"
	"github.com/rs/zerolog/log"
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
