import { Card, Form } from 'antd'
import TagValuesSelect from 'components/meta/TagValues'
import ViewBoard from 'components/metric/ViewBoard'
import { SystemBoardForRole } from 'config/monitoring/System'
import * as React from 'react'

interface SystemProps {
}

interface SystemState {
}

export default class System extends React.Component<SystemProps, SystemState> {

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
              <TagValuesSelect measurement="lindb.monitor.system.cpu_stat" tagKey="role"/>
            </Form.Item>
            <Form.Item label="Node">
              <TagValuesSelect measurement="lindb.monitor.system.cpu_stat" tagKey="node" mode="tags" watch={["role"]}/>
            </Form.Item>
          </Form>
        </Card>
        <ViewBoard board={SystemBoardForRole} />
      </React.Fragment>
    )
  }
}