import { Alert, Button, Card, Col, Form, Icon, Input, Row, Switch } from 'antd'
import { autobind } from 'core-decorators'
import { reaction, toJS } from 'mobx'
import * as React from 'react'
import Metric from '../../components/Chart/Metric'
import DatabaseNameSelect from '../../components/Metadata/DatabaseNames'
import { ChartStatus, ChartStatusEnum } from '../../model/Chart'
import { UnitEnum } from '../../model/Metric'
import StoreManager from '../../store/StoreManager'
import { observer } from 'mobx-react'

interface DataSearchProps {
}

interface DataSearchStatus {
  autoRefresh: boolean
  ql: string;
  chartStatus: ChartStatus;
}

const chartID = "999999999"

@observer
export default class DataSearch extends React.Component<DataSearchProps, DataSearchStatus> {
  disposers: any[]

  constructor(props: DataSearchProps) {
    super(props)

    const chartStatus = new ChartStatus()
    chartStatus.status = ChartStatusEnum.Init
    this.state = {
      ql: StoreManager.URLParamStore.getValue("ql"),
      autoRefresh: false,
      chartStatus: chartStatus,
    }

    this.disposers = [
      reaction(
        () => StoreManager.ChartStore.chartStatusMap.get(chartID),
        chartStatus => {
          this.setState({ chartStatus: toJS(chartStatus!) })
        }
      )
    ]
  }

  componentWillUnmount(): void {
    this.disposers.map(handle => handle())
    StoreManager.ChartStore.unRegister(chartID)
  }

  @autobind
  handleAutoRefreshChange(checked: boolean) {
    this.setState({
      autoRefresh: checked
    })
  }

  @autobind
  async searchQL() {
    const { ql } = this.state;
    StoreManager.ChartStore.reRegister(chartID, {
      target: {
        db: StoreManager.URLParamStore.getValue("db"),
        ql: ql,
      }
    })
    StoreManager.URLParamStore.changeURLParams({ ql: ql })
    StoreManager.URLParamStore.forceChange()
  }

  @autobind
  renderChartStatus() {
    const { chartStatus } = this.state
    if (chartStatus.status === ChartStatusEnum.Init) {
      return
    } else if (chartStatus.status === ChartStatusEnum.Loading) {
      return <Icon type="loading" style={{ fontSize: "24px", color: "#2e81f7" }} />
    } else if (chartStatus.status === ChartStatusEnum.NoData) {
      return <Alert message="Data not found" type="warning" showIcon />
    } else if (chartStatus.status === ChartStatusEnum.UnLimit) {
      return <Alert message="Display limit" type="warning" showIcon />
    } else if (chartStatus.status === ChartStatusEnum.LoadError) {
      return <Alert message={chartStatus.msg} type="error" showIcon />
    } else {
      return <Alert message="Success" type="success" showIcon />
    }
  }

  render() {
    const { autoRefresh } = this.state
    // const ql = StoreManager.URLParamStore.getValue("ql")

    return (
      <div>
        <Card size="small">
          <Form layout="inline"
            style={{
              width: "calc(100%)",
              textAlign: "left",
              display: "flex"
            }} >
            <Form.Item label="Database">
              <DatabaseNameSelect />
            </Form.Item>
            <Form.Item label="Lin QL" className="ql-form-item">
              <Input value={this.state.ql}
                onChange={(value) => this.setState({ ql: value.target.value })}
                placeholder="Please input LinDB query language" />
            </Form.Item>
          </Form>
          <Form
            layout="inline"
            style={{
              width: "calc(100%)",
              textAlign: "center",
              marginTop: 6,
            }}
          >
            <Form.Item>
              <Button type="primary" icon="search" onClick={() => this.searchQL()}>Search</Button>
            </Form.Item>
            <Form.Item>
              <Button type="danger" icon="close">Clear</Button>
            </Form.Item>
          </Form>
        </Card>

        {/* Refresh And TimePicker */}
        <Card size="small" >
          <Row>
            <Col span={20}>
              {this.renderChartStatus()}
            </Col>
            <Col span={4} style={{ textAlign: 'right' }}>
              <Form layout="inline">
                <Form.Item label="Auto Refresh">
                  <Switch
                    defaultChecked={autoRefresh}
                    onChange={this.handleAutoRefreshChange}
                  />
                </Form.Item>
              </Form>
            </Col>
          </Row>
        </Card>
        <Card>
          <Metric
            id={chartID}
            unit={UnitEnum.None}
          />
        </Card>
      </div>
    )
  }
}