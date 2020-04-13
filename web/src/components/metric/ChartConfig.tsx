
import { get } from 'lodash';
import { ChartStatus, ChartStatusEnum } from 'model/Chart';
import { UnitEnum } from 'model/Metric';
import { DataFormatter } from 'utils/DataFormatter';

const ChartJS = require('chart.js')

ChartJS.defaults.global.elements.line.borderWidth = 0;
ChartJS.defaults.global.elements.line.cubicInterpolationMode = "monotone";

export const INTERVALS = [10000, 30000, 60 * 1000, 2 * 60 * 1000, 5 * 60 * 1000, 10 * 60 * 1000, 15 * 60 * 1000, 30 * 60 * 1000];
export const EXTRA_PLOTLINES = "extraPlotLines";

export function isNoData(chartStatus: ChartStatus) {
    if (!chartStatus || chartStatus.status) {
        return false;
    }
    const status:any = chartStatus.status
    return [ChartStatusEnum.Loaded, ChartStatusEnum.UnLimit, ChartStatusEnum.Loading].indexOf(status) < 0;
}

export function getVisibleInterval(times: any, interval: number, width: number) {
    const ticksCount = times.length;
    const perTickPX = width / (ticksCount - 1);
    let ticksPerVisibleTick = parseInt("" + 100 / perTickPX);
    const t = ticksPerVisibleTick * interval;
    let result = interval;
    if (t < 60 * 60 * 1000) {
        INTERVALS.forEach(item => {
            if (t / item < 0.6) {
                return;
            } else {
                result = item;
            }
        });
    } else {
        result = t - t % (60 * 60 * 1000);
    }
    return result;
}

ChartJS.plugins.register({
    afterUpdate: function (chart: any) {
        const xAxes = chart.scales["xAxes-bottom"];
        const tickOffset = get(xAxes, "options.ticks.tickOffset", null);
        const display = get(xAxes, "options.display", false);
        if (display && tickOffset) {
            const width = get(chart, "scales.xAxes-bottom.width", 0);
            const interval = get(chart, "config.options.interval", 0);
            const times = get(chart, "config.options.times", []);
            const visibleInterval = getVisibleInterval(times, interval, width);

            xAxes.draw = function () {
                const xScale = chart.scales["xAxes-bottom"];
                const helpers = ChartJS.helpers;

                const tickFontColor = helpers.getValueOrDefault(xScale.options.ticks.fontColor, ChartJS.defaults.global.defaultFontColor);
                const tickFontSize = helpers.getValueOrDefault(xScale.options.ticks.fontSize, ChartJS.defaults.global.defaultFontSize);
                const tickFontStyle = helpers.getValueOrDefault(xScale.options.ticks.fontStyle, ChartJS.defaults.global.defaultFontStyle);
                const tickFontFamily = helpers.getValueOrDefault(xScale.options.ticks.fontFamily, ChartJS.defaults.global.defaultFontFamily);
                const tickLabelFont = helpers.fontString(tickFontSize, tickFontStyle, tickFontFamily);
                const tl = xScale.options.gridLines.tickMarkLength;

                const isRotated = xScale.labelRotation !== 0;
                const yTickStart = xScale.top;
                const yTickEnd = xScale.top + tl;
                const chartArea = chart.chartArea;

                helpers.each(
                    xScale.ticks,
                    (label: any, index: any) => {
                        if (times[index] % visibleInterval !== 0) {
                            return;
                        }
                        // console.log("xxxxxx",index,times,visibleInterval)

                        // copy of chart.js code
                        let xLineValue = this.getPixelForTick(index);
                        const xLabelValue = this.getPixelForTick(index, this.options.gridLines.offsetGridLines);

                        if (this.options.gridLines.display) {

                            this.ctx.lineWidth = this.options.gridLines.lineWidth;
                            this.ctx.strokeStyle = this.options.gridLines.color;

                            xLineValue += helpers.aliasPixel(this.ctx.lineWidth);

                            // Draw the label area
                            this.ctx.beginPath();

                            if (this.options.gridLines.drawTicks) {
                                this.ctx.moveTo(xLineValue, yTickStart);
                                this.ctx.lineTo(xLineValue, yTickEnd);
                            }

                            // Draw the chart area
                            if (this.options.gridLines.drawOnChartArea) {
                                this.ctx.moveTo(xLineValue, chartArea.top);
                                this.ctx.lineTo(xLineValue, chartArea.bottom);
                            }

                            // Need to stroke in the loop because we are potentially changing line widths & colours
                            this.ctx.stroke();
                        }

                        if (this.options.ticks.display) {
                            this.ctx.save();
                            this.ctx.translate(xLabelValue + this.options.ticks.labelOffset, (isRotated) ? this.top + 12 : this.options.position === "top" ? this.bottom - tl : this.top + tl);
                            this.ctx.rotate(helpers.toRadians(this.labelRotation) * -1);
                            this.ctx.font = tickLabelFont;
                            this.ctx.textAlign = (isRotated) ? "right" : "center";
                            this.ctx.textBaseline = (isRotated) ? "middle" : this.options.position === "top" ? "bottom" : "top";
                            this.ctx.fillStyle = tickFontColor;
                            this.ctx.fillText(label, 0, 0);
                            this.ctx.restore();
                        }
                    },
                    xScale);
            };
        }
    },
    afterDraw: function (chart: any) {
        const status = get(chart, "options.status", null);

        const ctx = chart.chart.ctx;

        if (isNoData(status)) {
            chart.clear();
            const width = chart.chart.width;
            const height = chart.chart.height;
            let text = "";
            let color = get(chart, "options.scales.yAxes[0].ticks.fontColor", null);
            ctx.textAlign = "center";
            ctx.textBaseline = "middle";
            switch (status.status) {
                case ChartStatusEnum.NoData:
                    text = "No data to display";
                    break;
                case ChartStatusEnum.BadRequest:
                    text = "Invalid Configuration";
                    color = "#D27613";
                    break;
                case ChartStatusEnum.LoadError:
                    color = "#F56C6C";
                    ctx.fillStyle = color;
                    ctx.fillText(status.msg, width / 2, (height / 2) + 20);
                    text = "Internal Server Error";
                    break;
                default:
                    break;
            }
            ctx.font = "13px Arial";
            ctx.fillStyle = color;
            ctx.fillText(text, width / 2, height / 2);
            ctx.restore();
        } else if (chart.options.isSeriesChart) {

            const chartArea = chart.chartArea;

            ctx.beginPath();
            ctx.lineWidth = 1;
            ctx.strokeStyle = get(chart, "options.scales.yAxes[0].gridLines.color", null);
            ctx.moveTo(chartArea.left, chartArea.bottom);
            ctx.lineTo(chartArea.right, chartArea.bottom);
            ctx.stroke();
        }
    }
});

export const LEFT_Y_AXES = {
    id: "yAxes-left",
    position: "left",
    gridLines: {
        drawTicks: false,
        lineWidth: 0.3,
        autoSkip: true,
        tickMarkLength: 0,
        zeroLineWidth: 0,
        drawBorder: false,
        color: "#d2d2d2",
        // drawOnChartArea: false,
        // borderDash: [1, 1],
    },
    ticks: {
        fontColor: "#d2d2d2",
        mirror: true, // draw tick in chart area
        padding: 2,
        display: true,
        min: 0,
        fontSize: 12,
        autoSkip: true,
        tickMarkLength: 0,
        maxTicksLimit: 6
    },
};

export const RIGHT_Y_AXES = {
    id: "yAxes-right",
    position: "right",
    display: false,
    gridLines: {
        drawBorder: false,
        drawTicks: false,
        display: false,
    },
    ticks: {
        fontColor: "#d2d2d2",
        mirror: true, // draw tick in chart area
        padding: 2,
        display: false,
        min: 0,
        fontSize: 12,
        autoSkip: true,
        maxTicksLimit: 6
    }
};

export const CANVAS_CHART_CONFIG = {
    data: {},
    options: {
        responsive: true,
        maintainAspectRatio: false,
        responsiveAnimationDuration: 0, // animation duration after a resize
        legend: {
            display: false
        },
        elements: {
            line: {
                tension: 0, // disables bezier curve
                borderWidth: 1
            },
            point: {
                radius: 0
            },
            arc: {
                borderWidth: 0
            }
        },
        tooltips: {
            enabled: false,
            mode: "dataset",
        },
        hover: {
            animationDuration: 0, // duration of animations when hovering an item
            mode: "index",
            intersect: false,
        },
        animation: {
            duration: 0, // general animation time
        },
        zoom: {
            enabled: true,
            drag: true,
            mode: "x",
            limits: {
                max: 10,
                min: 0.5
            }
        },
        annotation: {
            events: ["click"],
            annotations: []
        },
        scales: {
            xAxes: [{
                id: "xAxes-bottom",
                display: true,
                type: "category",
                gridLines: {
                    drawTicks: false,
                    lineWidth: 0.3,
                    tickMarkLength: 2,
                    color: "#d2d2d2",
                    // drawOnChartArea: false,
                    drawBorder: false,
                },
                ticks: {
                    fontSize: 12,
                    fontColor: "#d2d2d2",
                    maxRotation: 0, // angle in degrees
                    tickOffset: 100
                }
            }],
            yAxes: [LEFT_Y_AXES, RIGHT_Y_AXES]
        },
    }
};


/**
 * Build yAxes config
 * @param options yAxes options
 * @param yAxesType type
 */
export function getyAxesConfig(options: any, unit: UnitEnum) {
    if (unit) {
        options.ticks.callback = function (value: any, index: any, values: Array<number>) {
            if (index === 0) {
                return "";
            }
            return DataFormatter.formatter(value, unit);
        };
    }

    return options;
}