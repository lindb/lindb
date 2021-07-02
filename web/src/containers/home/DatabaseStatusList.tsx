import {Badge, Table} from 'antd';
import ReplicaStatusInfo from 'components/ReplicaStatusInfo';
import {DatabaseStatus} from 'model/Monitoring';
import * as React from 'react';
import {uuid} from 'uuidv4';

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
                render: (text: any, record: any, index: any) => {
                    return (
                        <div>
                            <div className="cluster-status">
                                <span className="success">{record.config.numOfShard}</span>
                            </div>
                        </div>
                    )
                },
            },
            {
                title: 'Replica Factor',
                dataIndex: "config.replicaFactor",
                render: (text: any, record: any, index: any) => {
                    return (
                        <div>
                            <div className="cluster-status">
                                <span className="success">{record.config.replicaFactor}</span>
                            </div>
                        </div>
                    )
                },
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
            <Table dataSource={databaseStatusList} bordered={true} rowKey={(record: any) => { return uuid() }} size="small" columns={columns} pagination={false} />
        )
    }
}