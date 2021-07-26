import {Card, Form} from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import {WriteDashboard} from 'config/monitoring/Storage'
import * as React from 'react'
import TagValuesSelect from "components/meta/TagValues";

interface StorageProps {
}

interface StorageState {
}

export default class MonitoringStorage extends React.Component<StorageProps, StorageState> {

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
                            <TagValuesSelect metric="lindb.tsdb.memdb" tagKey="node" mode="tags"/>
                        </Form.Item>
                    </Form>
                </Card>
                <ViewBoard board={WriteDashboard}/>
            </React.Fragment>
        )
    }
}