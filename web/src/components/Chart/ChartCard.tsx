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
  render() {
    const { chart } = this.props
    const { title, ql, unit, group } = chart

    return (
      <Card title={title} size="small">
        <LazyLoad height={300} once={true} resize={true}>
          {(ql && unit) ? (
            <Metric
              id="fkdsjfksj"
              unit={unit}
              group={group}
            />
          ) : null}
        </LazyLoad>
      </Card>
    )
  }
}