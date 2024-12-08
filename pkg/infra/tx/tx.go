package tx

import (
	"gorm.io/gorm"
)

type TX interface {
	Commit() error
	Rollback() error
	AutoCR(error) error
	GetConnection() *gorm.DB
}

type TXFactory interface {
	NewTX() TX
}

type txFactory struct {
	connection *gorm.DB
}

func NewTXFactory(db *gorm.DB) TXFactory {
	return txFactory{
		connection: db,
	}
}

type tx struct {
	connection *gorm.DB
}

func (tf txFactory) NewTX() TX {
	return &tx{
		connection: tf.connection.Begin(),
	}
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
