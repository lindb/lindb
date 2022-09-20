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
import ChartLegend from "@src/components/chart/ChartLegend";
import { DefaultChartConfig, getChartThemeConfig } from "@src/configs";
import { ChartStatus, ChartTypeEnum, UnitEnum } from "@src/models";
import { ChartEventStore, ChartStore } from "@src/stores";
import urlStore from "@src/stores/url.store";
import { setStyle } from "@src/utils";
import Chart from "chart.js/auto";
import * as _ from "lodash-es";
import { reaction } from "mobx";
import moment from "moment";
import React, { MutableRefObject, useContext, useEffect, useRef } from "react";
import { DateTimeFormat } from "@src/constants";
import { UIContext } from "@src/context/UIContextProvider";

interface CanvasChartProps {
  chartId: string;
  height?: number;
  disableDrag?: boolean;
}
const Zoom = {
  drag: false,
  isMouseDown: false,
  selectedStart: 0,
  selectedEnd: 0,
};

export default function CanvasChart(props: CanvasChartProps) {
  const { chartId, height, disableDrag } = props;
  const eventCallbacks: Map<string, any> = new Map();
  const { theme } = useContext(UIContext);
  const chartRef = useRef() as MutableRefObject<HTMLCanvasElement | null>;
  const crosshairRef = useRef() as MutableRefObject<HTMLDivElement>;
  const chartObjRef = useRef() as MutableRefObject<Chart | null>;
  const zoomRef = useRef(_.cloneDeep(Zoom));
  const zoomDivRef = useRef() as MutableRefObject<HTMLDivElement>;
  const seriesRef = useRef() as MutableRefObject<any>;
  const chartStatusRef = useRef() as MutableRefObject<ChartStatus | undefined>;

  useEffect(() => {
    // set chart theme if chart instance exist
    for (let key of Object.keys(Chart.instances)) {
      const currChart = Chart.instances[`${key}`];
      currChart.config = getChartThemeConfig(
        theme,
        _.get(currChart, "config", {})
      );
      currChart.update();
    }
  }, [theme]);

  const resetZoomRange = () => {
    if (!disableDrag) {
      setStyle(zoomDivRef.current, {
        display: "none",
        width: "0px",
      });
      zoomRef.current.isMouseDown = false;
      zoomRef.current.drag = false;
      zoomRef.current.selectedStart = 0;
      zoomRef.current.selectedEnd = 0;
    }
  };

  const createChart = () => {
    // console.log("config", { type: "line", data: series }, config);
    const canvas = chartRef.current;
    if (!canvas) {
      return;
    }
    const chartCfg = ChartStore.charts.get(chartId);
    const config: any = _.merge(
      {
        type: _.get(chartCfg, "type", ChartTypeEnum.Line),
        unit: _.get(chartCfg, "unit", UnitEnum.None),
      },
      getChartThemeConfig(theme, DefaultChartConfig)
    );
    config.options.crosshair = crosshairRef.current;

    const chartInstance = new Chart(canvas, config);
    chartObjRef.current = chartInstance;
    let start = 0;

    eventCallbacks.set("mousemove", function (e: MouseEvent) {
      if (chartStatusRef.current != ChartStatus.OK) {
        return;
      }
      const points: any = chartInstance.getElementsAtEventForMode(
        e,
        "index",
        { intersect: false },
        false
      );
      if (!points || points.length <= 0) {
        return;
      }
      const chartArea = chartInstance.chartArea;
      const currIdx = points[0].index;
      const x = e.offsetX;
      if (!disableDrag && zoomRef.current.isMouseDown) {
        zoomRef.current.selectedEnd = seriesRef.current.times[points[0].index];
        zoomRef.current.drag = true;
        const width = e.offsetX - start;
        if (width >= 0) {
          setStyle(zoomDivRef.current, {
            width: `${width}px`,
          });
        } else {
          setStyle(zoomDivRef.current, {
            width: `${-width}px`,
            transform: `translate(${e.offsetX}px, ${chartArea.top}px)`,
          });
        }
      }

      const canvaxRect = canvas.getBoundingClientRect();

      const interval =
        _.get(chartInstance, "config._config.data.interval", 0) / 1000;
      const v = currIdx * interval;
      // cross hair
      for (let key of Object.keys(Chart.instances)) {
        const currChart = Chart.instances[`${key}`];
        const crosshair = _.get(currChart, "options.crosshair", null);
        const len = _.get(currChart, "config._config.data.times", []).length;
        if (!crosshair || len == 0) {
          continue;
        }

        const chartArea = currChart.chartArea;
        const width = _.get(currChart, "chartArea.width", 0);
        const i = _.get(currChart, "config._config.data.interval", 0) / 1000;

        const x = (v / ((len - 1) * i)) * width + chartArea.left;
        if (x > chartArea.right) {
          continue;
        }
        const top = chartArea.top;
        const bottom = chartArea.bottom;
        setStyle(crosshair, {
          display: "block",
          height: `${bottom - top}px`,
          transform: `translate(${x}px, ${top}px)`,
        });
      }

      ChartEventStore.setShowTooltip(true);
      ChartEventStore.mouseMove({
        index: currIdx,
        mouseX: x,
        chart: chartInstance,
        chartArea: chartArea,
        chartCanvas: canvas,
        chartCanvasRect: canvaxRect,
        nativeEvent: e,
      });
    });
    eventCallbacks.set("mouseleave", function (e: any) {
      if (chartStatusRef.current != ChartStatus.OK) {
        return;
      }

      // hide cross hair
      for (let key of Object.keys(Chart.instances)) {
        const currChart = Chart.instances[`${key}`];
        const crosshair = _.get(currChart, "options.crosshair", null);
        setStyle(crosshair, {
          display: "none",
        });
      }

      // reset zoom range selection if leave chart area
      resetZoomRange();

      ChartEventStore.mouseLeave(e);
      ChartEventStore.setShowTooltip(false);
    });
    if (!disableDrag) {
      eventCallbacks.set("mousedown", function (e: any) {
        if (chartStatusRef.current != ChartStatus.OK) {
          return;
        }
        zoomRef.current.isMouseDown = true;
        const points: any = chartInstance.getElementsAtEventForMode(
          e,
          "index",
          { intersect: false },
          false
        );
        if (points && points.length > 0) {
          zoomRef.current.selectedStart =
            seriesRef.current.times[points[0].index];
        }
        start = e.offsetX;
        const chartArea = chartInstance.chartArea;
        const height = chartArea.height;
        setStyle(zoomDivRef.current, {
          display: "block",
          height: `${height}px`,
          // left: start,
          transform: `translate(${start}px, ${chartArea.top}px)`,
        });
      });
      eventCallbacks.set("mouseup", function (_e: any) {
        if (chartStatusRef.current != ChartStatus.OK) {
          return;
        }
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
          urlStore.changeURLParams({ params: { from: from, to: to } });
        }
        resetZoomRange();
      });
    }
    eventCallbacks.forEach((v, k) => {
      canvas.addEventListener(k, v);
    });
  };

  /**
   * set chart display data
   * @param series data which display in chart
   */
  const setChartData = (series: any) => {
    const chartInstance = chartObjRef.current;
    if (chartInstance) {
      chartInstance.data = series;
    }
  };

  /**
   * Init canvas chart component.
   * 1. watch chart status from ChartStore.
   */
  useEffect(() => {
    const disposer = [
      reaction(
        () => ChartStore.chartStatusMap.get(chartId),
        (s: ChartStatus | undefined) => {
          chartStatusRef.current = s;
          if (!s || s == ChartStatus.Loading) {
            return;
          }
          const series = ChartStore.seriesCache.get(chartId);
          seriesRef.current = series;
          const chartInstance = chartObjRef.current;

          if (chartInstance) {
            setChartData(series);
            chartInstance.update();
          } else {
            createChart();
            setChartData(series);
          }
        }
      ),
    ];
    const canvas = chartRef.current;

    return () => {
      if (canvas) {
        eventCallbacks.forEach((v, k) => {
          canvas.removeEventListener(k, v);
        });
      }
      disposer.forEach((d) => d());
      if (chartObjRef.current) {
        chartObjRef.current.destroy();
        // reset chart obj as null, maybe after hot load(develop) canvas element not ready.
        // invoke chart update will fail.
        chartObjRef.current = null;
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [chartId]);

  return (
    <div>
      <div className="lin-chart" style={{ height: height || 200 }}>
        <canvas className="chart" ref={chartRef} height={height || 200} />
        <div ref={crosshairRef} className="crosshair" />
        {!disableDrag && <div ref={zoomDivRef} className="zoom" />}
      </div>
      <ChartLegend />
    </div>
  );
}
