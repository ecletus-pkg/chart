(function(factory) {
    if (typeof define === 'function' && define.amd) {
        // AMD. Register as anonymous module.
        define(['jquery'], factory);
    } else if (typeof exports === 'object') {
        // Node / CommonJS
        factory(require('jquery'));
    } else {
        // Browser globals.
        factory(jQuery);
    }
})(function($) {

    'use strict';

    var NAMESPACE = 'aghaphe.chart';
    var EVENT_CLICK = 'click.' + NAMESPACE;
    var CHARTS_SELECTOR = '.aghaphe-charts';
    var CHART_SELECTOR = '.aghaphe-chart';
    var EVENT_ENABLE = 'enable.' + NAMESPACE;
    var EVENT_DISABLE = 'disable.' + NAMESPACE;

    var yesterday = (new Date()).AddDate(-1);

    function AghapheCharts(element, options) {
        this.$element = $(element);
        this.dataSourceURL = this.$element.attr('data-source-url');
        this.data = this.$element.data();
        this.options = $.extend(true, {}, AghapheCharts.DEFAULTS, $.isPlainObject(options) && options);
        this.$startDate = this.$element.find('[data-charts-value="start-date"]');
        this.$endDate = this.$element.find('[data-charts-value="end-date"]');
        this.init();
    }

    Chart.defaults.global.responsive = true;

    AghapheCharts.prototype = {
        constructor: AghapheCharts,

        init: function() {
            this.name = this.data.name;
            this.charts = {};
            this.bind();
            this.loadLevel("last_week");
        },

        bind: function() {
            this.$element
                .on(EVENT_CLICK, '[data-charts-action-load]', this.load.bind(this))
                .on(EVENT_CLICK, '[data-charts-action]', this.action.bind(this));
        },

        unbind: function() {
            this.$element
                .off(EVENT_CLICK, this.load)
                .off(EVENT_CLICK, this.action)
        },

        loadLevel: function(level) {
            console.log("loadLevel:", level);
            switch (level) {
                case "yesterday":
                    this.$startDate.val(yesterday.Format("yyyy-MM-dd"));
                    this.$endDate.val(yesterday.Format("yyyy-MM-dd"));
                    break;
                case "this_week":
                    let beginningOfThisWeek = yesterday.AddDate(-yesterday.getDay() + 1);
                    this.$startDate.val(beginningOfThisWeek.Format("yyyy-MM-dd"));
                    this.$endDate.val(beginningOfThisWeek.AddDate(6).Format("yyyy-MM-dd"));
                    break;
                case "last_week":
                    let endOfLastWeek = yesterday.AddDate(-yesterday.getDay());
                    this.$startDate.val(endOfLastWeek.AddDate(-6).Format("yyyy-MM-dd"));
                    this.$endDate.val(endOfLastWeek.Format("yyyy-MM-dd"));
                    break;
            }
            this.reload();
        },

        load: function(e) {
            let $this = $(e.currentTarget);
            this.loadLevel($this.attr('data-charts-action-load'));
            $this.blur();
        },

        reload: function() {
            let charts = new Array(),
                $charts = new Array(),
                _this = this;
            this.$element.find(CHART_SELECTOR).each(function () {
                let $this = $(this),
                    id = $this.attr('data-chart');
                if (id) {
                    $charts[$charts.length] = $this;
                    charts[charts.length] = id;
                }
            });
            if (charts.length > 0){
                $.getJSON(this.dataSourceURL, {
                    startDate: this.$startDate.val(),
                    endDate: this.$endDate.val(),
                    charts:charts.join('+')
                }, function (jsonData) {
                    $charts.forEach(function ($chart) {
                        let id = $chart.attr('data-chart');
                        if (jsonData[id]) {
                            _this.renderChart($chart, jsonData[id]);
                        }
                    })
                    //RenderChart(jsonData);
                });
            }
        },

        renderChart: function($chart, data) {
            let id = $chart.attr('data-chart'),
                chart,
                $context,
                dataSet;
            if (this.charts[id]) this.charts[id].destroy();
            var labels = [];
            var counts = [];
            for (var i = 0; i < data.Result.length; i++) {
                labels.push(data.Result[i].Date.substring(5,10));
                counts.push(data.Result[i].Total)
            }

            $context = $chart.find('canvas');
            if ($context.length =0) return;

            dataSet = $.trim($chart.data().set || '{}');
            eval('dataSet = ' + dataSet + ';');
            dataSet = $.extend({},  AghapheCharts.DEFAULTS.dataSet, dataSet, {
                label:data.Label,
                data: counts,
            });

            data = {
                type: $chart.data('type') || AghapheCharts.DEFAULTS.type,
                data: {
                    labels: labels,
                    datasets: [dataSet]
                }
            };
            chart = new Chart($context[0].getContext("2d"), data);
            this.charts[id] = chart;
        },

        action: function(e) {
            switch ($(e.currentTarget).attr('data-charts-action')) {
                case "reload":
                    this.reload();
            }
        },

        destroy: function() {
            this.unbind();
            for (let k in this.charts) {
                this.charts[k].destroy();
            }
            this.charts[k] = null;
        }
    };

    AghapheCharts.DEFAULTS = {
        sourceUrl:"reports.json",
        chartsName:"",
        type: "line",
        dataSet: {
            fillColor: "rgba(151,187,205,0.2)",
            strokeColor: "rgba(151,187,205,1)",
            pointColor: 'rgb(255, 159, 64)',
            pointStrokeColor: "#fff",
            pointHighlightFill: "#fff",
            pointHighlightStroke: "rgba(151,187,205,1)",
        }
    };

    AghapheCharts.COLORS = {
        red: 'rgb(255, 99, 132)',
        orange: 'rgb(255, 159, 64)',
        yellow: 'rgb(255, 205, 86)',
        green: 'rgb(75, 192, 192)',
        blue: 'rgb(54, 162, 235)',
        purple: 'rgb(153, 102, 255)',
        grey: 'rgb(201, 203, 207)'
    };

    AghapheCharts.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }

                $this.data(NAMESPACE, (data = new AghapheCharts(this, options)));
            }

            if (typeof options === 'string' && $.isFunction(fn = data[options])) {
                fn.apply(data);
            }
        });
    };

    $(function() {
        $(document).
        on(EVENT_DISABLE, function(e) {
            AghapheCharts.plugin.call($(CHARTS_SELECTOR, e.target), 'destroy');
        }).
        on(EVENT_ENABLE, function(e) {
            AghapheCharts.plugin.call($(CHARTS_SELECTOR, e.target));
        }).
        triggerHandler(EVENT_ENABLE);
    });

    return AghapheCharts;

});