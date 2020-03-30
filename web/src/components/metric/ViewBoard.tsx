import * as React from 'react'
import { Col, Row } from 'antd'
import { Board } from '../../model/Metric'
import { SPACING } from '../../config/config'

import ChartCard from './ChartCard'

interface ViewBoardProps {
  board: Board
}

interface ViewBoardStatus {
}

export default class ViewBoard extends React.Component<ViewBoardProps, ViewBoardStatus> {
  constructor(props: ViewBoardProps) {
    super(props)
    this.state = {}
  }



  render() {
    const { board } = this.props
    return (
      <React.Fragment>
        {board.map((rows, index) => (
          <Row key={index} gutter={SPACING}>
            {rows.map(metric => (
              <Col key={metric.id} span={metric.span}>
                <ChartCard chart={metric.chart} type="line" id={metric.id}/>
              </Col>
            ))}
          </Row>
        ))}
      </React.Fragment>
    )
  }
}