package chart

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ecletus/admin"
	"github.com/moisespsena-go/xroute"
)

func (c *Charts) AdminHandler(ctx *admin.Context) {
	var chartsList []string
	if charts := strings.TrimSpace(ctx.Request.URL.Query().Get("charts")); charts != "" {
		for _, chartId := range strings.Split(charts, " ") {
			if chartId != "" {
				chartsList = append(chartsList, chartId)
			}
		}
	}

	start := ctx.Request.URL.Query().Get("startDate")
	end := ctx.Request.URL.Query().Get("endDate")
	startDate, endDate, err := ParseDateRange(start, end)

	if err != nil {
		ctx.AddError(err)
	} else if charts := c.Factory.AdminFactory(ctx, startDate, endDate, chartsList...); !ctx.HasError() {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(charts)
		ctx.Writer.Write(b)
		return
	}

	ctx.Writer.WriteHeader(http.StatusInternalServerError)
	ctx.Writer.Write([]byte(ctx.Errors.String()))
}

func (c *Charts) RegisterAdminRoutes(Admin *admin.Admin) {
	Admin.OnRouter(func(r xroute.Router) {
		r.Get(c.Factory.Config.Uri.AdminUri, admin.NewHandler(c.AdminHandler))
	})
}

func (c *Charts) SiteHandler(w http.ResponseWriter, r *http.Request, rctx *xroute.RouteContext) {
	w.Write([]byte("hello to charts handler!"))
}

func (c *Charts) RegisterSiteRoutes(router *xroute.Mux) {
	router.Get(c.Factory.Config.Uri.SiteUri, c.SiteHandler)
}
