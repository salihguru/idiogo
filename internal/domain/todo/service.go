package todo

import (
	"context"

	"github.com/google/uuid"
	"github.com/salihguru/idiogo/pkg/entity"
	"github.com/salihguru/idiogo/pkg/list"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

type CreateReq struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=5000"`
}

type UpdateReq struct {
	ID          uuid.UUID `params:"id" validate:"required,uuid"`
	Title       *string   `json:"title" validate:"omitempty,min=3,max=255"`
	Description *string   `json:"description" validate:"omitempty,max=5000"`
	Status      *string   `json:"status" validate:"omitempty,oneof=pending completed cancelled archived"`
}

type ViewReq struct {
	ID uuid.UUID `params:"id" validate:"required,uuid"`
}

type ListReq struct {
	Filters
	list.PagiRequest
}

func (s *Service) Create(ctx context.Context, req CreateReq) (*Todo, error) {
	todo := &Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      StatusPending,
	}
	if err := s.repo.Save(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *Service) Update(ctx context.Context, req UpdateReq) (*Todo, error) {
	todo, err := s.repo.View(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Status != nil {
		todo.Status = Status(*req.Status)
	}
	if err := s.repo.Save(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *Service) View(ctx context.Context, req ViewReq) (*Todo, error) {
	return s.repo.View(ctx, req.ID)
}

func (s *Service) Find(ctx context.Context, req ListReq) ([]*Todo, error) {
	return s.repo.Find(ctx, req.Filters, req.PagiRequest)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	todo, err := s.repo.View(ctx, id)
	if err != nil {
		return err
	}
	todo.Status = StatusArchived
	todo.DeletedAt = entity.DeleteNow()
	return s.repo.Save(ctx, todo)
}
