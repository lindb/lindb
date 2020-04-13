import { Col, Row } from 'antd'
import ChartCard from 'components/metric/ChartCard'
import { SPACING } from 'config/config'
import { Board } from 'model/Metric'
import * as React from 'react'

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