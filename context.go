package chart

import (
	"time"

	"github.com/aghape/admin"

	"github.com/aghape/core"
)

type Context struct {
	*core.Errors

	Chart        *Chart
	Context      *core.Context
	AdminContext *admin.Context
	StartDate    time.Time
	EndDate      time.Time
}
