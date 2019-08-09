import * as React from 'react'
import { getOptions } from '../../config/chartConfig'
import { UnitEnum } from '../../model/Metric'
import { getChartColor, getRandomNumBetween } from '../../utils/util'

import { set } from 'lodash-es'

const ChartJS = require('chart.js')

interface LoginBackgroundProps {
}

interface LoginBackgroundStatus {
}

export default class LoginBackground extends React.Component<LoginBackgroundProps, LoginBackgroundStatus> {
  canvas: React.RefObject<HTMLCanvasElement>
  chartInstance: any // ChartJS instance

  constructor(props: LoginBackgroundProps) {
    super(props)
    this.state = {}

    this.canvas = React.createRef()
  }

  componentDidMount(): void {
    this.init()
  }

  init() {
    const canvas = this.canvas.current
    if (!canvas) {
      return
    }

    const ctx = canvas.getContext('2d')

    const datasets = [ ...Array(5) ].join('.').split('.').map((_, idx) => {
      const color = getChartColor(idx)
      const points = 100
      const now = new Date()
      const step = Math.random() * 10000

      return {
        label: `Test${idx}`,
        fill: false,
        borderColor: color,
        backgroundColor: color,
        pointBackgroundColor: color,
        data: [ ...Array(points) ].join('.').split('.').map((__, i) => ({
          x: new Date(+now - (points - i) * 10 * 1000),
          y: step + step * Math.random() + getRandomNumBetween(-200 * Math.random(), 200 * Math.random()),
        })),
      }
    })
    const options = getOptions(UnitEnum.Bytes) || {}
    set(options, 'elements.line.tension', 0.4)

    this.chartInstance = new ChartJS(ctx, { type: 'line', data: { datasets }, options })
  }

  render() {
    return (
      <div className="lindb-login__background">
        <canvas ref={this.canvas}/>
      </div>
    )
  }
}