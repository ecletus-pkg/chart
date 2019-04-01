package chart

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strings"

	"github.com/moisespsena-go/error-wrap"

	"github.com/ecletus-pkg/admin"
	"github.com/ecletus/admin"
	"github.com/ecletus/core"
	"github.com/moisespsena-go/default-logger"
	"github.com/moisespsena-go/path-helpers"
)

var (
	PKG = path_helpers.GetCalledDir()
	log = defaultlogger.NewLogger(PKG)
)

type Charts struct {
	Factory *OutputChartFactory
}

func (c *Charts) RegisterFuncMaps(register func(string, interface{})) {
	register("charts", func() *Charts {
		return c
	})
	register("charts__admin_url", func(ctx *admin.Context) string {
		return ctx.GenURL(c.Factory.Config.Uri.AdminUri)
	})
	register("charts__site_url", func(ctx *core.Context) string {
		return ctx.GenURL(c.Factory.Config.Uri.AdminUri)
	})
	register("charts__get", func(id string, adminName ...string) (chart *Chart) {
		if len(adminName) > 0 {
			chart, _ = c.Factory.AdminCharts[adminName[0]].Get(id)
		} else {
			chart, _ = c.Factory.SiteCharts.Get(id)
		}
		return chart
	})
	register("charts__opt", func(id string, key string, adminName ...string) interface{} {
		var chart *Chart
		if len(adminName) > 0 {
			chart, _ = c.Factory.AdminCharts[adminName[0]].Get(id)
		} else {
			chart, _ = c.Factory.SiteCharts.Get(id)
		}
		return chart.Option(key)
	})
	register("charts__opt_html", func(id string, key string, adminName ...string) template.HTMLAttr {
		var chart *Chart
		if len(adminName) > 0 {
			chart, _ = c.Factory.AdminCharts[adminName[0]].Get(id)
		} else {
			chart, _ = c.Factory.SiteCharts.Get(id)
		}
		op := chart.Option(key)
		var w = bytes.Buffer{}
		v := json.NewEncoder(&w)
		err := v.Encode(op)
		if err != nil {
			return template.HTMLAttr(err.Error())
		}
		r := w.String()
		return template.HTMLAttr(r)
	})
	register("charts__js", func() []string {
		return []string{"js/date.js", "vendors/chartjs.org/Chart.js", "chart/ecletus-charts.js"}
	})
}

type Chart struct {
	ChartFactory     `json:"-"`
	ID               string
	ChartLabel       string
	ChartDescription string
	Config           *ChartConfig
	IsSite           bool
	AdminName        string
}

func (c *Chart) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func (c *Chart) Label() string {
	return c.ChartLabel
}

func (c *Chart) Description() string {
	return c.ChartDescription
}

func (c *Chart) GetInfo(infos ...ChartInfo) (label, description string) {
	label, description = c.ChartLabel, c.ChartDescription
	for _, info := range infos {
		if info == nil {
			continue
		}
		if s := info.Label(); s != "" {
			label = s
		}
		if s := info.Description(); s != "" {
			description = s
		}
	}
	return
}

func (c *Chart) Option(key string) interface{} {
	opt := c.Config.Options
	parts := strings.Split(key, ".")

	for _, key := range parts[0 : len(parts)-1] {
		if optv, ok := opt[key]; ok {
			opt = optv.(map[string]interface{})
		}
		return nil
	}
	return opt[parts[len(parts)-1]]
}

func (c *Chart) OptionMap(key string) map[string]interface{} {
	opt := c.Config.Options
	parts := strings.Split(key, ".")

	for _, key := range parts[0 : len(parts)-1] {
		if optv, ok := opt[key]; ok {
			opt = optv.(map[string]interface{})
		}
		return nil
	}
	return opt[parts[len(parts)-1]].(map[string]interface{})
}

func (c *Chart) Find(ctx *Context) (out *OutputChart, err error) {
	finder, err := c.Factory(ctx)
	if err != nil {
		err = errwrap.Wrap(err, "Chart finder factory")
		return
	}
	data, err := finder.Find(ctx)
	if err != nil {
		err = errwrap.Wrap(err, "Chart find")
		return
	}

	var info ChartInfo
	if finfo, ok := finder.(ChartInfo); ok {
		info = finfo
	}
	label, description := c.GetInfo(info)
	return &OutputChart{ctx.Context.Ts(label), ctx.Context.Ts(description), data}, nil
}

func NewChart(cfg *ChartConfig, factory ChartFactory) *Chart {
	c := &Chart{ID: cfg.ID, ChartFactory: factory, Config: cfg}
	c.ChartLabel = cfg.Label
	c.IsSite = cfg.Site
	c.AdminName = cfg.AdminName
	if !c.IsSite && c.AdminName == "" {
		c.AdminName = admin_plugin.DEFAULT_ADMIN
	}
	return c
}
