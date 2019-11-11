import { Card, Tag } from 'antd'
import { observer } from 'mobx-react'
import * as React from 'react'
import { NodeList, StorageCluster } from '../../model/Monitoring'
import { getBrokerCluster, getMaster, listStorageCluster } from '../../service/monitoring'
import { dateFormatter } from '../../utils/util'
import NodeInfo from './NodeInfo'
import StorageClusterInfo from './StorageClusterInfo'

interface OverviewProps {
}

interface OverviewState {
  ip: string
  port: number
  electTime: number
  brokers: NodeList
  storageClusters: Array<StorageCluster>
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
      storageClusters: [],
    }
  }

  componentDidMount(): void {
    this.init()
  }

  init() {
    this.getMaster()
    this.getBrokersList()
    this.listStorageCluster()
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

  async listStorageCluster() {
    const storageClusters: any = await listStorageCluster()
    if (storageClusters) {
      this.setState({ storageClusters })
    }
  }

  render() {
    const { ip, port, electTime, brokers, storageClusters } = this.state

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
        <Card size="small" title="Broker Node List">
          <NodeInfo nodes={brokers} />
        </Card>
        {/* Storage Cluster Overview*/}
        <Card size="small" title="Storeage Cluster List">
          <StorageClusterInfo storageClusterList={storageClusters} />
        </Card>
      </div>
    )
  }
}

export default Overview 
