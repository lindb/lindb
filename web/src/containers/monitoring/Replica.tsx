import {Card, Form} from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import * as React from 'react'
import TagValuesSelect from "components/meta/TagValues";
import {ReplicaDashboard} from "config/monitoring/Replica";

interface ReplicaProps {
}

interface ReplicaState {
}

export default class MonitoringReplica extends React.Component<ReplicaProps, ReplicaState> {

    render() {
        return (
            <React.Fragment>
                <Card>
                    <Form layout="inline"
                          style={{
                              width: "calc(100%)",
                              textAlign: "left",
                          }}>
                        <Form.Item label="Node">
                            <TagValuesSelect metric="lindb.replica.local" tagKey="node" mode="tags"/>
                        </Form.Item>
                    </Form>
                </Card>
                <ViewBoard board={ReplicaDashboard}/>
            </React.Fragment>
        )
    }
}