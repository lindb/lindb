import { CANVAS_CHART_CONFIG, getyAxesConfig } from 'components/metric/ChartConfig';
import { autobind } from 'core-decorators';
import { cloneDeep } from 'lodash';
import { reaction } from 'mobx';
import { ChartStatusEnum } from 'model/Chart';
import { ChartTooltipData, UnitEnum } from 'model/Metric';
import * as React from 'react';
import StoreManager from 'store/StoreManager';
const ChartJS = require('chart.js')

interface ChartProps {
  height?: number
  uuid: string
  type: string
  unit: UnitEnum
  onMouseMove?: (data: ChartTooltipData) => void
  onMouseOut?: () => void
}

interface ChartState {
}

export default class Chart extends React.Component<ChartProps, ChartState> {
  chartCanvas: any // Chart Canvas Ref
  crosshair: React.RefObject<HTMLDivElement> // Chart Crosshair Ref
  chartInstance: any // ChartJS instance
  disposers: any[]
  chartConfig: any
  series: any = {}
  chartStatus: any

  constructor(props: ChartProps) {
    super(props)
    this.chartCanvas = React.createRef()
    this.crosshair = React.createRef()
    this.buildChartConfig()

    this.disposers = [
      reaction(
        () => StoreManager.ChartEventStore.hiddenSeries.get(props.uuid),
        () => {
          this.series = StoreManager.ChartStore.seriesCache.get(props.uuid) || {}
          if (!this.series) {
            return 
          }
          const selected = StoreManager.ChartStore.selectedSeries.get(props.uuid)
          this.series.datasets && this.series.datasets.forEach((item: any) => {
            item.hidden = selected && selected.size > 0 && !selected.has(item.label)
          })
          this.renderChart()
        }
      ),
      reaction(
        () => StoreManager.ChartEventStore.tooltipData,
        this.handleCrosshairMove,
      ),
      reaction(
        () => StoreManager.ChartStore.chartStatusMap.get(props.uuid),
        chartStatus => {
          if (chartStatus!.status !== ChartStatusEnum.Loaded) {
            return
          }
          this.series = StoreManager.ChartStore.seriesCache.get(props.uuid) || {};
          this.renderChart()
          this.addEventListener()
        }
      )
    ]
  }

  componentWillUnmount(): void {
    this.disposers.map(handle => handle())
  }

  shouldComponentUpdate(nextProps: Readonly<ChartProps>, nextState: Readonly<ChartState>, nextContext: any): boolean {
    return false;
  }

  buildChartConfig() {
    if (this.chartConfig) {
      return;
    }
    const { type } = this.props
    const chartConfig: any = cloneDeep(CANVAS_CHART_CONFIG);
    chartConfig.type = type;
    this.chartConfig = chartConfig
  }

  @autobind
  setData(chart: any) {
    chart.options.interval = this.series.interval;
    chart.options.times = this.series.times;
    chart.data = this.series;
    chart.options.scales.yAxes[0].ticks.suggestedMax = this.series.leftMax * 1.05;
    chart.options.scales.yAxes[0] = getyAxesConfig(chart.options.scales.yAxes[0], this.props.unit);
  }

  @autobind
  getShowSeries(reportData: any): any {
    const datasets: Array<any> = reportData.datasets;
    let series: any = [];
    datasets.forEach(dataset => {
      let display = dataset.seriesDisplay;
      if (display === undefined || display === true) {
        series.push(dataset);
      }
    });
    return series;
  }
  /**
   * init Chart
   */
  renderChart() {
    this.setData(this.chartConfig)
    console.log(this.chartConfig)
    if (this.chartInstance) {
      this.chartInstance.update()
    } else {
      const canvas = this.chartCanvas.current
      const ctx = canvas.getContext('2d')
      this.chartInstance = new ChartJS(ctx, this.chartConfig)
    }
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
    const index = points.length > 0 ? points[0]._index : 0  // current mouseover index
    const data = points.length > 0 ? points[0]._chart.data : []
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
        value: item.data[index],
        time: this.series.times[index],
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
      time: series.length > 0 ? series[0].time : 0,
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
  handleLegendItemClick() {
    // hidden.forEach((hide, index) => {
    //   const meta = this.chartInstance.getDatasetMeta(index)
    //   meta.hidden = hide
    // })
    console.log("xxxxxx", this.chartInstance.get)

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
    const { height } = this.props
    let h = height
    if (height! <= 0) {
      h = 280
    }
    console.log("xxxx", h, height)
    // Canvas Wrapped By a div element to avoid invoke `.resize` many times
    return (
      <div className="lindb-chart-canvas" style={{ height: h }}>
        <div className="lindb-chart-canvas__crosshair" ref={this.crosshair} />
        <canvas ref={this.chartCanvas} height="280px" />
      </div>
    )
  }
}