package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"default:null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" swaggertype:"string" format:"date-time"`
}

type Headless struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"default:null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" swaggertype:"string" format:"date-time"`
}

type ViewBase struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	return nil
}

func (h *Headless) BeforeCreate(tx *gorm.DB) (err error) {
	if h.CreatedAt.IsZero() {
		h.CreatedAt = time.Now()
	}
	return nil
}

func DeleteNow() gorm.DeletedAt {
	return gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
}

func Restore() gorm.DeletedAt {
	return gorm.DeletedAt{}
}
