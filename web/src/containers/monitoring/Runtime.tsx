import { Card } from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import { RuntimeStorageBoard } from 'config/monitoring/Runtime'
import * as React from 'react'

interface RuntimeProps {
}

interface RuntimeState {
}

export default class Runtime extends React.Component<RuntimeProps, RuntimeState> {

  render() {
    return (
      <React.Fragment>
        <Card>Node</Card>
        <ViewBoard board={RuntimeStorageBoard}/>
      </React.Fragment>
    )
  }
}