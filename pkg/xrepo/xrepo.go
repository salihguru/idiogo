package xrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/restayway/stx"
	"gorm.io/gorm"
)

type ScopeFunc = func(*gorm.DB) *gorm.DB

func WithContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	if txdb := stx.Current(ctx); txdb != nil {
		return txdb
	}
	return db.WithContext(ctx)
}

func ViewByWhere[T any](ctx context.Context, db *gorm.DB, where string, args ...interface{}) (*T, error) {
	var entity T
	if err := WithContext(ctx, db).Where(where, args...).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func ViewByID[T any](ctx context.Context, db *gorm.DB, id uuid.UUID) (*T, error) {
	return ViewByWhere[T](ctx, db, "id = ?", id)
}

func Find[T any](ctx context.Context, db *gorm.DB, scopes ...ScopeFunc) ([]T, error) {
	var entities []T
	if err := WithContext(ctx, db).Scopes(scopes...).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func Save[T any](ctx context.Context, db *gorm.DB, entity *T, id uuid.UUID) error {
	if id == uuid.Nil {
		return WithContext(ctx, db).Create(entity).Error
	}
	return WithContext(ctx, db).Where("id = ?", id).Save(entity).Error
}
