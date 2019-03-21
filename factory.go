package chart

import (
	"time"

	"github.com/ecletus/admin"
	"github.com/moisespsena/go-error-wrap"
)

type ResourceFactory struct {
	Resource *admin.Resource
	Charts   map[string]*Chart
}

func (f *ResourceFactory) RegisterChart(cfg *ChartConfig) *Chart {
	if f.Charts == nil {
		f.Charts = map[string]*Chart{}
	}

	factory := ChartFuncFactory(func(ctx *Context) (Finder, error) {
		event := NewReportFactoryEvent(ctx)
		event.With(E_REPORT_FACTORY + "@" + cfg.ResourceName)
		if err := f.Resource.Trigger(event); err != nil {
			return nil, errwrap.Wrap(err, "Record Report factory of %q", f.Resource.ID)
		}
		return event.Finder, nil
	})
	chart := NewChart(cfg, factory)
	f.Charts[cfg.ID] = chart
	return chart
}

type OutputChartFactory struct {
	Config      *Config
	Resources   map[string]map[string]*ResourceFactory
	Funcs       map[string]ChartFuncFactory
	AdminCharts map[string]*ChartsMap
	SiteCharts  ChartsMap
}

func NewOutputChartFactory(config *Config) (f *OutputChartFactory) {
	return &OutputChartFactory{
		Config:      config,
		Resources:   map[string]map[string]*ResourceFactory{},
		Funcs:       map[string]ChartFuncFactory{},
		AdminCharts: map[string]*ChartsMap{},
	}
}

func (f *OutputChartFactory) LoadFuncs() {
	for _, chartConfig := range f.Config.Charts {
		if chartConfig.ResourceName == "" {
			f.RegisterFunc(chartConfig)
		}
	}
}

func (f *OutputChartFactory) RegisterFunc(cfg *ChartConfig) bool {
	if fun, ok := registeredFuncs[cfg.Factory]; ok {
		f.Funcs[cfg.ID] = fun
		chart := NewChart(cfg, fun)
		if chart.IsSite {
			f.SiteCharts.Add(chart)
		} else {
			if f.AdminCharts[chart.AdminName] == nil {
				f.AdminCharts[chart.AdminName] = &ChartsMap{}
			}
			f.AdminCharts[chart.AdminName].Add(chart)
		}
		return true
	}
	log.Warningf("Func %q is not registered", cfg.ID)
	return false
}

func (f *OutputChartFactory) GetOrRegisterResource(res *admin.Resource) (rf *ResourceFactory) {
	rf = &ResourceFactory{Resource: res}
	adminName := res.GetAdmin().Name
	if _, ok := f.Resources[adminName]; !ok {
		f.Resources[adminName] = map[string]*ResourceFactory{}
	}
	f.Resources[adminName][res.ID] = rf
	return
}

func (f *OutputChartFactory) LoadAdminResources(Admin *admin.Admin) {
	for _, cfg := range f.Config.Charts {
		if cfg.ResourceName == "" || cfg.AdminName != Admin.Name {
			continue
		}

		if res := Admin.GetResourceByID(cfg.Factory); res != nil {
			rf := f.GetOrRegisterResource(res)
			chart := rf.RegisterChart(cfg)
			if f.AdminCharts[Admin.Name] == nil {
				f.AdminCharts[Admin.Name] = &ChartsMap{}
			}
			f.AdminCharts[Admin.Name].Add(chart)
		} else {
			log.Warningf("Admin %q does not have resource %q", Admin.Name, cfg.Factory)
		}
	}
}

func (f *OutputChartFactory) AdminFactory(context *admin.Context, startDate, endDate time.Time, chartsUID ...string) (result map[string]*OutputChart) {
	if len(chartsUID) == 0 {
		return
	}

	var (
		ctx = &Context{
			Errors:       &context.Errors,
			Context:      context.Context,
			AdminContext: context,
			StartDate:    startDate,
			EndDate:      endDate,
		}
		chart *Chart
		ok    bool
	)
	result = make(map[string]*OutputChart)

	for _, chartUID := range chartsUID {
		if chart, ok = f.AdminCharts[context.Admin.Name].Get(chartUID); !ok {
			if chart, ok = f.SiteCharts.Get(chartUID); !ok {
				continue
			}
		}

		ctx.Chart = chart
		out, err := chart.Find(ctx)
		if err != nil {
			context.AddError(err)
			context.LogErrors()
			continue
		}
		result[chart.ID] = out
	}

	return
}
