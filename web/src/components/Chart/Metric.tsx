import * as React from 'react'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import { autobind } from 'core-decorators'
import { search } from '../../service/metric'
import StoreManager from '../../store/StoreManager'
import { LineChart } from '../../utils/ProcessChartData'
import { ChartTooltipData, ResultSet, UnitEnum } from '../../model/Metric'

import Chart from './Chart'
import ChartLegend from './ChartLegend'

const uuidv4 = require('uuid/v4')

interface MetricProps {
  db: string
  ql: string
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
  @observable chartResult: ResultSet | null = null
  @observable loading: boolean = false
  private readonly uuid: string

  constructor(props: MetricProps) {
    super(props)

    this.uuid = uuidv4() // Set unique ID
  }

  async componentDidMount() {
    const { db, ql } = this.props

    this.loading = false
    const result = await search(db, ql)
    this.loading = true

    this.chartResult = result.data
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
    const { unit } = this.props
    const { data, options, plugins } = LineChart(this.chartResult, unit)

    return (
      <div className="lindb-metric">
        {this.chartResult && (
          <React.Fragment>
            <Chart
              type="line"
              unit={unit}
              data={data}
              uuid={this.uuid}
              options={options}
              plugins={plugins}
              onMouseMove={this.handleChartMouseMove}
              onMouseOut={this.handleChartMouseOut}
            />
            <ChartLegend data={data} onLegendClick={this.handleLegendClick}/>
          </React.Fragment>
        )}
      </div>
    )
  }
}