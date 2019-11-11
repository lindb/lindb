import { Col, Row } from 'antd'
import * as React from 'react'
import { ReplicaStatus } from '../model/Monitoring'

interface ReplicaStatusInfoProps {
    replicaStatus: ReplicaStatus,
}

interface ReplicaStatusInfoStatus {
}

export default class ReplicaStatusInfo extends React.Component<ReplicaStatusInfoProps, ReplicaStatusInfoStatus> {

    render() {
        const { replicaStatus } = this.props
        return (
            <Row>
                <Col span={8}>
                    <div className="cluster-status"><span className="success">{replicaStatus.total}</span></div>
                    <div className="cluster-status-desc">
                        Total
                </div>
                </Col>
                <Col span={8}>
                    <div className="cluster-status"><span className="warning">{replicaStatus.underReplicated}</span></div>
                    <div className="cluster-status-desc">
                        Under-replicated
                    </div>
                </Col>
                <Col span={8}>
                    <div className="cluster-status"><span className="error">{replicaStatus.unavailable}</span></div>
                    <div className="cluster-status-desc">
                        Unavailable
                    </div>
                </Col>
            </Row>
        )
    }
}