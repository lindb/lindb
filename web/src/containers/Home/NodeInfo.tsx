import * as React from 'react'
import { Icon, Table, Tabs } from 'antd'
import { BrokersList } from '../../model/service'
import { dateFormatter } from '../../utils/util'

const { TabPane } = Tabs
const { Column } = Table

interface NodeInfoProps extends BrokersListTableProps {
}

interface NodeInfoStatus {
}

export default class NodeInfo extends React.Component<NodeInfoProps, NodeInfoStatus> {
  constructor(props: NodeInfoProps) {
    super(props)
    this.state = {}
  }

  render() {
    const { brokers } = this.props

    return (
      <Tabs defaultActiveKey="Basic" size="small">
        <TabPane tab={<span><Icon type="info-circle" />Brokers</span>} key="Basic">
          {brokers && brokers.length > 0 && <BrokersListTable brokers={brokers} />}
        </TabPane>
        <TabPane tab={<span><Icon type="database" />Memory</span>} key="Memory">
          Memory
        </TabPane>
        <TabPane tab={<span><Icon type="wallet" />Storage</span>} key="Storage">
          Storage
        </TabPane>
      </Tabs>
    )
  }
}

interface BrokersListTableProps {
  brokers: BrokersList
}

class BrokersListTable extends React.Component<BrokersListTableProps> {
  render() {
    const { brokers } = this.props

    return (
      <Table dataSource={brokers} size="small" rowKey="onlineTime" style={{ border: "none" }} pagination={false}>
        <Column title="Host Name" dataIndex="node.hostName" />
        <Column title="IP" dataIndex="node.ip" />
        <Column title="Port" dataIndex="node.port" />
        <Column title="Online Time" dataIndex="onlineTime" render={(time: number) => dateFormatter(time)} />
      </Table>
    )
  }
}