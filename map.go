package chart

type ChartsMap struct {
	m map[string]*Chart
}

func (c *ChartsMap) Add(charts ...*Chart) {
	if c.m == nil {
		c.m = map[string]*Chart{}
	}

	for _, chart := range charts {
		c.m[chart.ID] = chart
	}
}

func (c *ChartsMap) Get(uid string) (chart *Chart, ok bool) {
	if c.m != nil {
		chart, ok = c.m[uid]
	}
	return
}
