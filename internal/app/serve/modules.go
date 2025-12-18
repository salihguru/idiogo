package serve

import (
	"github.com/salihguru/idiogo/internal/domain/todo"
	"github.com/salihguru/idiogo/internal/rest"
)

type Modules struct {
	Todo rest.Module[*todo.Repo, *todo.Service]
}

func newModules(deps *Depends) Modules {
	todoRepo := todo.NewRepo(deps.DB)
	todoSrv := todo.NewService(todoRepo)
	return Modules{
		Todo: rest.Module[*todo.Repo, *todo.Service]{
			Repo:    todoRepo,
			Service: todoSrv,
			Router:  todo.NewHandler(*todoSrv),
		},
	}
}

func (m Modules) Routers() []rest.Router {
	return []rest.Router{
		m.Todo.Router,
	}
}
