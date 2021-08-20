import {Card, Form} from 'antd'
import ViewBoard from 'components/metric/ViewBoard'
import * as React from 'react'
import TagValuesSelect from "components/meta/TagValues";
import {TrafficDashboard} from "config/monitoring/Traffic";

interface TrafficProps {
}

interface TrafficState {
}

export default class MonitoringTraffic extends React.Component<TrafficProps, TrafficState> {

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
                            <TagValuesSelect metric="lindb.traffic.grpc_server" tagKey="node" mode="tags"/>
                        </Form.Item>
                    </Form>
                </Card>
                <ViewBoard board={TrafficDashboard}/>
            </React.Fragment>
        )
    }
}