/*
Licensed to LinDB under one or more contributor
license agreements. See the NOTICE file distributed with
this work for additional information regarding copyright
ownership. LinDB licenses this file to you under
the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
import Chart from "chart.js/auto";
import * as helpers from "chart.js/helpers";
import * as _ from "lodash-es";
import { UnitEnum } from "@src/models";
import { formatter } from "@src/utils";

export const INTERVALS = [
  10000,
  30000,
  60 * 1000,
  2 * 60 * 1000,
  5 * 60 * 1000,
  10 * 60 * 1000,
  15 * 60 * 1000,
  30 * 60 * 1000,
];
function getVisibleInterval(times: any, interval: number, width: number) {
  const ticksCount = times.length;
  const perTickPX = width / (ticksCount - 1);
  let ticksPerVisibleTick = parseInt(`${100 / perTickPX}`);
  const t = ticksPerVisibleTick * interval;
  let result = interval;
  if (t < 60 * 60 * 1000) {
    INTERVALS.forEach((item) => {
      if (t / item < 0.6) {
        return;
      } else {
        result = item;
      }
    });
  } else {
    result = t - (t % (60 * 60 * 1000));
  }
  return result;
}

// Chart.register({
//   id: "lazy",
//   afterUpdate: function (chart: any) {
//     const xAxes = chart.scales["xAxes-bottom"];
//     const tickOffset = _.get(xAxes, "options.ticks.tickOffset", null);
//     const display = _.get(xAxes, "options.display", false);
//     if (display && tickOffset) {
//       const width = _.get(chart, "scales.xAxes-bottom.width", 0);
//       const interval = _.get(chart, "config.options.interval", 0);
//       const times = _.get(chart, "config.options.times", []);
//       const visibleInterval = getVisibleInterval(times, interval, width);

//       xAxes.draw = function () {
//         const xScale = chart.scales["xAxes-bottom"];

//         // const tickFontColor = helpers.getValueOrDefault(
//         //   xScale.options.ticks.fontColor,
//         //   Chart.defaults.defaultFontColor
//         // );
//         // const tickFontSize = helpers.getValueOrDefault(
//         //   xScale.options.ticks.fontSize,
//         //   ChartJS.defaults.global.defaultFontSize
//         // );
//         // const tickFontStyle = helpers.getValueOrDefault(
//         //   xScale.options.ticks.fontStyle,
//         //   ChartJS.defaults.global.defaultFontStyle
//         // );
//         // const tickFontFamily = helpers.getValueOrDefault(
//         //   xScale.options.ticks.fontFamily,
//         //   ChartJS.defaults.global.defaultFontFamily
//         // );
//         // const tickLabelFont = helpers.fontString(
//         //   tickFontSize,
//         //   tickFontStyle,
//         //   tickFontFamily
//         // );
//         const tl = xScale.options.gridLines.tickMarkLength;

//         const isRotated = xScale.labelRotation !== 0;
//         const yTickStart = xScale.top;
//         const yTickEnd = xScale.top + tl;
//         const chartArea = chart.chartArea;

//         helpers.each(
//           xScale.ticks,
//           (label: any, index: any) => {
//             if (times[index] % visibleInterval !== 0) {
//               return;
//             }
//             // console.log("xxxxxx",index,times,visibleInterval)

//             // copy of chart.js code
//             let xLineValue = this.getPixelForTick(index);
//             const xLabelValue = this.getPixelForTick(
//               index,
//               this.options.gridLines.offsetGridLines
//             );

//             if (this.options.gridLines.display) {
//               this.ctx.lineWidth = this.options.gridLines.lineWidth;
//               this.ctx.strokeStyle = this.options.gridLines.color;

//               xLineValue += helpers.aliasPixel(this.ctx.lineWidth);

//               // Draw the label area
//               this.ctx.beginPath();

//               if (this.options.gridLines.drawTicks) {
//                 this.ctx.moveTo(xLineValue, yTickStart);
//                 this.ctx.lineTo(xLineValue, yTickEnd);
//               }

//               // Draw the chart area
//               if (this.options.gridLines.drawOnChartArea) {
//                 this.ctx.moveTo(xLineValue, chartArea.top);
//                 this.ctx.lineTo(xLineValue, chartArea.bottom);
//               }

//               // Need to stroke in the loop because we are potentially changing line widths & colours
//               this.ctx.stroke();
//             }

//             if (this.options.ticks.display) {
//               this.ctx.save();
//               this.ctx.translate(
//                 xLabelValue + this.options.ticks.labelOffset,
//                 isRotated
//                   ? this.top + 12
//                   : this.options.position === "top"
//                   ? this.bottom - tl
//                   : this.top + tl
//               );
//               this.ctx.rotate(helpers.toRadians(this.labelRotation) * -1);
//               // this.ctx.font = tickLabelFont;
//               this.ctx.textAlign = isRotated ? "right" : "center";
//               this.ctx.textBaseline = isRotated
//                 ? "middle"
//                 : this.options.position === "top"
//                 ? "bottom"
//                 : "top";
//               // this.ctx.fillStyle = tickFontColor;
//               this.ctx.fillText(label, 0, 0);
//               this.ctx.restore();
//             }
//           },
//           xScale
//         );
//       };
//     }
//   },
//   afterDraw: function (chart: any) {
//     // const status = _.get(chart, "options.status", null);
//     // const ctx = chart.chart.ctx;
//     // if (isNoData(status)) {
//     //   chart.clear();
//     //   const width = chart.chart.width;
//     //   const height = chart.chart.height;
//     //   let text = "";
//     //   let color = get(chart, "options.scales.yAxes[0].ticks.fontColor", null);
//     //   ctx.textAlign = "center";
//     //   ctx.textBaseline = "middle";
//     //   switch (status.status) {
//     //     case ChartStatusEnum.NoData:
//     //       text = "No data to display";
//     //       break;
//     //     case ChartStatusEnum.BadRequest:
//     //       text = "Invalid Configuration";
//     //       color = "#D27613";
//     //       break;
//     //     case ChartStatusEnum.LoadError:
//     //       color = "#F56C6C";
//     //       ctx.fillStyle = color;
//     //       ctx.fillText(status.msg, width / 2, height / 2 + 20);
//     //       text = "Internal Server Error";
//     //       break;
//     //     default:
//     //       break;
//     //   }
//     //   ctx.font = "13px Arial";
//     //   ctx.fillStyle = color;
//     //   ctx.fillText(text, width / 2, height / 2);
//     //   ctx.restore();
//     // } else if (chart.options.isSeriesChart) {
//     //   const chartArea = chart.chartArea;
//     //   ctx.beginPath();
//     //   ctx.lineWidth = 1;
//     //   ctx.strokeStyle = _.get(
//     //     chart,
//     //     "options.scales.yAxes[0].gridLines.color",
//     //     null
//     //   );
//     //   ctx.moveTo(chartArea.left, chartArea.bottom);
//     //   ctx.lineTo(chartArea.right, chartArea.bottom);
//     //   ctx.stroke();
//     // }
//   },
// });
Chart.register({
  id: "message",
  beforeDraw: function(chart: any, _args: any, _options: any): boolean {
    const datasets = _.get(chart, "config._config.data.datasets", []);
    if (datasets.length <= 0) {
      // display no data message
      chart.clear();
      const ctx = chart.ctx;
      const width = chart.width;
      const height = chart.height;
      ctx.font = "14px Arial";
      ctx.fillStyle = "rgba(249, 249, 249, 0.6)";
      ctx.textAlign = "center";
      ctx.fillText("No data in response", width / 2, height / 2);
      ctx.restore();
      return false;
    }
    return true;
  },
} as any);

export const DefaultChartConfig = {
  type: undefined,
  data: {},
  plugins: {
    message: {},
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
    scales: {
      x: {
        type: "category",
        grid: {
          //       drawTicks: false,
          lineWidth: 0.3,
          //       // tickMarkLength: 2,
          tickLength: 0,
          //       // drawOnChartArea: false,
          //       drawBorder: false,
          color: "rgba(255, 255, 255, 0.1)",
        },
        ticks: {
          font: {
            size: 12,
          },
          //       // fontSize: 10,
          maxRotation: 0, // angle in degrees
          color: "rgb(249, 249, 249)",
          callback: function(_value: any, index: number, _values: any) {
            const times = _.get(this, "chart.config._config.data.times", []);
            const labels = _.get(this, "chart.config._config.data.labels", []);
            if (times[index] % (5 * 60 * 1000) == 0) {
              return labels[index];
            }
            return null;
          },
          //       // align: "end", // default: center, start/end
          //       // tickOffset: 100,
          //       // fontColor: undefined,
        },
        //     display: undefined,
        //     stacked: undefined,
      },
      y: {
        grid: {
          //       drawTicks: false,
          //       lineWidth: 0.3,
          //       autoSkip: true,
          tickLength: 0,
          //       // zeroLineWidth: 0,
          //       drawBorder: false,
          //       // drawOnChartArea: false,
          //       // borderDash: [1, 1],
          color: "rgba(255, 255, 255, 0.1)",
        },
        ticks: {
          //       mirror: true, // draw tick in chart area
          //       display: true,
          //       // min: 0,
          font: { size: 12 },
          color: "rgb(249, 249, 249)",
          // autoSkip: true,
          callback: function(value: any, index: number, _values: any) {
            // if (index == 0) {
            //   //ignore first tick
            //   return null;
            // }
            if (index % 2 == 0) {
              return formatter(
                value,
                _.get(this, "chart.config._config.unit", UnitEnum.None)
              );
            }
            return null;
          },
          //       // tickMarkLength: 0,
          //       // maxTicksLimit: 6,
          suggestedMin: 0,
        },
        // suggestedMax: 10,
      },
    },
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        mode: "dataset",
        enabled: false,
      },
      title: {
        display: false,
      },
    },
    elements: {
      line: {
        tension: 0, // disables bezier curve
        borderWidth: 1,
        fill: undefined,
      },
      point: {
        radius: 1,
        hoverRadius: 2,
        pointStyle: undefined,
      },
      arc: {
        borderWidth: 0,
      },
    },
    hover: {
      // animationDuration: 0, // duration of animations when hovering an item
      mode: "index",
      intersect: false,
      // onHover: undefined,
    },
  },
};
