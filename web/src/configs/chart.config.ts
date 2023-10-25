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
import { Chart, registerables } from "chart.js";
import * as helpers from "chart.js/helpers";
import * as _ from "lodash-es";
import { Unit } from "@src/models";
import { FormatKit } from "@src/utils";
import { Theme } from "@src/constants";

Chart.register(...registerables);

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

Chart.register({
  id: "message",
  beforeDraw: function (chart: any, _args: any, _options: any): boolean {
    const datasets = _.get(chart, "data.datasets", []);
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

export const DarkChart = {
  options: {
    scales: {
      x: {
        grid: {
          color: "rgba(255, 255, 255, 0.1)",
        },
        ticks: {
          color: "rgb(249, 249, 249)",
        },
      },
      y: {
        grid: {
          color: "rgba(255, 255, 255, 0.1)",
        },
        ticks: {
          color: "rgb(249, 249, 249)",
        },
      },
    },
  },
};

export const LightChart = {
  options: {
    scales: {
      x: {
        grid: {
          color: "rgba(232,233,234,1)",
        },
        ticks: {
          color: "rgba(28,31,35,0.8)",
        },
      },
      y: {
        grid: {
          // color: "rgba(28,31,35,0.2)",
          color: "rgba(232,233,234,1)",
        },
        ticks: {
          color: "rgba(28,31,35,0.8)",
        },
      },
    },
  },
};

export function getChartThemeConfig(theme: Theme, raw: any) {
  let chartTheme: any = LightChart;
  if (theme === Theme.dark) {
    chartTheme = DarkChart;
  }
  //IMPORTANT: need clone object, because merge return target object.
  return _.cloneDeep(_.merge(raw, chartTheme));
}

export const DefaultChartConfig = {
  type: "line",
  data: {},
  plugins: {
    message: {},
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
    zoom: true,
    legend: {
      asTable: true,
    },
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
        },
        ticks: {
          font: {
            size: 12,
          },
          //       // fontSize: 10,
          maxRotation: 0, // angle in degrees
          callback: function (_value: any, index: number, _values: any) {
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
        },
        ticks: {
          //       mirror: true, // draw tick in chart area
          //       display: true,
          //       // min: 0,
          font: { size: 12 },
          // autoSkip: true,
          callback: function (value: any, index: number, _values: any) {
            // if (index == 0) {
            //   //ignore first tick
            //   return null;
            // }
            if (index % 2 == 0) {
              return FormatKit.format(
                value,
                _.get(this, "chart.lin.extend.unit", Unit.Short)
              );
            }
            return null;
          },
        },
        // suggestedMin: 10,
        beginAtZero: true,
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
        radius: 0,
        hoverRadius: 0,
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
