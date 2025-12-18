package serve

import (
	"context"
	"sync"
	"time"

	"github.com/salihguru/idiogo/internal/config"
	"github.com/salihguru/idiogo/pkg/cancel"
	"github.com/salihguru/idiogo/pkg/i18np"
	"github.com/salihguru/idiogo/pkg/validation"
)

type App struct {
	Modules Modules
	Deps    Depends
	Config  config.Config
}

var once sync.Once
var instance *App

func Init(ctx context.Context) error {
	var err error
	once.Do(func() {
		var configs config.Config
		var i18n *i18np.I18n
		if err = config.Bind(&configs, "./config.yaml"); err != nil {
			return
		}
		i18n, err = i18np.New(i18np.Config{})
		if err != nil {
			return
		}
		i18n.Load(configs.I18n.Dir, configs.I18n.Locales...)
		deps := Depends{
			I18n:          i18n,
			ValidationSrv: validation.New(i18n),
		}
		if err = deps.Up(ctx, configs); err != nil {
			return
		}
		instance = &App{
			Modules: newModules(&deps),
			Deps:    deps,
			Config:  configs,
		}
	})
	return err
}

func Get() *App {
	return instance
}

type disconFunc func(context.Context) error

func (a *App) disconnectAll(ctx context.Context, fns ...disconFunc) error {
	for _, fn := range fns {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context, fns ...disconFunc) error {
	return cancel.NewWithTimeout(ctx, 5*time.Second, func(ctx context.Context) error {
		fns = append(fns, a.Deps.Shutdown)
		return a.disconnectAll(ctx, fns...)
	})
}
