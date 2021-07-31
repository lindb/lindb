import {Card, Form} from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import {ConcurrentBoard} from 'config/monitoring/Concurrent'
import * as React from 'react'
import TagValuesSelect from "components/meta/TagValues";

interface RuntimeProps {
}

interface RuntimeState {
}

export default class MonitoringConcurrent extends React.Component<RuntimeProps, RuntimeState> {

    render() {
        return (
            <React.Fragment>
                <Card>
                    <Form layout="inline"
                          style={{
                              width: "calc(100%)",
                              textAlign: "left",
                          }}>
                        <Form.Item label="Pool Name">
                            <TagValuesSelect metric="lindb.concurrent" tagKey="pool_name" mode="tags"/>
                        </Form.Item>
                        <Form.Item label="Node">
                            <TagValuesSelect metric="lindb.concurrent" tagKey="node" mode="tags" watch={["pool_name"]}/>
                        </Form.Item>
                    </Form>
                </Card>
                <ViewBoard board={ConcurrentBoard}/>
            </React.Fragment>
        )
    }
}