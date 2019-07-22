import * as React from 'react'
import NodeInfo from './NodeInfo'
import { observer } from 'mobx-react'
import { Card, Table, Tag } from 'antd'

const { Column } = Table

interface HomeProps {
}

interface HomeState {
}

@observer
class Home extends React.Component<HomeProps, HomeState> {

  constructor(props: Readonly<HomeProps>) {
    super(props)
    this.state = {}
  }

  onTabChange = (key: string) => {
    this.setState({ key })
  }

  render() {
    const data = []

    return (
      <div>
        {/* Master */}
        <Card size="small" title="Master" loading={false}>
          adca-infra-etrace-lindb-1.vm.elenet.me(2890)
          <Tag color="green">Elect Time: 2019-06-27 15:23:40</Tag>
        </Card>

        {/* Node */}
        <Card size="small" title="Node">
          <NodeInfo/>
        </Card>

        {/* Dead Node */}
        <Card size="small" title="Dead Node" loading={false}>
          <Table dataSource={data} size="small">
            <Column title="ID" dataIndex="id" key="id"/>
            <Column title="Host Name" dataIndex="hostname" key="hostname"/>
            <Column title="IP" dataIndex="ip" key="ip"/>
            <Column title="TCP Port" dataIndex="tcpPort" key="tcpPort"/>
            <Column title="Dead Time" dataIndex="deadTime" key="deadTime"/>
          </Table>
        </Card>

        {/* Database */}
        <Card size="small" title="Database" loading={false}>
          <Table dataSource={data} size="small">
            <Column title="Name" dataIndex="Name" key="Name"/>
            <Column title="Num. Shard" dataIndex="Shard" key="Shard"/>
            <Column title="Num. Leader" dataIndex="Leader" key="Leader"/>
            <Column title="Num. Live Replica" dataIndex="Live" key="Live"/>
            <Column title="Num. ISR Replica" dataIndex="ISR" key="ISR"/>
            <Column title="Num. Replica" dataIndex="Replica" key="Replica"/>
            <Column title="Description" dataIndex="Description" key="Description"/>
          </Table>
        </Card>

        {/* Software Attributes */}
        <Card size="small" title="Software Attributes" loading={false}>
          <Table dataSource={data} size="small">
            <Column title="Name" dataIndex="Name" key="Name"/>
            <Column title="Value" dataIndex="Value" key="Value"/>
          </Table>
        </Card>
      </div>
    )
  }
}

export default Home
