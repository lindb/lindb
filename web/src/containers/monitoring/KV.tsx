import {Card, Form} from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import * as React from 'react'
import TagValuesSelect from "components/meta/TagValues";
import {KVDashboard} from "config/monitoring/KV";

interface KVProps {
}

interface KVState {
}

export default class MonitoringKV extends React.Component<KVProps, KVState> {

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
                            <TagValuesSelect metric="lindb.kv.table.reader" tagKey="node" mode="tags"/>
                        </Form.Item>
                    </Form>
                </Card>
                <ViewBoard board={KVDashboard}/>
            </React.Fragment>
        )
    }
}