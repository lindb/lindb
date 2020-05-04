import { Card } from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import { RuntimeStorageBoard } from 'config/monitoring/Runtime'
import * as React from 'react'

interface MonitoringRuntimeProps {
}

interface MonitoringRuntimeState {
}

export default class MonitoringRuntime extends React.Component<MonitoringRuntimeProps, MonitoringRuntimeState> {
  constructor(props: MonitoringRuntimeProps) {
    super(props)
    this.state = {}
  }

  render() {
    return (
      <React.Fragment>
        <Card>Node</Card>
        <ViewBoard board={RuntimeStorageBoard}/>
      </React.Fragment>
    )
  }
}