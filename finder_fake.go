package chart

import (
	"fmt"
)

type FakeFinder struct {
	Min int
	Max int
}

func (f *FakeFinder) Find(ctx *Context) (res []Record, err error) {
	days := int(ctx.EndDate.Sub(ctx.StartDate).Hours() / 24)
	res = make([]Record, days)

	Random(f.Min, f.Max, func(i, v int) bool {
		if i < days {
			res[i] = Record{fmt.Sprint(v), ctx.StartDate.AddDate(0, 0, i)}
			return true
		}
		return false
	})
	return
}
