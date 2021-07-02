import {Chart} from 'model/Chart'

export type Board = ChartRow[]
export type ChartRow = Metric[]

export interface Metric {
  span?: number
  id: string
  chart: Chart
}

export enum UnitEnum {
  None = 'none',
  Bytes = 'bytes',
  KBytesPerSecond = 'KB/s',
  Percent = 'percent',
  Seconds = 'seconds(s)',
  Milliseconds = 'milliseconds(ms)',
  Nanoseconds = "nanoseconds(ns)",
}

/* Result Set */
export class ResultSet {
  metricName?: string
  series?: Series[]
  startTime?: number
  endTime?: number
  interval?: number
  stats?: any
}

export class Series {
  tags?: { [propName: string]: string }
  fields?: { [propName: string]: { [timestamp: string]: number } }
}

export interface ChartDatasets {
  fill: boolean
  label: string
  data: any
  backgroundColor: string
  borderColor: string
  pointBackgroundColor: string
  hidden: boolean
}

/* Tooltip Data Source */
export interface ChartTooltipData {
  time: number
  index: number
  unit: UnitEnum
  border: ChartBorderInfo
  series: ChartTooltipItem[]
}

export interface ChartBorderInfo {
  pointX: number
  canvasTop: number
  canvasLeft: number
  chartWidth: number
  canvasWidth: number
  canvasBottom: number
  chartOffsetTopWithCanvas: number
  chartOffsetLeftWithCanvas: number
  chartOffsetRightWithCanvas: number
  chartOffsetBottomWithCanvas: number
}

export interface ChartTooltipItem {
  time: number
  value: number
  name: string
  color: string
}