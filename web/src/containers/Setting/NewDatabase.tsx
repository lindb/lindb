import * as React from 'react'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import { autobind } from 'core-decorators'
import { redirectTo } from '../../utils/URLUtil'
import { FormComponentProps } from 'antd/es/form'
import { Button, Card, Col, Form, Input, InputNumber, message, Row, Select, Switch } from 'antd'
import { createDatabase, DatabaseCluster, getStorageClusterList, StorageClusterInfo } from '../../service/storage'

interface CreateDatabaseProps extends FormComponentProps {
}

class CreateDatabaseConfig {
  public static DATABASE_NAME: string = 'DatabaseName'
  public static CLUSTER_NAME: string = 'DatabaseClusterName'
  public static CLUSTER_NUM_OF_SHARD: string = 'DatabaseClusterNumOfShard'
  public static CLUSTER_REPLICA_FACTOR: string = 'DatabaseClusterReplicaFactor'

  public static CLUSTER_ENGINE_INTERVAL: string = 'DatabaseClusterEngineInterval'
  public static CLUSTER_ENGINE_TIMEWINDOW: string = 'DatabaseClusterEngineTimeWindow'
  public static CLUSTER_ENGINE_AUTO_CREATE_NS: string = 'DatabaseClusterEngineAutoCreateNS'
  public static CLUSTER_ENGINE_BEHIND: string = 'DatabaseClusterEngineBehind'
  public static CLUSTER_ENGINE_AHEAD: string = 'DatabaseClusterEngineAhead'
  public static CLUSTER_ENGINE_INDEX_TIME_THRESHOLD: string = 'DatabaseClusterEngineIndexTimeThreshold'
  public static CLUSTER_ENGINE_INDEX_SIZE_THRESHOLD: string = 'DatabaseClusterEngineIndexSizeThreshold'
  public static CLUSTER_ENGINE_DATA_TIME_THRESHOLD: string = 'DatabaseClusterEngineDataTimeThreshold'
  public static CLUSTER_ENGINE_DATA_SIZE_THRESHOLD: string = 'DatabaseClusterEngineDataSizeThreshold'
}

@observer
class CreateDatabase extends React.Component<CreateDatabaseProps> {
  @observable isPosting: boolean = false
  @observable isLoadingClusterList: boolean = false
  @observable storageClusterList: StorageClusterInfo[] | undefined = []
  @observable databaseClusterList: any[] = []

  @autobind
  handleSubmit() {
    const { form } = this.props

    form.validateFields(async (err: any, value: any) => {
      if (err) {
        return
      }

      const {
        [ CreateDatabaseConfig.DATABASE_NAME ]: name,
        [ CreateDatabaseConfig.CLUSTER_NAME ]: cluster,
        [ CreateDatabaseConfig.CLUSTER_ENGINE_AHEAD ]: ahead,
        [ CreateDatabaseConfig.CLUSTER_ENGINE_BEHIND ]: behind,
        [ CreateDatabaseConfig.CLUSTER_NUM_OF_SHARD ]: numOfShard,
        [ CreateDatabaseConfig.CLUSTER_ENGINE_INTERVAL ]: interval,
        [ CreateDatabaseConfig.CLUSTER_REPLICA_FACTOR ]: replicaFactor,
        [ CreateDatabaseConfig.CLUSTER_ENGINE_TIMEWINDOW ]: timeWindow,
        [ CreateDatabaseConfig.CLUSTER_ENGINE_AUTO_CREATE_NS ]: autoCreateNS,
      } = value

      const clusterPayload: DatabaseCluster = {
        name,
        cluster,
        numOfShard,
        replicaFactor,
        engine: {
          interval,
          timeWindow,
          autoCreateNS,
          behind,
          ahead,
          index: {
            timeThreshold: value[ CreateDatabaseConfig.CLUSTER_ENGINE_INDEX_TIME_THRESHOLD ],
            sizeThreshold: value[ CreateDatabaseConfig.CLUSTER_ENGINE_INDEX_SIZE_THRESHOLD ],
          },
          data: {
            timeThreshold: value[ CreateDatabaseConfig.CLUSTER_ENGINE_DATA_TIME_THRESHOLD ],
            sizeThreshold: value[ CreateDatabaseConfig.CLUSTER_ENGINE_DATA_SIZE_THRESHOLD ],
          },
        },
      }

      this.isPosting = true
      await createDatabase(clusterPayload)
      .then(() => {
        this.isPosting = false
        message.success('Creation successful')
        setTimeout(() => redirectTo('/setting/database'), 1000)
      })
      .catch(error => { console.warn(error) })
    })
  }

  @autobind
  async fetchClusterList() {
    const hasList = this.storageClusterList && this.storageClusterList.length > 0

    if (!hasList) {
      this.isLoadingClusterList = true
      this.storageClusterList = await getStorageClusterList()
      this.isLoadingClusterList = false
    }
  }

  render() {
    const { form: { getFieldDecorator } } = this.props

    const formLayout = {
      labelCol: {
        xs: { span: 11 },
        md: { span: 8 },
        lg: { span: 8 },
        xl: { span: 10 },
        xxl: { span: 6 },
      },
      wrapperCol: {
        xs: { span: 13 },
        md: { span: 16 },
        lg: { span: 16 },
        xl: { span: 14 },
        xxl: { span: 16 },
      },
    }

    return (
      <Card size="small" title="New Database">
        <Row>
          <Col
            xl={{ span: 10, offset: 7 }}
            lg={{ span: 16, offset: 4 }}
            md={{ span: 18, offset: 3 }}
            xs={{ span: 24, offset: 0 }}
          >
            <Form layout="horizontal" {...formLayout}>
              {/* Database Name */}
              <Form.Item label="Name" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.DATABASE_NAME, {
                  rules: [ { required: true, message: 'Please input database name' } ],
                })(<Input placeholder="Database Name" style={{ width: '100%' }}/>)}
              </Form.Item>

              <Form.Item label="Cluster" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_NAME, {
                  rules: [ { required: true } ],
                })(
                  <Select
                    // style={{ width: 180 }}
                    style={{ width: '100%' }}
                    placeholder="Select a cluster"
                    showSearch={true}
                    loading={this.isLoadingClusterList}
                    onFocus={this.fetchClusterList}
                  >
                    {this.storageClusterList && this.storageClusterList.map(item => (
                      <Select.Option key={item.name}>{item.name}</Select.Option>
                    ))}
                  </Select>,
                )}
              </Form.Item>

              {/* Number of Shard */}
              <Form.Item label="Number of Shard" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_NUM_OF_SHARD, {
                  initialValue: 1,
                  rules: [ { required: true } ],
                })(<InputNumber style={{ width: '100%' }}/>)}
              </Form.Item>

              {/* Replica Factor */}
              <Form.Item label="Replica Factor" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_REPLICA_FACTOR, {
                  initialValue: 1,
                  rules: [ { required: true } ],
                })(<InputNumber style={{ width: '100%' }}/>)}
              </Form.Item>

              {/* Interval */}
              <Form.Item label="Interval" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_ENGINE_INTERVAL, {
                  initialValue: '10s',
                  rules: [ { required: true } ],
                })(<Input placeholder="Interval"/>)}
              </Form.Item>

              {/* Time Window */}
              <Form.Item label="Time Window" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_ENGINE_TIMEWINDOW, {
                  initialValue: 0,
                })(<InputNumber style={{ width: '100%' }}/>)}
              </Form.Item>

              {/* Aut Create NS */}
              <Form.Item label="Auto Create NS" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_ENGINE_AUTO_CREATE_NS, {
                  initialValue: true,
                  valuePropName: 'checked',
                  rules: [ { required: true } ],
                })(<Switch/>)}
              </Form.Item>

              {/* Behind */}
              <Form.Item label="Behind" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_ENGINE_BEHIND, {
                  initialValue: '',
                })(<Input placeholder="Behind"/>)}
              </Form.Item>

              {/* Ahead */}
              <Form.Item label="Ahead" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_ENGINE_AHEAD, {
                  initialValue: '',
                })(<Input placeholder="Ahead"/>)}
              </Form.Item>

              {/* Index Time Threshold */}
              <Form.Item label="Index Time Threshold" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_ENGINE_INDEX_TIME_THRESHOLD, {
                  // initialValue: ,
                })(<InputNumber style={{ width: '100%' }}/>)}
              </Form.Item>

              {/* Index Size Threshold */}
              <Form.Item label="Index Size Threshold" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_ENGINE_INDEX_SIZE_THRESHOLD, {
                  // initialValue: ,
                })(<InputNumber style={{ width: '100%' }}/>)}
              </Form.Item>

              {/* Data Time Threshold */}
              <Form.Item label="Data Time Threshold" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_ENGINE_DATA_TIME_THRESHOLD, {
                  // initialValue: ,
                })(<InputNumber style={{ width: '100%' }}/>)}
              </Form.Item>

              {/* Data Size Threshold */}
              <Form.Item label="Data Size Threshold" help={null}>
                {getFieldDecorator(CreateDatabaseConfig.CLUSTER_ENGINE_DATA_SIZE_THRESHOLD, {
                  // initialValue: ,
                })(<InputNumber style={{ width: '100%' }}/>)}
              </Form.Item>

              <Form.Item
                wrapperCol={{
                  xs: { offset: 11 },
                  md: { offset: 8 },
                  lg: { offset: 8 },
                  xl: { offset: 10 },
                  xxl: { offset: 6 },
                }}
              >
                <Button type="primary" icon="save" onClick={this.handleSubmit}>Create</Button>
              </Form.Item>
            </Form>
          </Col>
        </Row>
      </Card>
    )
  }
}

export const CreateDatabaseWithForm = Form.create<CreateDatabaseProps>()(CreateDatabase)