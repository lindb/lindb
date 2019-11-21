import { Badge, Col, Progress, Row, Table } from 'antd'
import * as React from 'react'
import { Link } from 'react-router-dom'
import { StorageCluster } from '../../model/Monitoring'
import { DataFormatter } from '../../utils/DataFormatter'
import ReplicaStatusInfo from '../../components/ReplicaStatusInfo'

const uuidv4 = require('uuid/v4')

interface StorageClusterInfoProps {
    storageClusterList: Array<StorageCluster>
}

interface StoreageClusterInfoStatus {
}

export default class StorageClusterInfo extends React.Component<StorageClusterInfoProps, StoreageClusterInfoStatus> {

    /**
     * render capacity usage
     * @param storageCluster 
     */
    renderCapacityUsage(storageCluster: StorageCluster) {
        return (
            <div>
                <Row>
                    <Col span={6}>
                        <div className="cluster-status"><span className="success"> {DataFormatter.transformPercent(storageCluster.capacity.usedPercent)}</span></div>
                    </Col>
                    <Col span={14}>
                        <Progress className="lindb-storage-capacity" percent={storageCluster.capacity.usedPercent} status="success" showInfo={false} />
                    </Col>
                </Row>
                <Row>
                    <Col span={6}>
                        <div className="cluster-status-desc">Used Capacity</div>
                    </Col>
                    <Col span={6}>
                        <div className="cluster-status-desc value ">{DataFormatter.transformBytes(storageCluster.capacity.used)}</div>
                    </Col>
                    <Col span={6}>
                        <div className="cluster-status-desc">Usable Capacity</div>
                    </Col>
                    <Col span={6}>
                        <div className="cluster-status-desc value">{DataFormatter.transformBytes(storageCluster.capacity.total - storageCluster.capacity.used)}</div>
                    </Col>
                </Row>
            </div>
        )
    }

    /**
     * render node status
     * @param storageCluster 
     */
    renderNodeStatus(storageCluster: StorageCluster) {
        return (
            <div>
                <Row>
                    <Col span={8}>
                        <div className="cluster-status"><span className="success">{storageCluster.nodeStatus.alive}</span></div>
                        <div className="cluster-status-desc">
                            Alive Nodes
                        </div>
                    </Col>
                    <Col span={8}>
                        <div className="cluster-status"><span className="warning">{storageCluster.nodeStatus.suspect}</span></div>
                        <div className="cluster-status-desc">
                            Suspect nodes
                        </div>
                    </Col>
                    <Col span={8}>
                        <div className="cluster-status"><span className="error">{storageCluster.nodeStatus.dead}</span></div>
                        <div className="cluster-status-desc">
                            Dead nodes
                        </div>
                    </Col>
                </Row>
            </div>
        )
    }

    /**
     * render replication status
     * @param storageCluster 
     */
    renderReplicationStatus(storageCluster: StorageCluster) {
        return (
            <div>
                <Row>
                    <Col span={8}>
                        <div className="cluster-status"><span className="success">{storageCluster.replicaStatus.total}</span></div>
                        <div className="cluster-status-desc">
                            Total
                        </div>
                    </Col>
                    <Col span={8}>
                        <div className="cluster-status"><span className="warning">{storageCluster.replicaStatus.underReplicated}</span></div>
                        <div className="cluster-status-desc">
                            Under-replicated
                        </div>
                    </Col>
                    <Col span={8}>
                        <div className="cluster-status"><span className="error">{storageCluster.replicaStatus.unavailable}</span></div>
                        <div className="cluster-status-desc">
                            Unavailable
                        </div>
                    </Col>
                </Row>
            </div>
        )
    }
    render() {
        const { storageClusterList } = this.props
        const columns = [
            {
                title: 'Name',
                key: 'name',
                render: (text: any, record: any, index: any) => {
                    return (
                        <div>
                            <Link to={"/storage/cluster/" + record.name}>
                                <Badge status="success" />
                                {record.name}
                            </Link>
                        </div>
                    )
                },
            },
            {
                title: 'Capacity Usage',
                width: "30%",
                render: (text: any, record: any, index: any) => {
                    return this.renderCapacityUsage(record)
                },
            },
            {
                title: 'Node Status',
                width: "30%",
                render: (text: any, record: any, index: any) => {
                    return this.renderNodeStatus(record)
                },
            },
            {
                title: 'Replication Status',
                width: "30%",
                render: (text: any, record: any, index: any) => {
                    return <ReplicaStatusInfo replicaStatus={record.replicaStatus} />
                },
            },
        ]
        return (
            <Table dataSource={storageClusterList} bordered={true} rowKey={(record: any) => { return uuidv4() }} size="small" columns={columns} pagination={false} />
        )
    }
}