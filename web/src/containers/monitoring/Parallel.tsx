import {Card, Form} from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import {ParallelBoard} from 'config/monitoring/Parallel'
import * as React from 'react'
import TagValuesSelect from "components/meta/TagValues";

interface RuntimeProps {
}

interface RuntimeState {
}

export default class MonitoringParallel extends React.Component<RuntimeProps, RuntimeState> {

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
                            <TagValuesSelect metric="lindb.parallel" tagKey="node" mode="tags" watch={["pool"]}/>
                        </Form.Item>
                    </Form>
                </Card>
                <ViewBoard board={ParallelBoard}/>
            </React.Fragment>
        )
    }
}