import { Card } from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import { SystemStorageBoard } from 'config/monitoring/Runtime'
import * as React from 'react'

interface MonitoringSystemProps {
}

interface MonitoringSystemStatus {
}

export default class MonitoringSystem extends React.Component<MonitoringSystemProps, MonitoringSystemStatus> {
  constructor(props: MonitoringSystemProps) {
    super(props)
    this.state = {}
  }

  render() {
    return (
      <React.Fragment>
        <Card>Node</Card>
        <ViewBoard board={SystemStorageBoard}/>
      </React.Fragment>
    )
  }
}