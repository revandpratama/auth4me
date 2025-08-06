package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	fiberApp *fiber.App
	DB       *gorm.DB
}

type Option func(*App) error

func NewApp(opts ...Option) (*App, error) {
	app := &App{}
	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (a *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.fiberApp.ShutdownWithContext(ctx); err != nil {
		return err
	}
	
	return nil
}
