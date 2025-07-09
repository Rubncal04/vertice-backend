package migrations

import (
	"vertice-backend/internal/domain"

	"gorm.io/gorm"
)

func AutoMigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
	)
}
