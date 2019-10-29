import * as React from 'react'
import NodeInfo from './NodeInfo'
import { observer } from 'mobx-react'
import { Card, Table, Tag } from 'antd'
import { getBrokersList, getMaster } from '../../service/home'
import { BrokersList } from '../../model/service'
import { dateFormatter } from '../../utils/util'

const { Column } = Table

interface HomeProps {
}

interface HomeState {
  ip: string
  port: number
  electTime: number
  brokers: BrokersList
}

@observer
class Home extends React.Component<HomeProps, HomeState> {
  constructor(props: HomeProps) {
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
    const brokers: any = await getBrokersList()
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
          <NodeInfo brokers={brokers} />
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

        {/* Software Attributes */}
        <Card size="small" title="Software Attributes" loading={false}>
          <Table dataSource={[]} size="small">
            <Column title="Name" dataIndex="Name" key="Name" />
            <Column title="Value" dataIndex="Value" key="Value" />
          </Table>
        </Card>
      </div>
    )
  }
}

export default Home
