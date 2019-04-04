package chart

import "github.com/moisespsena-go/edis"

var (
	E_REPORT_FACTORY  = PKG + ".reportFactory"
	E_REPORTS_FACTORY = PKG + ".reportsFactory"
)

type ChartReportEvent struct {
	edis.EventInterface
	Context *Context
}

type ChartReportFactoryEvent struct {
	*ChartReportEvent
	Finder Finder
}

type ChartReportsFactoryEvent struct {
	*ChartReportEvent
	Finder Finder
}

type events struct {
	dis edis.EventDispatcherInterface
}

func Events(dis edis.EventDispatcherInterface) *events {
	return &events{dis}
}

func (e events) OnReportFactory(cb func(e *ChartReportFactoryEvent) error) {
	e.dis.On(E_REPORT_FACTORY, func(e edis.EventInterface) error {
		return cb(e.(*ChartReportFactoryEvent))
	})
}

func NewReportFactoryEvent(ctx *Context) *ChartReportFactoryEvent {
	return &ChartReportFactoryEvent{&ChartReportEvent{edis.NewEvent(E_REPORT_FACTORY), ctx}, nil}
}
