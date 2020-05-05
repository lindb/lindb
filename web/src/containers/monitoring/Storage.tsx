import { Card, Form } from 'antd'
import ServerRoleSelect from 'components/meta/ServerRole'
import ViewBoard from 'components/metric/ViewBoard'
import { WriteDashboard } from 'config/monitoring/Storage'
import * as React from 'react'

interface StorageProps {
}

interface StorageState {
}

export default class MonitoringSystem extends React.Component<StorageProps, StorageState> {

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
        <ViewBoard board={WriteDashboard} />
      </React.Fragment>
    )
  }
}