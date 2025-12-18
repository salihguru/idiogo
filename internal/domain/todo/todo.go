package todo

import "github.com/salihguru/idiogo/pkg/entity"

type Todo struct {
	entity.Base
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
	StatusArchived  Status = "archived"
)
