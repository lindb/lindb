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
import React, {
  MutableRefObject,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { Chart } from "chart.js";
import { UIContext } from "@src/context/UIContextProvider";
import { DefaultChartConfig, getChartThemeConfig } from "@src/configs";
import classNames from "classnames";
import Legend from "./Legend";
import Tooltip from "./Tooltip";
import * as _ from "lodash-es";
import { ChartKit, CSSKit } from "@src/utils";
import { PlatformStore, URLStore } from "@src/stores";
import { MouseEventType } from "@src/stores/platform.store";
import { ChartType, Unit } from "@src/models";
import moment from "moment";
import { DateTimeFormat } from "@src/constants";

const Zoom = {
  drag: false,
  isMouseDown: false,
  selectedStart: 0,
  selectedEnd: 0,
  x: 0,
};

const CanvasChart: React.FC<{
  type: ChartType;
  unit?: Unit;
  config?: any;
  datasets?: any;
}> = (props) => {
  const { type, unit, config, datasets } = props;
  const [selectedSeries, setSelectedSeries] = useState<string[]>([]);
  const { theme } = useContext(UIContext);

  const canvasRef = useRef() as MutableRefObject<HTMLCanvasElement | null>;
  const chartInstance = useRef() as MutableRefObject<Chart | null>;
  const zoomRef = useRef(_.cloneDeep(Zoom));
  const zoomDivRef = useRef() as MutableRefObject<HTMLDivElement>;
  const crosshairRef = useRef() as MutableRefObject<HTMLDivElement>;
  const hiddenSeriesMap = new Map<string, string>();

  const isZoom = (chart: Chart | null) => {
    return _.get(chart, "options.zoom", false);
  };

  const resetZoomRange = () => {
    if (isZoom(chartInstance.current)) {
      CSSKit.setStyle(zoomDivRef.current, {
        display: "none",
        width: "0px",
      });
      zoomRef.current.isMouseDown = false;
      zoomRef.current.drag = false;
      zoomRef.current.selectedStart = 0;
      zoomRef.current.selectedEnd = 0;
      zoomRef.current.x = 0;
    }
  };

  const onSeriesSelect = useCallback(
    (series: any, e: React.MouseEvent) => {
      const chart = chartInstance.current;
      if (!chart) {
        return;
      }
      // document.getSelection().removeAllRanges();

      const label = series.label;
      hiddenSeriesMap.set(label, label);
      const datasets = _.get(chart, "data.datasets", []);
      const currentSelectedSet = new Set<string>(selectedSeries);
      if (e.shiftKey) {
        currentSelectedSet.add(label);
      } else if (e.metaKey) {
        let hidden = false;
        datasets.forEach((series: any) => {
          if (!series.hidden) {
            currentSelectedSet.add(series.label);
          } else if (label === series.label) {
            hidden = true;
          }
        });
        if (hidden) {
          currentSelectedSet.add(label);
        } else {
          currentSelectedSet.delete(series.label);
        }
      } else {
        currentSelectedSet.clear();
        currentSelectedSet.add(label);
      }

      setSelectedSeries(Array.from(currentSelectedSet));
      datasets.forEach((series: any) => {
        series.hidden = !currentSelectedSet.has(series.label);
      });
      chart.update();
    },
    [selectedSeries]
  );

  const handleMouseDown = (e: MouseEvent) => {
    if (!chartInstance.current) {
      return;
    }
    const chart = chartInstance.current;
    const datasets = _.get(chart, "data", {}) as any;
    if (!_.isEmpty(datasets) && isZoom(chart)) {
      zoomRef.current.isMouseDown = true;
      const points: any = chart.getElementsAtEventForMode(
        e,
        "index",
        { intersect: false },
        false
      );
      if (points && points.length > 0) {
        zoomRef.current.selectedStart = datasets.times[points[0].index];
      }
      zoomRef.current.x = e.offsetX;
      const chartArea = chart.chartArea;
      const height = chartArea.height;
      CSSKit.setStyle(zoomDivRef.current, {
        display: "block",
        height: `${height}px`,
        transform: `translate(${zoomRef.current.x}px, ${chartArea.top}px)`,
      });
    }
  };

  const handleMouseUp = (_e: MouseEvent) => {
    if (zoomRef.current.drag) {
      const start = Math.min(
        zoomRef.current.selectedStart,
        zoomRef.current.selectedEnd
      );
      const end = Math.max(
        zoomRef.current.selectedStart,
        zoomRef.current.selectedEnd
      );
      const from = moment(start).format(DateTimeFormat);
      const to = moment(end).format(DateTimeFormat);
      URLStore.changeURLParams({ params: { from: from, to: to } });
    }
    resetZoomRange();
  };

  const handleMouseMove = (e: MouseEvent) => {
    if (!chartInstance.current) {
      return;
    }
    const chart = chartInstance.current;
    const points: any = chart.getElementsAtEventForMode(
      e,
      "index",
      { intersect: false },
      false
    );
    if (!points || points.length <= 0) {
      return;
    }
    const currIdx = points[0].index;
    const interval = _.get(chart, "config.data.interval", 0) / 1000;
    const v = currIdx * interval;

    if (zoomRef.current.isMouseDown && isZoom(chart)) {
      const chartArea = chart.chartArea;
      const datasets = _.get(chart, "data", {}) as any;
      zoomRef.current.selectedEnd = datasets.times[points[0].index];
      zoomRef.current.drag = true;
      const width = e.offsetX - zoomRef.current.x;
      if (width >= 0) {
        CSSKit.setStyle(zoomDivRef.current, {
          width: `${width}px`,
        });
      } else {
        CSSKit.setStyle(zoomDivRef.current, {
          width: `${-width}px`,
          transform: `translate(${e.offsetX}px, ${chartArea.top}px)`,
        });
      }
    }

    // move cross hair
    for (let key of Object.keys(Chart.instances)) {
      const currChart = Chart.instances[`${key}`];
      const crosshair = _.get(currChart, "lin.extend.crosshair", null);
      const len = _.get(currChart, "config.data.times", []).length;
      if (!crosshair || len == 0) {
        continue;
      }

      const chartArea = currChart.chartArea;
      const width = _.get(currChart, "chartArea.width", 0);
      const i = _.get(currChart, "config.data.interval", 0) / 1000;

      const x = (v / ((len - 1) * i)) * width + chartArea.left;
      if (x > chartArea.right) {
        continue;
      }
      const top = chartArea.top;
      const bottom = chartArea.bottom;
      CSSKit.setStyle(crosshair, {
        display: "block",
        height: `${bottom - top}px`,
        transform: `translate(${x}px, ${top}px)`,
      });
    }
    // set platform state context
    PlatformStore.setChartMouseEvent({
      type: MouseEventType.Move,
      index: currIdx,
      native: e,
      chart: chartInstance.current,
    });
  };

  const handleMouseOut = (e: MouseEvent) => {
    try {
      resetZoomRange();
      for (let key of Object.keys(Chart.instances)) {
        const currChart = Chart.instances[`${key}`];
        const crosshair = _.get(currChart, "lin.extend.crosshair", null);
        if (!crosshair) {
          continue;
        }

        CSSKit.setStyle(crosshair, {
          display: "none",
        });
      }
    } finally {
      PlatformStore.setChartMouseEvent({
        type: MouseEventType.Out,
        native: e,
      });
    }
  };

  useEffect(() => {
    // if canvas is null, return it.
    if (!canvasRef.current) {
      return;
    }
    const chartCfg: any = _.merge(
      getChartThemeConfig(theme, DefaultChartConfig),
      config
    );
    if (!chartInstance.current) {
      // component init
      chartCfg.type = ChartKit.getChartType(type);
      chartCfg.data = datasets || [];

      const canvas = canvasRef.current;
      const chart = new Chart(canvas, chartCfg);
      _.set(chart, "lin.extend.crosshair", crosshairRef.current);
      _.set(chart, "lin.extend.unit", unit || Unit.Short);
      _.set(chart, "lin.extend.onSeriesSelect", onSeriesSelect);
      _.set(chart, "lin.extend.theme", theme);
      chartInstance.current = chart;
      // add mouse event handles, after chart created
      canvas.addEventListener("mousemove", handleMouseMove);
      canvas.addEventListener("mouseout", handleMouseOut);
      canvas.addEventListener("mousedown", handleMouseDown);
      canvas.addEventListener("mouseup", handleMouseUp);
    } else {
      const chart = chartInstance.current;
      chart.data = datasets || [];
      chart.options = chartCfg.options;
      _.set(chart, "lin.extend.unit", unit || Unit.Short);
      if (theme !== _.get(chart, "lin.extend.theme")) {
        // theme changed
        _.set(chart, "lin.extend.theme", theme);
        chartInstance.current.config.options = getChartThemeConfig(theme, {
          options: _.get(chart, "config.options", {}),
        }).options;
      }
      // data set or config update
      chart.update();
    }

    const chartDatasets = _.get(chartInstance.current, "data.datasets", []);
    const currentSelectedSet = new Set<string>(selectedSeries);
    chartDatasets.forEach((series: any) => {
      if (!series.hidden) {
        currentSelectedSet.add(series.label);
      }
    });
    setSelectedSeries(Array.from(currentSelectedSet));
  }, [config, datasets, theme]);

  /**
   * destroy resource
   */
  useEffect(() => {
    return () => {
      if (chartInstance.current) {
        const canvas = chartInstance.current.canvas;
        canvas.removeEventListener("mousemove", handleMouseMove);
        canvas.removeEventListener("mouseout", handleMouseOut);
        canvas.removeEventListener("mousedown", handleMouseDown);
        canvas.removeEventListener("mouseup", handleMouseUp);
        chartInstance.current.destroy();
        chartInstance.current = null;
      }
    };
  }, []);

  const chartMetricCls = classNames("chart-metric-container", {
    "chart-cursor-pointer": false,
  });
  return (
    <div className={chartMetricCls}>
      <div className="chart-canvas-wrapper">
        <canvas ref={canvasRef} />
        <div ref={crosshairRef} className="crosshair" />
        <div ref={zoomDivRef} className="zoom" />
      </div>
      <Legend chart={chartInstance.current} />
      <Tooltip chart={chartInstance.current} />
    </div>
  );
};

export default CanvasChart;
