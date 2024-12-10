package main

import (
	"github.com/moazedy/todo/internal/adapter/driver/httpimplement"
	"github.com/moazedy/todo/pkg/infra/config"
)

func main() {
	cfg := config.Init()

	httpimplement.Start(cfg)
}
