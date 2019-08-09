import * as React from 'react'
import classNames from 'classnames'
import { autobind } from 'core-decorators'

interface ChartLegendProps {
  data: any
  onLegendClick?: (index: number, status: boolean[]) => void
}

interface ChartLegendStatus {
  status: boolean[] // Array of series status
}

export default class ChartLegend extends React.Component<ChartLegendProps, ChartLegendStatus> {
  private chartLegendCls = 'lindb-chart-legend'

  constructor(props: ChartLegendProps) {
    super(props)

    const { data } = this.props
    const legendLength = data.datasets ? data.datasets.length : 0

    this.state = {
      status: [ ...Array(legendLength + 1) ].join('0').split('').map(() => false),
    }
  }

  @autobind
  handleLegendItemClick(index: number) {
    const status = this.state.status.slice()
    status[ index ] = !this.state.status[ index ]
    this.setState({ status })

    const { onLegendClick } = this.props
    onLegendClick && onLegendClick(index, this.state.status)
  }

  render() {
    const { data } = this.props
    const { status } = this.state
    const cls = this.chartLegendCls

    return (
      <div className={cls}>
        {/* render Legend Item */}
        {data.datasets && data.datasets.map((item: any, index: number) => {
          const { label, borderColor } = item

          return (
            <span
              key={label + borderColor}
              className={classNames(`${cls}__item`, { hidden: status[ index ] })}
              onClick={() => this.handleLegendItemClick(index)}
              title={label}
            >
              <i className={`${cls}__item__icon`} style={{ backgroundColor: borderColor }}/>{label}
            </span>
          )
        })}
      </div>
    )
  }
}