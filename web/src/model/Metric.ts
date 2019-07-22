export type Board = ChartRow[]
export type ChartRow = Metric[]

export interface Metric {
  span?: number
  chart: Chart
}

export interface Chart {
  id?: number
  ql?: string
  title?: string
  group?: object
  series?: Array<any>
  unit?: UnitEnum
}

export enum UnitEnum {
  None = 'none',
  Bytes = 'bytes',
  Percent = 'percent',
  Seconds = 'seconds(s)',
  Milliseconds = 'milliseconds(ms)',
}

/* Result Set */
export class ResultSet {
  queryType?: string
  result?: Result
}

export class Result {
  data?: any
  endTime?: number
  interval?: number
  startTime?: number
  pointCount?: number
  measurementName?: string
  groups?: Group[]
}

export class Group {
  group?: Map<string, string>
  fields?: Map<string, Array<number>>
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