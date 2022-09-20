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
import { Chart, ChartArea } from "chart.js/auto";
import { Series } from "@src/models";

export type MouseMoveEvent = {
  index: number;
  mouseX: number;
  chart: Chart;
  chartArea: ChartArea;
  chartCanvas: HTMLCanvasElement;
  chartCanvasRect: DOMRect;
  nativeEvent: MouseEvent;
};

export enum ChartStatus {
  Init = "init",
  Loading = "loading",
  OK = "ok",
  Empty = "empty",
  Error = "error",
}
export enum ChartType {
  Line = "line",
  Area = "area",
}

export type ChartConfig = {
  title?: string;
  description?: string;
  config?: any;
  timeShift?: string;
  targets?: Target[];
  unit?: UnitEnum;
};

export enum ChartTypeEnum {
  Line = "line",
}

export enum UnitEnum {
  None = "none",
  Short = "short",
  Bytes = "bytes",
  KBytesPerSecond = "KB/s",
  Percent = "percent",
  Percent2 = "percent(0-1)",
  Seconds = "seconds(s)",
  Milliseconds = "milliseconds(ms)",
  Nanoseconds = "nanoseconds(ns)",
}

export type Target = {
  db?: string;
  sql: string | QueryStatement;
  bind?: boolean;
  watch?: string[];
};

export type QueryStatement = {
  namespace?: string;
  metric?: string;
  field?: string[];
  tags?: Object;
  groupBy?: string[];
};

export interface ChartDataSource {
  datasets: SeriesData[];
  times?: number[];
  interval: any;
}

export interface SeriesData {
  data: number[];
  targetType: SeriesDataType;
  metric: Series;
  lineNum: string;
  name?: string;
  label: string;
  display: boolean;
  spanGaps: boolean;
  type: string;
  yAxisID: string;
  value: {
    total?: number;
    max?: number;
    min?: number;
    avg?: number;
    current?: number;
  };
}

export enum SeriesDataType {
  DEFAULT = "default",
  COMPUTED = "computed",
}
