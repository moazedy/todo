package tx

import (
	"fmt"

	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/pkg/infra/config"
	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

type tx struct {
	connection *gorm.DB
}

func (t *tx) GetConnection() *gorm.DB {
	return t.connection
}

func (t *tx) Commit() error {
	return t.connection.Commit().Error
}

func (t *tx) Rollback() error {
	return t.connection.Rollback().Error
}

func (t *tx) AutoCR(err error) error {
	if err == nil {
		commitErr := t.Commit()
		if commitErr != nil {
			return err
		}
		return nil
	} else {
		rollbackErr := t.Rollback()
		if rollbackErr != nil {
			println("error while tx rollback: ", rollbackErr)
		}

		return err
	}
}

// opening connection with database
func GetDB(cfg config.PostgresConfig) *gorm.DB {
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
