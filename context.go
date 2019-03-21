package chart

import (
	"time"

	"github.com/ecletus/admin"

	"github.com/ecletus/core"
)

type Context struct {
	*core.Errors

	Chart        *Chart
	Context      *core.Context
	AdminContext *admin.Context
	StartDate    time.Time
	EndDate      time.Time
}
