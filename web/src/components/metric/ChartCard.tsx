import { Card } from 'antd'
import Metric from 'components/metric/Metric'
import { Chart } from 'model/Chart'
import * as React from 'react'
import LazyLoad from 'react-lazyload'

interface ChartCardProps {
  chart: Chart
  id: string
  type: string
  height?: number
}

interface ChartCardStatus {
}

export default class ChartCard extends React.Component<ChartCardProps, ChartCardStatus> {
  render() {
    const { id, type, chart, height } = this.props
    const { title, target, unit } = chart

    return (
      <Card title={title} size="small">
        <LazyLoad height={300} once={true} resize={true}>
          {(target && unit) ? (
            <Metric
              id={id}
              height={height}
              unit={unit}
              chart={chart}
              type={type}
            />
          ) : null}
        </LazyLoad>
      </Card>
    )
  }
}