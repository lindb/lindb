import * as React from 'react'
import { autobind } from 'core-decorators'
import Tips from '../../components/base/tips'
import TimePicker from '../../components/TimePicker'
import { Button, Card, Form, Icon, Select, Switch, Tabs } from 'antd'

const { TabPane } = Tabs
const { Option } = Select

interface SearchProps {
}

interface SearchStatus {
  autoRefresh: boolean
}

export default class Search extends React.Component<SearchProps, SearchStatus> {
  constructor(props: SearchProps) {
    super(props)

    this.state = {
      autoRefresh: false,
    }
  }

  @autobind
  handleAutoRefreshChange(checked: boolean) {
    this.setState({
      autoRefresh: checked
    })
  }

  render() {
    const { autoRefresh } = this.state

    return (
      <div>
        <Card size="small">
          <Tabs defaultActiveKey="Basic" size="small">
            {/* Basic */}
            <TabPane tab={<span><Icon type="info-circle"/>Basic</span>} key="Basic">

              <Form layout="inline">
                <Form.Item label="DataBase">
                  <Select defaultValue="lucy" style={{ width: 120 }}>
                    <Option value="jack">Jack</Option>
                    <Option value="lucy">Lucy</Option>
                    <Option value="disabled" disabled={true}>
                      Disabled
                    </Option>
                    <Option value="Yiminghe">yiminghe</Option>
                  </Select>
                </Form.Item>
                <Form.Item>
                  <Button type="primary" icon="search">Search</Button>
                </Form.Item>
                <Form.Item>
                  <Button type="danger" icon="close">Clear</Button>
                </Form.Item>
              </Form>

            </TabPane>

            {/* Meta Data */}
            <TabPane tab={<span><Icon type="database"/>Meta Data</span>} key="Meta Data">
              Meta Data
            </TabPane>

            {/* Advance */}
            <TabPane tab={<span><Icon type="wallet"/>Advance</span>} key="Advance">
              Advance
            </TabPane>

            {/* Explain */}
            <TabPane tab={<span><Icon type="wallet"/>Explain</span>} key="Explain">
              Explain
            </TabPane>
          </Tabs>
        </Card>

        {/* Refresh And TimePicker */}
        <Card size="small" style={{ textAlign: 'right' }}>
          <Form layout="inline">
            <Form.Item label="Auto Refresh">
              <Switch
                defaultChecked={autoRefresh}
                onChange={this.handleAutoRefreshChange}
              />
            </Form.Item>
            <Form.Item>
              <Select defaultValue="lucy" style={{ width: 120 }} disabled={!autoRefresh}>
                <Option value="jack">Jack</Option>
                <Option value="lucy">Lucy</Option>
                <Option value="disabled" disabled={true}>
                  Disabled
                </Option>
                <Option value="Yiminghe">yiminghe</Option>
              </Select>
            </Form.Item>
            <Form.Item>
              <TimePicker/>
            </Form.Item>
          </Form>
        </Card>

        {/* Info Area */}
        <Card>
          <Tips tip="No Data" icon="warning"/>
        </Card>
      </div>
    )
  }
}