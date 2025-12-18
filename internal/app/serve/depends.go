package serve

import (
	"context"

	"github.com/salihguru/idiogo/internal/config"
	"github.com/salihguru/idiogo/internal/infra/db"
	"github.com/salihguru/idiogo/internal/infra/db/migration"
	"github.com/salihguru/idiogo/pkg/i18np"
	"github.com/salihguru/idiogo/pkg/validation"
	"gorm.io/gorm"
)

type Depends struct {
	DB            *gorm.DB
	ValidationSrv *validation.Srv
	I18n          *i18np.I18n
}

func (d *Depends) Up(ctx context.Context, cnf config.Config) error {
	db, err := db.NewPostgres(ctx, db.PostgresConfig{
		Host:     cnf.DB.Host,
		Port:     cnf.DB.Port,
		User:     cnf.DB.User,
		Password: cnf.DB.Pass,
		DBName:   cnf.DB.Name,
		SSLMode:  cnf.DB.SSLMode,
		Debug:    cnf.DB.Debug,
	})
	if err != nil {
		return err
	}
	if cnf.DB.Migrate {
		if err := migration.RunSql(ctx, db); err != nil {
			return err
		}
	}
	d.DB = db
	return nil
}

func (d Depends) Shutdown(ctx context.Context) error {
	return d.closeDB(ctx)
}

func (d Depends) closeDB(_ context.Context) error {
	if d.DB == nil {
		return nil
	}
	sql, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}
