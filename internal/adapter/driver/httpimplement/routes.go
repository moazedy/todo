package httpimplement

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/moazedy/todo/internal/adapter/driven/db/repoimplement"
	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/internal/domain/srvimplement"
	"github.com/moazedy/todo/internal/port/driven/db/repository"
	"github.com/moazedy/todo/internal/port/driver/http"
	"github.com/moazedy/todo/internal/port/driver/service"
	"github.com/moazedy/todo/pkg/infra/config"
	"github.com/moazedy/todo/pkg/infra/tx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type serverItems struct {
	// basics
	dbConnection *gorm.DB
	txFactory    tx.TXFactory

	// repo layer
	todoItemRepoFactory repository.GenericRepoFactory[model.TodoItem]

	// service layer
	todoService service.TodoItem

	// presentation layer
	todoController http.TodoItem
}

var items serverItems

func register(app *echo.Echo, cfg config.Config) {
	items.dbConnection = initDB(cfg.Postgres)
	items.txFactory = tx.NewTXFactory(items.dbConnection)

	items.todoItemRepoFactory = repoimplement.NewGenericRepoFactory[model.TodoItem]()

	items.todoService = srvimplement.NewTodoItem(items.txFactory, items.todoItemRepoFactory)

	items.todoController = NewTodoItem(items.todoService)

	app.POST("/todo/item", items.todoController.Create)
	app.PUT("/todo/item", items.todoController.Update)
	app.DELETE("/todo/item/:id", items.todoController.Delete)
	app.GET("/todo/item/:id", items.todoController.GetByID)
}

// opening connection with database
func initDB(cfg config.PostgresConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.ToString()), &gorm.Config{})
	if err != nil {
		println(err)
		panic("failed to connect with db")
	}

	// Check if the desired database exists
	err = db.Exec(
		fmt.Sprintf("CREATE DATABASE %s",
			cfg.Name,
		)).Error
	if err != nil {
		println(err.Error())
	}

	// Reconnect to the newly created database
	db, err = gorm.Open(postgres.Open(cfg.ToStringWithDbName()), &gorm.Config{})
	if err != nil {
		println(err)
		panic("failed to connect with db")
	}

	err = db.Exec(
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	).Error
	if err != nil {
		panic(err.Error())
	}

	// NOTE : register all entities in here
	err = db.AutoMigrate(
		&model.TodoItem{},
	)
	if err != nil {
		println(err.Error())
		panic("failed to migrate")
	}

	return db
}
