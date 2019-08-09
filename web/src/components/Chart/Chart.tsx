import * as React from 'react'
import { reaction } from 'mobx'
import { observer } from 'mobx-react'
import { autobind } from 'core-decorators'
import StoreManager from '../../store/StoreManager'
import { ChartTooltipData, UnitEnum } from '../../model/Metric'

const ChartJS = require('chart.js')

interface ChartProps {
  data: any
  uuid: string
  type: string
  unit: UnitEnum
  options?: any
  plugins?: any
  onMouseMove?: (data: ChartTooltipData) => void
  onMouseOut?: () => void
}

interface ChartStatus {
}

@observer
export default class Chart extends React.Component<ChartProps, ChartStatus> {
  chartCanvas: React.RefObject<HTMLCanvasElement> // Chart Canvas Ref
  crosshair: React.RefObject<HTMLDivElement> // Chart Crosshair Ref
  chartInstance: any // ChartJS instance
  disposers: any[]

  constructor(props: ChartProps) {
    super(props)
    this.chartCanvas = React.createRef()
    this.crosshair = React.createRef()

    this.disposers = [
      reaction(
        () => StoreManager.ChartEventStore.hiddenSeries.get(props.uuid) || [],
        this.handleLegendItemClick,
      ),
      reaction(
        () => StoreManager.ChartEventStore.tooltipData,
        this.handleCrosshairMove,
      ),
    ]
  }

  componentDidMount(): void {
    this.renderChart()
    this.addEventListener()
  }

  componentWillUnmount(): void {
    this.disposers.map(handle => handle())
  }

  /**
   * init Chart
   */
  renderChart() {
    const { type, data, options, plugins } = this.props
    const canvas = this.chartCanvas.current
    if (!canvas) {
      return
    }

    const ctx = canvas.getContext('2d')

    this.chartInstance = new ChartJS(ctx, { type, data, options, plugins })
  }

  addEventListener() {
    const canvas = this.chartCanvas.current
    if (!canvas) {
      return
    }

    canvas.addEventListener('mousemove', this.handleMouseMove)
    canvas.addEventListener('mouseout', this.handleMouseOut)
  }

  @autobind
  handleMouseMove(e: MouseEvent) {
    const canvas = this.chartCanvas.current
    if (!canvas) {
      return
    }

    // Get all vertical points
    const points = this.chartInstance.getElementsAtXAxis(e) // vertical Points
    const index = points.length > 0 ? points[ 0 ]._index : 0  // current mouseover index
    const data = points.length > 0 ? points[ 0 ]._chart.data : []
    const { datasets } = data

    if (!datasets) {
      return
    }

    // Get all vertical series information
    const series = datasets
    .filter((_: any, idx: number) => {
      // filter series(hidden)
      const meta = this.chartInstance.getDatasetMeta(idx)
      return meta ? !meta.hidden : true
    })
    .map((item: any) => ({
      color: item.borderColor,
      name: item.label,
      value: item.data[ index ].y,
      time: +item.data[ index ].x,
    }))

    /**
     *  calculate size info
     */
    const { chartArea } = this.chartInstance
    const {
      left: canvasLeft,
      top: canvasTop,
      bottom: canvasBottom,
      width: canvasWidth,
    } = canvas.getBoundingClientRect()
    // Chart Border
    const {
      top: chartOffsetTopWithCanvas,
      left: chartOffsetLeftWithCanvas,
      right: chartOffsetRightWithCanvas,
      bottom: chartOffsetBottomWithCanvas,
    } = chartArea
    const chartWidth = chartOffsetRightWithCanvas - chartOffsetLeftWithCanvas
    // Mouse info
    const pointX = e.clientX

    const border = {
      pointX,
      canvasTop,
      canvasLeft,
      chartWidth,
      canvasWidth,
      canvasBottom,
      chartOffsetTopWithCanvas,
      chartOffsetLeftWithCanvas,
      chartOffsetRightWithCanvas,
      chartOffsetBottomWithCanvas,
    }

    const result: ChartTooltipData = {
      index,
      series,
      border,
      unit: this.props.unit,
      time: series.length > 0 ? series[ 0 ].time : 0,
    }

    const { onMouseMove } = this.props
    onMouseMove && onMouseMove(result)
  }

  @autobind
  handleMouseOut() {
    const { onMouseOut } = this.props
    onMouseOut && onMouseOut()
  }

  /**
   * Legend Click Handle
   * @param {any[]} hidden Current array of hidden series
   */
  @autobind
  handleLegendItemClick(hidden: any[]) {
    hidden.forEach((hide, index) => {
      const meta = this.chartInstance.getDatasetMeta(index)
      meta.hidden = hide
    })

    this.chartInstance.update({
      duration: 300,
    })
  }

  @autobind
  handleCrosshairMove(data: ChartTooltipData | null) {
    const crosshair = this.crosshair.current
    const canvas = this.chartCanvas.current
    const { chartArea } = this.chartInstance

    if (!crosshair || !canvas || !chartArea) {
      return
    }

    crosshair.style.display = data ? 'block' : 'none'
    crosshair.style.willChange = data ? 'transform' : 'auto'

    if (data) {
      const top = chartArea.top
      const height = chartArea.bottom - top

      const currChartWidth = chartArea.right - chartArea.left // Current chart width (valid area, not all Canvas)

      const originOffsetLeftByChart = Math.max(
        data.border.chartOffsetLeftWithCanvas,
        Math.min(data.border.pointX - data.border.canvasLeft, data.border.chartOffsetRightWithCanvas),
      ) - data.border.chartOffsetLeftWithCanvas // Starting by 0, relative to the offset of the valid area

      // Ratio of mouse position to valid area width
      const percentage = originOffsetLeftByChart / data.border.chartWidth

      // calculate crosshair position
      const left = currChartWidth === data.border.chartWidth
        ? originOffsetLeftByChart + data.border.chartOffsetLeftWithCanvas // when same width (higher precision)
        : percentage * currChartWidth + chartArea.left // when different width, percentage * valid area width + chart offsetLeft

      crosshair.style.height = `${height}px`

      crosshair.style.transform = `translate(${left}px, ${top}px)`

    }
  }

  render() {
    // Canvas Wrapped By a div element to avoid invoke `.resize` many times
    return (
      <div className="lindb-chart-canvas">
        <div className="lindb-chart-canvas__crosshair" ref={this.crosshair}/>
        <canvas ref={this.chartCanvas}/>
      </div>
    )
  }
}