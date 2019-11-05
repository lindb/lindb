import * as React from 'react'
import { Icon, Table, Tabs, Row, Col, Tag, Progress } from 'antd'
import { NodeList } from '../../model/Monitoring'
import { dateFormatter } from '../../utils/util'

interface StorageClusterProps {
}

interface StoreageClusterStatus {
}

export default class StoreageCluster extends React.Component<StorageClusterProps, StoreageClusterStatus> {
    renderCapacityUsage() {
        // const usedPercentage = usableCapacity !== 0 ? usedCapacity / usableCapacity : 0;
        return (
            <div>
                <h4 className="cluster-title">Capacity Usage</h4>
                <Row>
                    <Col span={6}>
                        <div className="cluster-status"><span className="success">30.00%</span></div>
                    </Col>
                    <Col span={14}>
                        <Progress className="lindb-storage-capacity" percent={30} status="success" showInfo={false} />
                    </Col>
                </Row>
                <Row>
                    <Col span={6}>
                        <div className="cluster-status-desc">Used Capacity</div>
                    </Col>
                    <Col span={6}>
                        <div className="cluster-status-desc value ">300G</div>
                    </Col>
                    <Col span={6}>
                        <div className="cluster-status-desc">Usable Capacity</div>
                    </Col>
                    <Col span={6}>
                        <div className="cluster-status-desc value">3.1T</div>
                    </Col>
                </Row>
            </div>
        )
    }
    renderNodeStatus() {
        return (
            <div>
                <h4 className="cluster-title">Node Status</h4>
                <Row>
                    <Col span={8}>
                        <div className="cluster-status"><span className="success">10</span></div>
                        <div className="cluster-status-desc">
                            Alive Nodes
                        </div>
                    </Col>
                    <Col span={8}>
                        <div className="cluster-status"><span className="warning">2</span></div>
                        <div className="cluster-status-desc">
                            Suspect nodes
                        </div>
                    </Col>
                    <Col span={8}>
                        <div className="cluster-status"><span className="error">1</span></div>
                        <div className="cluster-status-desc">
                            Dead nodes
                        </div>
                    </Col>
                </Row>
            </div>
        )
    }
    renderReplicationStatus() {
        // const usedPercentage = usableCapacity !== 0 ? usedCapacity / usableCapacity : 0;
        return (
            <div>
                <h4 className="cluster-title">Replication Status</h4>
                <Row>
                    <Col span={8}>
                        <div className="cluster-status"><span className="success">1000</span></div>
                        <div className="cluster-status-desc">
                            Total
                        </div>
                    </Col>
                    <Col span={8}>
                        <div className="cluster-status"><span className="warning">5</span></div>
                        <div className="cluster-status-desc">
                            Under-replicated
                        </div>
                    </Col>
                    <Col span={8}>
                        <div className="cluster-status"><span className="error">10</span></div>
                        <div className="cluster-status-desc">
                            Unavailable
                        </div>
                    </Col>
                </Row>
            </div>
        )
    }
    render() {
        return (
            <Row>
                <Col span={8}>
                    {this.renderCapacityUsage()}
                </Col>
                <Col span={8}>
                    {this.renderNodeStatus()}
                </Col>
                <Col span={8}>
                    {this.renderReplicationStatus()}
                </Col>
            </Row>
        )
    }
}