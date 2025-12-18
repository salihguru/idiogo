package migration

import (
	"context"

	"github.com/salihguru/idiogo/internal/domain/todo"
	"gorm.io/gorm"
)

func RunSql(ctx context.Context, db *gorm.DB) error {
	err := db.AutoMigrate(
		&todo.Todo{},
	)
	if err != nil {
		return err
	}

	return nil
}
