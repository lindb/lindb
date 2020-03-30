import { Badge, Table } from 'antd'
import * as React from 'react'
import ReplicaStatusInfo from '../../components/ReplicaStatusInfo'
import { DatabaseStatus } from '../../model/Monitoring'

const uuidv4 = require('uuid/v4')

interface DatabaseStatusListProps {
    databaseStatusList: Array<DatabaseStatus>
}

interface DatabaseStatusListStatus {
}

export default class DatabaseStatusList extends React.Component<DatabaseStatusListProps, DatabaseStatusListStatus> {

    render() {
        const { databaseStatusList } = this.props
        const columns = [
            {
                title: 'Name',
                key: 'config.name',
                render: (text: any, record: any, index: any) => {
                    return (
                        <div>
                            <Badge status="success" />
                            {record.config.name}
                        </div>
                    )
                },
            },
            {
                title: 'Num. Of Shards',
                dataIndex: "config.numOfShard",
            },
            {
                title: 'Replica Factor',
                dataIndex: "config.replicaFactor",
            },
            {
                title: 'Replication Status',
                width: "30%",
                render: (text: any, record: any, index: any) => {
                    return <ReplicaStatusInfo replicaStatus={record.replicaStatus} />
                },
            },
            {
                title: 'Description',
                dataIndex: 'config.desc',
            },
        ]
        return (
            <Table dataSource={databaseStatusList} bordered={true} rowKey={(record: any) => { return uuidv4() }} size="small" columns={columns} pagination={false} />
        )
    }
}