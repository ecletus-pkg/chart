package chart

import (
	"os"

	"github.com/ecletus-pkg/admin"
	"github.com/ecletus/plug"
	"github.com/ecletus/router"
)

type Plugin struct {
	admin_plugin.AdminNames
	plug.EventDispatcher
	ConfigDir string
	RouterKey string
	ChartsKey string
	config    *Config
	charts    *Charts
}

func (p *Plugin) RequireOptions() []string {
	return []string{p.RouterKey}
}

func (p *Plugin) ProvideOptions() []string {
	return []string{p.ChartsKey}
}

func (p *Plugin) Init(options *plug.Options) error {
	cfg, err := LoadConfigDir(p.ConfigDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if cfg == nil {
		cfg = &Config{}
	}
	cfg.LoadDefaults()
	p.config = cfg
	p.charts = &Charts{NewOutputChartFactory(cfg)}
	options.Set(p.ChartsKey, p.charts)
	return nil
}

func (p *Plugin) OnRegister() {
	admin_plugin.Events(p).
		Router(func(e *admin_plugin.AdminRouterEvent) {
			p.charts.RegisterAdminRoutes(e.Admin)
		}).
		FuncMap(func(e *admin_plugin.AdminFuncMapEvent) {
			p.charts.RegisterFuncMaps(e.Register)
		}).
		Done(func(e *admin_plugin.AdminEvent) {
			p.charts.Factory.LoadAdminResources(e.Admin)
		})

	router.OnRoute(p, func(e *router.RouterEvent) {
		p.charts.RegisterSiteRoutes(e.Router.Mux)
	})

	plug.OnPostInit(p, func() {
		p.charts.Factory.LoadFuncs()
	})
}
