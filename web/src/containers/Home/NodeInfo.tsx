import * as React from 'react'
import { Icon, Tabs } from 'antd'

const { TabPane } = Tabs

interface NodeInfoProps {
}

interface NodeInfoStatus {
}

export default class NodeInfo extends React.Component<NodeInfoProps, NodeInfoStatus> {
  constructor(props: NodeInfoProps) {
    super(props)
    this.state = {}
  }

  render() {
    return (
      <Tabs defaultActiveKey="Basic" size="small">
        <TabPane tab={<span><Icon type="info-circle" />Basic</span>} key="Basic">
          Basic
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