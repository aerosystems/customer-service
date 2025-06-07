package adapters

import (
	"fmt"

	"gorm.io/gorm"
)

type Migration struct {
	db *gorm.DB
}

func NewMigration(db *gorm.DB) *Migration {
	return &Migration{
		db: db,
	}
}

func (m *Migration) Run() error {
	return autoMigrateGORM(m.db)
}

func autoMigrateGORM(db *gorm.DB) error {
	if err := db.AutoMigrate(Customer{}); err != nil {
		return fmt.Errorf("failed to autoMigrateGORM: %v", err)
	}
	return nil
}
