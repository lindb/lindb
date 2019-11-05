import * as React from 'react'
import NodeInfo from './NodeInfo'
import { observer } from 'mobx-react'
import { Card, Table, Tag, Divider } from 'antd'
import { getBrokerCluster, getMaster } from '../../service/monitoring'
import { NodeList } from '../../model/Monitoring'
import { dateFormatter } from '../../utils/util'
import StorageCluster from './StorageCluster'

const { Column } = Table

interface OverviewProps {
}

interface OverviewState {
  ip: string
  port: number
  electTime: number
  brokers: NodeList
}

@observer
class Overview extends React.Component<OverviewProps, OverviewState> {
  constructor(props: OverviewProps) {
    super(props)

    this.state = {
      ip: '',
      port: 0,
      electTime: 0,
      brokers: [],
    }
  }

  componentDidMount(): void {
    this.init()
  }

  init() {
    this.getMaster()
    this.getBrokersList()
  }

  async getMaster() {
    const result: any = await getMaster()
    if (result) {
      const { node: { ip, port }, electTime } = result
      this.setState({ ip, port, electTime })
    }
  }

  async getBrokersList() {
    const brokers: any = await getBrokerCluster()
    if (brokers) {
      this.setState({ brokers })
    }
  }

  render() {
    const { ip, port, electTime, brokers } = this.state

    return (
      <div>
        {/* Master */}
        <Card size="small" title="Master" loading={false}>
          {ip}:{port}
          <Tag color="green" style={{ color: "#fff", background: "#353535", padding: "2px", marginLeft: "8px" }}>
            <span style={{ margin: "4px" }}>
              Elect Time: {dateFormatter(electTime)}
            </span>
          </Tag>
        </Card>

        {/* Node */}
        <Card size="small" title="Node">
          <NodeInfo nodes={brokers} />
        </Card>
        {/* Storage Cluster Overview*/}
        <Card size="small" title="Storeage Cluster">
          <Divider orientation="left">Storage Cluster 1</Divider>
          <StorageCluster />
          <Divider orientation="left" style={{ paddingTop: 10 }}>Storage Cluster 2</Divider>
          <StorageCluster />
          <Divider orientation="left" style={{ paddingTop: 10 }}>Storage Cluster 3</Divider>
          <StorageCluster />
        </Card>

        {/* Dead Node */}
        <Card size="small" title="Dead Node" loading={false}>
          <Table dataSource={[]} size="small">
            <Column title="ID" dataIndex="id" key="id" />
            <Column title="Host Name" dataIndex="hostname" key="hostname" />
            <Column title="IP" dataIndex="ip" key="ip" />
            <Column title="TCP Port" dataIndex="tcpPort" key="tcpPort" />
            <Column title="Dead Time" dataIndex="deadTime" key="deadTime" />
          </Table>
        </Card>

        {/* Database */}
        <Card size="small" title="Database" loading={false}>
          <Table dataSource={[]} size="small">
            <Column title="Name" dataIndex="Name" key="Name" />
            <Column title="Num. Shard" dataIndex="Shard" key="Shard" />
            <Column title="Num. Leader" dataIndex="Leader" key="Leader" />
            <Column title="Num. Live Replica" dataIndex="Live" key="Live" />
            <Column title="Num. ISR Replica" dataIndex="ISR" key="ISR" />
            <Column title="Num. Replica" dataIndex="Replica" key="Replica" />
            <Column title="Description" dataIndex="Description" key="Description" />
          </Table>
        </Card>
      </div>
    )
  }
}

export default Overview 
