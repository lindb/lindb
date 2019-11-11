import { Badge, Col, Progress, Row, Table } from 'antd'
import * as React from 'react'
import { NodeList } from '../../model/Monitoring'
import { DataFormatter } from '../../utils/dataFormatter'

const uuidv4 = require('uuid/v4')

interface NodeInfoProps extends NodeListTableProps {
}

interface NodeInfoStatus {
}

export default class NodeInfo extends React.Component<NodeInfoProps, NodeInfoStatus> {
  constructor(props: NodeInfoProps) {
    super(props)
    this.state = {}
  }

  render() {
    const { nodes, isStorage } = this.props

    return (<BrokersListTable nodes={nodes} isStorage={isStorage} />)
  }
}

interface NodeListTableProps {
  nodes: NodeList
  isStorage?: boolean
}

class BrokersListTable extends React.Component<NodeListTableProps> {
  render() {
    const { nodes, isStorage } = this.props
    const columns = [
      {
        title: 'Address',
        key: 'record.node.node.hostName',
        render: (text: any, record: any, index: any) => {
          return (
            <div>
              <Badge status="success" />
              {record.node.node.ip + " - " + record.node.node.hostName}
            </div>
          )
        },
      },
      {
        title: 'Uptime',
        dataIndex: 'node.onlineTime',
        render: (text: any, record: any, index: any) => {
          return DataFormatter.transformSeconds(((new Date()).getTime() - record.node.onlineTime) / 1000)
        },
      },
      {
        title: 'CPU',
        dataIndex: 'system.cpus',
      },
      {
        title: 'Capacity Usage',
        dataIndex: 'system.diskStat',
        render: (text: any, record: any, index: any) => {
          return (
            <div>
              <Row>
                <Col span={6} style={{ textAlign: "right", marginRight: 6 }}>
                  {DataFormatter.transformPercent(record.system.diskStat.usedPercent)}
                </Col>
                <Col span={16}>
                  <Progress style={{ marginTop: 0 }} className="lindb-storage-capacity" percent={record.system.diskStat.usedPercent} status="success" showInfo={false} />
                </Col>
              </Row>
              <Row style={{ textAlign: "left" }}>
                <Col span={12}>
                  <span className="cluster-status-desc">Used: {DataFormatter.transformBytes(record.system.diskStat.used)}</span>
                </Col>
                <Col span={12}>
                  <span className="cluster-status-desc">Total: {DataFormatter.transformBytes(record.system.diskStat.total)}</span>
                </Col>
              </Row>
            </div>
          )
        },
      },
      {
        title: 'Memory Usage',
        dataIndex: 'system.memoryStat',
        render: (text: any, record: any, index: any) => {
          return (
            <div>
              <Row>
                <Col span={6} style={{ textAlign: "right", marginRight: 6 }}>
                  {DataFormatter.transformPercent(record.system.memoryStat.usedPercent)}
                </Col>
                <Col span={16}>
                  <Progress style={{ marginTop: 0 }} className="lindb-storage-capacity" percent={record.system.memoryStat.usedPercent} status="success" showInfo={false} />
                </Col>
              </Row>
              <Row style={{ textAlign: "left" }}>
                <Col span={12}>
                  <span className="cluster-status-desc">Used: {DataFormatter.transformBytes(record.system.memoryStat.used)}</span>
                </Col>
                <Col span={12}>
                  <span className="cluster-status-desc">Total: {DataFormatter.transformBytes(record.system.memoryStat.total)}</span>
                </Col>
              </Row>
            </div>
          )
        },
      },
      {
        title: 'Version',
        dataIndex: 'node.version',
      },
    ];
    if (isStorage) {
      console.log("ssssss", columns)
      columns.splice(2, 0, {
        title: 'Replicas',
        dataIndex: 'replicas',
      })
    }
    console.log("after", columns)
    return (
      <Table dataSource={nodes} bordered={true} rowKey={(record: any) => { return uuidv4() }} size="small" columns={columns} pagination={false} />
    )
  }
}