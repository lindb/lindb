import {Card, Form} from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import {QueryBoard} from 'config/monitoring/Query'
import * as React from 'react'
import TagValuesSelect from "components/meta/TagValues";

interface RuntimeProps {
}

interface RuntimeState {
}

export default class MonitoringQuery extends React.Component<RuntimeProps, RuntimeState> {

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
                            <TagValuesSelect metric="lindb.concurrent" tagKey="node" mode="tags"/>
                        </Form.Item>
                    </Form>
                </Card>
                <ViewBoard board={QueryBoard}/>
            </React.Fragment>
        )
    }
}