import { Card } from 'antd'
import * as React from 'react'
import ViewBoard from '../../components/Chart/ViewBoard'
import { SystemStorageBoard } from '../../config/monitoring/system'

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