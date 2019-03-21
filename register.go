package chart

var registeredFuncs = map[string]ChartFuncFactory{}

func RegisterFuncFactory(id string, factory ChartFuncFactory) {
	registeredFuncs[id] = factory
}
