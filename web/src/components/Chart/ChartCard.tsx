import { Card } from 'antd'
import * as React from 'react'
import { Chart } from '../../model/Metric'

import Metric from './Metric'
import LazyLoad from 'react-lazyload'

interface ChartCardProps {
  chart: Chart
}

interface ChartCardStatus {
}

export default class ChartCard extends React.Component<ChartCardProps, ChartCardStatus> {
  constructor(props: ChartCardProps) {
    super(props)
    this.state = {}
  }

  render() {
    const { chart } = this.props
    return (
      <Card title={chart.title} size="small">
        <LazyLoad height={300} once={true} resize={true}>
          <Metric
            db="_internal"
            ql={chart.ql}
            unit={chart.unit}
            group={chart.group}
          />
        </LazyLoad>
      </Card>
    )
  }
}