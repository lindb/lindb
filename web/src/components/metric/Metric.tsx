import { autobind } from 'core-decorators'
import { observer } from 'mobx-react'
import * as React from 'react'
import { ChartTooltipData, UnitEnum } from '../../model/Metric'
import StoreManager from '../../store/StoreManager'
import Chart from './Chart'
interface MetricProps {
  id: string
  type: string
  unit: UnitEnum
  chart?: any,
  height?: number
  showQL?: boolean
  showError?: boolean
  timeRange?: any
}

interface MetricStatus {
}

@observer
export default class Metric extends React.Component<MetricProps, MetricStatus> {

  componentDidMount(): void {
    const { id, chart } = this.props
    StoreManager.ChartStore.register(id, chart)
    // force change url to triger metric data load
    StoreManager.URLParamStore.forceChange()
  }

  // Clean up
  componentWillUnmount(): void {
    StoreManager.ChartEventStore.deleteSeriesById(this.props.id)
    StoreManager.ChartStore.unRegister(this.props.id)
  }

  @autobind
  handleLegendClick(index: number, status: boolean[]) {
    StoreManager.ChartEventStore.changeSeriesByClick(this.props.id, status)
  }

  @autobind
  handleChartMouseMove(data: ChartTooltipData) {
    StoreManager.ChartEventStore.setTooltipData(data)
  }

  @autobind
  handleChartMouseOut() {
    StoreManager.ChartEventStore.setTooltipData(null)
  }

  render() {
    const { unit, id, type, height } = this.props
    console.log(height)
    return (
      <div className="lindb-metric">
        <React.Fragment>
          <Chart
            type={type}
            unit={unit}
            uuid={id}
            height={height}
            onMouseMove={this.handleChartMouseMove}
            onMouseOut={this.handleChartMouseOut}
          />
          {/* <ChartLegend data={data} onLegendClick={this.handleLegendClick} /> */}
        </React.Fragment>
      </div>
    )
  }
}