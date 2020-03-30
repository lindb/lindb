import * as React from 'react'
import { observer } from 'mobx-react'
import { autobind } from 'core-decorators'
import { observable, reaction } from 'mobx'
import StoreManager from '../../store/StoreManager'
import { DataFormatter } from '../../utils/DataFormatter'
import { ChartBorderInfo, ChartTooltipData } from '../../model/Metric'
import Moment from 'react-moment';

interface ChartTooltipProps {
}

interface ChartTooltipStatus {
}

@observer
export default class ChartTooltip extends React.Component<ChartTooltipProps, ChartTooltipStatus> {
  @observable dataSource: ChartTooltipData | null = null // Tooltip data source, do not display when value is null

  tooltip: React.RefObject<HTMLDivElement>
  tooltipArrow: React.RefObject<HTMLDivElement>
  disposers: any[]
  private timer = 0 // tooltip countdown timer
  private INTERVAL = 350 // Interval for disappear(ms)
  private SCREEN_SPACING = 5 // Minimum gap between Tooltip and windows

  constructor(props: ChartTooltipProps) {
    super(props)

    this.tooltip = React.createRef()
    this.tooltipArrow = React.createRef()

    this.disposers = [
      reaction(() => StoreManager.ChartEventStore.tooltipData, this.handleTooltipDataChange),
    ]
  }

  /**
   * Get size information for container
   */
  componentDidUpdate(prevProps: Readonly<ChartTooltipProps>, prevState: Readonly<ChartTooltipStatus>): void {
    const tooltip = this.tooltip.current
    const tooltipArrow = this.tooltipArrow.current
    const data = this.dataSource

    if (tooltip && tooltipArrow && data) {
      const { left: containerLeft, top: containerTop, isBottom } = this.calculateContainerPosition(data.border)
      const arrowLeft = this.calculateArrowPosition(data.border)

      // Set Container position
      isBottom ? tooltip.classList.add('is-bottom') : tooltip.classList.remove('is-bottom')
      tooltip.style.left = `${containerLeft}px`
      tooltip.style.top = `${containerTop}px`

      // Set Arrow position
      tooltipArrow.style.left = `${arrowLeft}px`
    }
  }

  componentWillUnmount(): void {
    this.disposers.map(handle => handle())
  }

  /**
   * Set countdown to disappear
   */
  @autobind
  setDisappearTimer() {
    this.timer = +setTimeout(() => this.dataSource = null, this.INTERVAL)
  }

  /**
   * Clear countdown
   */
  @autobind
  clearDisappearTimer() {
    if (this.timer > 0) {
      clearTimeout(this.timer)
      this.timer = 0
    }
  }

  /**
   * Tooltip data source change handle
   * Data is corresponding info when mouse move, null when mouse out
   * @param {ChartTooltipData} data Tooltip data source
   */
  @autobind
  handleTooltipDataChange(data: ChartTooltipData | null) {
    if (data === null) {
      this.setDisappearTimer()
    } else {
      this.clearDisappearTimer()
      this.dataSource = data
    }
  }

  render() {
    const data = this.dataSource

    if (!data) {
      return null
    }

    return (
      <div
        ref={this.tooltip}
        className="lindb-chart-tooltip"
        onMouseEnter={this.clearDisappearTimer}
        onMouseLeave={this.setDisappearTimer}
      >
        {/* Timestamp */}
        <div className="lindb-chart-tooltip__timestamp">
          <Moment format="YYYY-MM-DD HH:mm:ss">{data.time}</Moment>
        </div>

        {/* Data Source */}
        <ul className="lindb-chart-tooltip__list">
          {data.series.map(item => (
            <li className="lindb-chart-tooltip__list-item" key={item.name + item.color}>
              <div className="lindb-chart-tooltip__list-item__icon">
                <i style={{ backgroundColor: item.color }}/>
              </div>
              <span className="lindb-chart-tooltip__list-item__name">{item.name}</span>
              <span className="lindb-chart-tooltip__list-item__value">
                {DataFormatter.formatter(item.value, data.unit)}
              </span>
            </li>
          ))}
        </ul>

        {/* Arrow */}
        <div className="lindb-chart-tooltip__arrow" ref={this.tooltipArrow}/>
      </div>
    )
  }

  /**
   * Calculate ToolTip position left
   * @param {ChartBorderInfo} border Chart border info
   */
  private calculateContainerPosition(border: ChartBorderInfo) {
    const { width: tooltipWidth, height: tooltipHeight } = this.getTooltipSize()

    if (tooltipWidth === 0 || tooltipHeight === 0) {
      return { left: 0, top: 0 }
    }

    const spacing = this.SCREEN_SPACING // gap between tooltip and window
    const pointX = this.getMousePointPosition(border)
    const rightRest = window.innerWidth - (tooltipWidth / 2) - spacing // right safe area
    const leftRest = spacing + tooltipWidth / 2 // left safe area

    const { top, isBottom } = this.getContainerShouldBeWhere(border)

    return {
      left: Math.max(leftRest, Math.min(pointX, rightRest)),
      top: top === null ? 'auto' : top,
      isBottom,
    }
  }

  /**
   * Calculate Arrow position left
   * @param {ChartBorderInfo} border Chart border info
   * @return {number} ToolTip Left
   */
  private calculateArrowPosition(border: ChartBorderInfo): number {
    const { width: tooltipWidth } = this.getTooltipSize()

    const spacing = this.SCREEN_SPACING // gap between tooltip and window
    const pointX = this.getMousePointPosition(border)
    const rightRest = window.innerWidth - (tooltipWidth / 2) - spacing // right safe area
    const leftRest = spacing + tooltipWidth / 2 // left safe area

    // base distance：tooltipWidth / 2
    // minimum left border：Math.min(pointX - leftRest, 0)
    // maximum right border：Math.max(pointX - rightRest, 0)
    return (tooltipWidth / 2) + Math.min(pointX - leftRest, 0) + Math.max(pointX - rightRest, 0)
  }

  private getContainerShouldBeWhere(border: ChartBorderInfo) {
    const topAreaHeight = border.canvasTop // top safe area
    const bottomAreaHeight = window.innerHeight - border.canvasBottom // bottom safe area
    const spacing = this.SCREEN_SPACING
    const offset = 10 // tooltip vertical offset by chart

    const { height: tooltipHeight } = this.getTooltipSize()
    const containerHeight = tooltipHeight + spacing + offset
    const scrollTop = document.documentElement.scrollTop // window scroll distance

    let top: number | null = null
    let left: number | null = null
    let isBottom: boolean = false

    // when vertical direction could contain tooltip
    if (topAreaHeight > containerHeight || bottomAreaHeight > containerHeight) {
      top = topAreaHeight > containerHeight
        ? topAreaHeight - tooltipHeight - offset + scrollTop
        : border.canvasBottom + offset + scrollTop

      isBottom = topAreaHeight <= containerHeight
    }

    return { top, left, isBottom }
  }

  private getTooltipSize() {
    const tooltip = this.tooltip.current
    return { width: tooltip ? tooltip.clientWidth : 0, height: tooltip ? tooltip.clientHeight : 0 }
  }

  /**
   * Calculate Mouse Point X (valid)
   * @param {ChartBorderInfo} border Chart border info
   * @return {any} 当前 Point Position Left
   */
  private getMousePointPosition(border: ChartBorderInfo) {
    // Mouse position
    const point = border.pointX
    // chart left border
    const leftBorder = border.canvasLeft + border.chartOffsetLeftWithCanvas
    // chart right border
    const rightBorder = leftBorder + border.chartWidth
    return Math.min(rightBorder, Math.max(leftBorder, point)) // left and right border contain mouse point
  }
}