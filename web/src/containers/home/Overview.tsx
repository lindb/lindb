import { Card, Tag } from 'antd'
import NodeInfo from 'containers/home/NodeInfo'
import StorageClusterInfo from 'containers/home/StorageClusterInfo'
import { observer } from 'mobx-react'
import { NodeList, StorageCluster } from 'model/Monitoring'
import * as React from 'react'
import { getBrokerCluster, getMaster, listStorageCluster } from 'service/Monitor'
import { dateFormatter } from 'utils/Util'

interface OverviewProps {
}

interface OverviewState {
  hostIp: string
  hostName: string
  grpcPort: number
  electTime: number
  brokers: NodeList
  storageClusters: Array<StorageCluster>
}

@observer
export default class Overview extends React.Component<OverviewProps, OverviewState> {
  constructor(props: OverviewProps) {
    super(props)

    this.state = {
      hostIp: '',
      hostName: '',
      grpcPort: 0,
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
      const { node: { hostIp, hostName, grpcPort }, electTime } = result
      this.setState({ hostIp, hostName, grpcPort, electTime })
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
    const { hostIp, hostName, grpcPort, electTime, brokers, storageClusters } = this.state

    return (
      <div>
        {/* Master */}
        <Card size="small" title="Master" loading={false}>
          {hostName}({hostIp}):{grpcPort}
          <Tag color="lime" style={{ padding: "2px", marginLeft: "8px" }}>
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
        <Card size="small" title="Storage Cluster List">
          <StorageClusterInfo storageClusterList={storageClusters} />
        </Card>
      </div>
    )
  }
}
