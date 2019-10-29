import { autobind } from 'core-decorators'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import * as React from 'react'
import { ChartTooltipData, ResultSet, UnitEnum } from '../../model/Metric'
import StoreManager from '../../store/StoreManager'
import Chart from './Chart'


const uuidv4 = require('uuid/v4')

interface MetricProps {
  id: string
  unit: UnitEnum
  group?: object
  showQL?: boolean
  showError?: boolean
  timeRange?: any
}

interface MetricStatus {
}

@observer
export default class Metric extends React.Component<MetricProps, MetricStatus> {
  chartResult: ResultSet | null = null
  @observable loading: boolean = false
  private readonly uuid: string

  constructor(props: MetricProps) {
    super(props)

    this.uuid = uuidv4() // Set unique ID
  }

  // Clean up
  componentWillUnmount(): void {
    StoreManager.ChartEventStore.deleteSeriesById(this.uuid)
  }

  @autobind
  handleLegendClick(index: number, status: boolean[]) {
    StoreManager.ChartEventStore.changeSeriesByClick(this.uuid, status)
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
    const { unit, id } = this.props
    return (
      <div className="lindb-metric">
          <React.Fragment>
            <Chart
              type="line"
              unit={unit}
              uuid={id}
              onMouseMove={this.handleChartMouseMove}
              onMouseOut={this.handleChartMouseOut}
            />
            {/* <ChartLegend data={data} onLegendClick={this.handleLegendClick} /> */}
          </React.Fragment>
      </div>
    )
  }
}