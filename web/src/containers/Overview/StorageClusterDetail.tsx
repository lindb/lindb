import { Card } from 'antd'
import * as React from 'react'
import { StorageCluster } from '../../model/Monitoring'
import { getStorageCluster } from '../../service/monitoring'
import DatabaseStatusList from './DatabaseStatusList'
import NodeInfo from './NodeInfo'
import StorageClusterInfo from './StorageClusterInfo'

interface StorageClusterDetailProps {
    match?: any
}

interface StoreageClusterDetailStatus {
    storageCluster?: StorageCluster
}

export default class StorageClusterDetail extends React.Component<StorageClusterDetailProps, StoreageClusterDetailStatus> {
    private clusterName: string

    constructor(props: StorageClusterDetailProps) {
        super(props)
        const { clusterName } = this.props.match.params
        this.clusterName = clusterName
        this.state = {}
    }

    componentDidMount(): void {
        this.getStorageCluster()
    }

    async getStorageCluster() {
        const storageCluster: any = await getStorageCluster(this.clusterName)
        if (storageCluster) {
            this.setState({ storageCluster })
        }
    }

    render() {
        const { storageCluster } = this.state
        if (!storageCluster) {
            return (<div>empty</div>)
        }
        return (
            <div>
                {/* Storage Cluster Overview*/}
                <Card size="small" title="Overview">
                    <StorageClusterInfo storageClusterList={[storageCluster!]} />
                </Card>

                {/* Node */}
                <Card size="small" title="Node List">
                    <NodeInfo nodes={storageCluster.nodes} isStorage={true} />
                </Card>

                {/* Database Status List */}
                <Card size="small" title="Database Status List">
                    <DatabaseStatusList databaseStatusList={storageCluster.databaseStatusList} />
                </Card>
            </div>
        )
    }
}