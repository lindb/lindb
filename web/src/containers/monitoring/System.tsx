import { Card, Form } from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import ServerRoleSelect from 'components/meta/ServerRole'
import { SystemBoardForRole } from 'config/monitoring/System'
import * as React from 'react'

interface MonitoringSystemProps {
}

interface MonitoringSystemState {
}

export default class MonitoringSystem extends React.Component<MonitoringSystemProps, MonitoringSystemState> {
  constructor(props: MonitoringSystemProps) {
    super(props)
    this.state = {}
  }

  render() {
    return (
      <React.Fragment>
        <Card>
          <Form layout="inline"
            style={{
              width: "calc(100%)",
              textAlign: "left",
            }} >
            <Form.Item label="Server Role">
              <ServerRoleSelect />
            </Form.Item>
          </Form>
        </Card>
        <ViewBoard board={SystemBoardForRole} />
      </React.Fragment>
    )
  }
}