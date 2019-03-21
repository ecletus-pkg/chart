package sample

import (
	"github.com/ecletus-pkg/chart"
	"github.com/moisespsena/go-path-helpers"
)

func init() {
	chart.RegisterFuncFactory(path_helpers.GetCalledDir()+".SampleFuncFactory", func(ctx *chart.Context) (chart.Finder, error) {
		min, max := 0, 2000
		if v, ok := ctx.Chart.Config.Options["min"]; ok {
			min = v.(int)
		}
		if v, ok := ctx.Chart.Config.Options["max"]; ok {
			max = v.(int)
		}
		return chart.NewFinderInfo(&chart.FakeFinder{min, max},
			"Sample Chart",
			"The Sample Chart description"), nil
	})
}
