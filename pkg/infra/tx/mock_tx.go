package tx

import "gorm.io/gorm"

type mockTx struct{}

func (t mockTx) GetConnection() *gorm.DB {
	return nil
}

func (t mockTx) Commit() error {
	return nil
}

func (t mockTx) Rollback() error {
	return nil
}

func (t mockTx) AutoCR(err error) error {
	return nil
}
