package chart

type ChartFactory interface {
	Factory(ctx *Context) (Finder, error)
}

type ChartFuncFactory func(ctx *Context) (Finder, error)

func (f ChartFuncFactory) Factory(ctx *Context) (Finder, error) {
	return f(ctx)
}
