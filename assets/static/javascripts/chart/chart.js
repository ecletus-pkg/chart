var CHARTS = {};

function NewChart(name, sourceData) {
    var labels = [];
    var counts = [];
    for (var i = 0; i < sourceData.Result.length; i++) {
        labels.push(sourceData.Result[i].Date.substring(5,10));
        counts.push(sourceData.Result[i].Total)
    }
    var el = document.getElementById("chart__" + name.toLowerCase());

    if (el) {
        var context = document.getElementById("chart__" + name.toLowerCase()).getContext("2d");
        var data = ChartData(sourceData.Label, labels, counts);
        return new Chart(context).Bar(data, "");
    }
}

function RenderChart(chartsSources) {
    Chart.defaults.global.responsive = true;
    var chart;

    for (var name in chartsSources) {
        if (CHARTS[name]) {
            CHARTS[name].destroy();
        }
        chart = NewChart(name, chartsSources[name]);
        if (chart) {
            CHARTS[name] = chart;
        }
    }
}

function ChartData(label, lables, counts) {
    return {
        labels: lables,
        datasets: [
            {
                label: label,
                fillColor: "rgba(151,187,205,0.2)",
                strokeColor: "rgba(151,187,205,1)",
                pointColor: "rgba(151,187,205,1)",
                pointStrokeColor: "#fff",
                pointHighlightFill: "#fff",
                pointHighlightStroke: "rgba(151,187,205,1)",
                data: counts
            }
        ]
    };
}