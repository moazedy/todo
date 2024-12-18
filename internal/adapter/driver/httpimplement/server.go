package httpimplement

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/moazedy/todo/pkg/infra/config"
)

func Start(cfg config.Config) {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
	}))

	p := prometheus.NewPrometheus(cfg.Server.Name, nil)
	p.Use(e)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(CustomErrorHandlerMiddleWare)

	register(e, cfg)

	go func() {
		if err := e.Start(cfg.Server.Address + ":" + cfg.Server.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("shutting down server: " + err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of X seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err := e.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
