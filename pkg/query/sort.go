package query

import (
	"fmt"

	"github.com/salihguru/idiogo/pkg/ptr"
	"github.com/salihguru/idiogo/pkg/xopt"
	"gorm.io/gorm"
)

func GetOrder(isAsc bool) string {
	if isAsc {
		return SortAsc
	}
	return SortDesc
}

type SortCond struct {
	Key       string
	Skip      bool
	Direction string
	IsDefault bool
}

func sortExpr(key string, direction *string) string {
	if direction == nil || *direction == "" {
		return key
	}
	return fmt.Sprintf("%s %s", key, *direction)
}

func (s SortCond) Expr() string {
	return sortExpr(s.Key, ptr.String(s.Direction))
}

func SortDirect(key string, isAsc ...bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(sortExpr(key, ptr.String(GetOrder(xopt.Get(false, isAsc...)))))
	}
}

func Sort(conds []SortCond) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(conds) == 0 {
			return db
		}
		isApplied := false
		var defaultSort string
		for _, cond := range conds {
			if cond.IsDefault {
				defaultSort = cond.Expr()
			}
			if cond.Skip {
				continue
			}
			db = db.Order(cond.Expr())
			isApplied = true
			break
		}
		if !isApplied && defaultSort != "" {
			db = db.Order(defaultSort)
		}
		return db
	}
}

func SortGeo(k string, long *float64, lat *float64, skip bool, isDefault ...bool) SortCond {
	return SortCond{
		Key:       OrderGeo("start_point", ptr.Float64Ref(long), ptr.Float64Ref(lat)),
		Skip:      skip || long == nil || lat == nil,
		IsDefault: getOption(false, isDefault...),
	}
}

func SortBasic(k, dir string, skip bool, isDefault ...bool) SortCond {
	return SortCond{
		Key:       k,
		Direction: dir,
		Skip:      skip,
		IsDefault: getOption(false, isDefault...),
	}
}
