package app

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func WithRESTServer(addr string) Option {
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
			if err := fiberApp.Listen(addr); err != nil {
				log.Error().Err(err).Msg("failed to start server")
				os.Exit(1)
			}
		}()

		log.Info().Msgf("REST server started on PORT %s", addr)

		return nil
	}
}
