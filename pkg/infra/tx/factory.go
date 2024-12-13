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
	isMock     bool
	connection *gorm.DB
}

func NewTXFactory(isMock bool, db *gorm.DB) TXFactory {
	return txFactory{
		connection: db,
		isMock:     isMock,
	}
}

func (tf txFactory) NewTX() TX {
	if tf.isMock {
		return mockTx{}
	} else {
		return &tx{
			connection: tf.connection.Begin(),
		}
	}
}
