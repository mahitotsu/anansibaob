package db

import (
	"github.com/mahitotsu/anansibaob/internal/models"
	"gorm.io/gorm"
)

func DropAndCreate(db *gorm.DB) error {

	tables, err := db.Migrator().GetTables()
	if err != nil {
		return err
	}
	for table := range tables {
		db.Migrator().DropTable(table)
	}

	return db.Migrator().AutoMigrate(models.AllModels)
}
