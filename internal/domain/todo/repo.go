package todo

import (
	"context"

	"github.com/google/uuid"
	"github.com/salihguru/idiogo/pkg/list"
	"github.com/salihguru/idiogo/pkg/query"
	"github.com/salihguru/idiogo/pkg/xrepo"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Save(ctx context.Context, todo *Todo) error {
	return xrepo.Save(ctx, r.db, todo, todo.ID)
}

func (r *Repo) View(ctx context.Context, id uuid.UUID) (*Todo, error) {
	return xrepo.ViewByID[Todo](ctx, r.db, id)
}

func (r *Repo) Find(ctx context.Context, f Filters, pagi list.PagiRequest) ([]*Todo, error) {
	return xrepo.Find[*Todo](ctx, r.db,
		query.Apply(r.conds(f)),
		list.Paginate(&pagi),
	)
}

func (r *Repo) conds(f Filters) []query.Conds {
	return []query.Conds{
		query.ILike("title", f.Q),
		query.Eq("status", f.Status, f.Status == ""),
	}
}
