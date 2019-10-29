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
  metricName?: string
  series?: Series[]
  startTime?: number
  endTime?: number
  interval?: number
}

export class Series {
  tags?: { [propName: string]: string }
  fields?: { [propName: string]: { [timestamp: string]: number } }
}

export interface ChartDatasets {
  fill: boolean
  label: string
  data: ChartDatasetsData[]
  backgroundColor: string
  borderColor: string
  pointBackgroundColor: string
}

export interface ChartDatasetsData {
  x: Date,
  y: number
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