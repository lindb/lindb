import { Card, Form } from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import { BrokerDashboard } from 'config/monitoring/Broker'
import * as React from 'react'
import TagValuesSelect from "components/meta/TagValues";

interface BrokerProps {
}

interface BrokerState {
}

export default class MonitoringSystem extends React.Component<BrokerProps, BrokerState> {

  render() {
    return (
      <React.Fragment>
        <Card>
          <Form layout="inline"
            style={{
              width: "calc(100%)",
              textAlign: "left",
            }} >
            <Form.Item label="Node">
              <TagValuesSelect measurement="system_cpu_stat" tagKey="node' where role='broker" mode="tags"/>
            </Form.Item>
          </Form>
        </Card>
        <ViewBoard board={BrokerDashboard} />
      </React.Fragment>
    )
  }
}