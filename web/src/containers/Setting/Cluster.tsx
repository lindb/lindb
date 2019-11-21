import * as React from 'react'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import { autobind } from 'core-decorators'
import { FormComponentProps } from 'antd/es/form'
import { Button, Table, Modal, Form, Input, Select, Card, message } from 'antd'
import { createStorage, getStorageClusterList, StorageClusterInfo } from '../../service/Storage'

const { Column } = Table

interface ClusterProps {
}

interface ClusterStatus {
}

@observer
export default class Cluster extends React.Component<ClusterProps, ClusterStatus> {
  @observable listLoading: boolean = false
  @observable createClusterModalVisible: boolean = false
  @observable storageClusterList: StorageClusterInfo[] | undefined = []

  constructor(props: ClusterProps) {
    super(props)
    this.state = {}
  }

  componentDidMount(): void {
    this.fetchStorageClusterList()
  }

  @autobind
  handleToggleModalVisible() {
    this.createClusterModalVisible = !this.createClusterModalVisible
  }

  @autobind
  handleCreateClusterSuccess() {
    this.handleToggleModalVisible()
    // tip
    message.success('success')
    // fetch latest list
    this.fetchStorageClusterList()
  }

  async fetchStorageClusterList() {
    this.listLoading = true
    this.storageClusterList = await getStorageClusterList()
    this.listLoading = false
  }

  render() {
    return (
      <React.Fragment>
        <Card size="small" className="align-right">
          <Button icon="plus-circle" type="primary" onClick={this.handleToggleModalVisible}>New Cluster</Button>
          <CreateClusterModalWithForm
            visible={this.createClusterModalVisible}
            onClose={this.handleToggleModalVisible}
            onSuccess={this.handleCreateClusterSuccess}
          />
        </Card>
        <Table
          dataSource={this.storageClusterList}
          size="small"
          pagination={false}
          loading={this.listLoading}
          rowKey="name"
        >
          <Column title="Name" dataIndex="name" />
          <Column title="Namespace" dataIndex="config.namespace" />
          <Column
            title="EndPoints"
            dataIndex="config.endpoints"
            render={(points: string[]) => JSON.stringify(points)}
          />
        </Table>
      </React.Fragment>
    )
  }
}

/* Cluster Modal */
interface CreateClusterModalProps extends FormComponentProps {
  visible: boolean
  onClose?: () => void
  onSuccess?: () => void
}

class CreateClusterConfig {
  public static CLUSTER_NAME: string = 'ClusterName'
  public static NAMESPACE: string = 'clusterConfigNamespace'
  public static ENDPOINTS: string = 'clusterConfigEndpoints'
}

@observer
class CreateClusterModal extends React.Component<CreateClusterModalProps> {
  @observable isLoading: boolean = false

  @autobind
  handleClose() {
    const { onClose } = this.props
    if (onClose && typeof onClose === 'function') {
      onClose()
    }
  }

  @autobind
  handleSubmit() {
    const { form, onSuccess } = this.props

    form.validateFields(async (err: any, value: any) => {
      if (err) {
        return
      }

      const {
        [CreateClusterConfig.CLUSTER_NAME]: name,
        [CreateClusterConfig.NAMESPACE]: namespace,
        [CreateClusterConfig.ENDPOINTS]: endpoints,
      } = value

      const config = { namespace, endpoints }

      this.isLoading = true
      await createStorage(name, config)
      this.isLoading = false

      // clear info
      form.resetFields()

      // invoke success handler
      if (onSuccess && typeof onSuccess === 'function') {
        onSuccess()
      }
    })
  }

  render() {
    const { visible, form: { getFieldDecorator } } = this.props
    return (
      <Modal
        title="New Cluster"
        visible={visible}
        okText="Submit"
        onOk={this.handleSubmit}
        onCancel={this.handleClose}
        confirmLoading={this.isLoading}
      >
        <Form layout="horizontal" labelCol={{ span: 8 }} wrapperCol={{ span: 14 }}>
          {/* Cluster Name */}
          <Form.Item label="Name">
            {getFieldDecorator(CreateClusterConfig.CLUSTER_NAME, {
              rules: [{ required: true, message: 'Please input cluster name' }],
            })(<Input placeholder="Cluster Name" autoComplete="off" />)}
          </Form.Item>

          {/* Namespace */}
          <Form.Item label="Namespace">
            {getFieldDecorator(CreateClusterConfig.NAMESPACE, {
              rules: [{ required: true, message: 'Please input cluster namespace' }],
            })(<Input placeholder="Namespace" autoComplete="off" />)}
          </Form.Item>

          {/* Endpoints */}
          <Form.Item label="Endpoints">
            {getFieldDecorator(CreateClusterConfig.ENDPOINTS, {
              rules: [{ type: 'array', required: true, message: 'Please input cluster endpoints' }],
            })(
              <Select mode="tags" style={{ width: '100%' }} placeholder="Endpoints">
                <Select.Option key="http://localhost:2379">http://localhost:2379</Select.Option>
              </Select>,
            )}
          </Form.Item>
        </Form>
      </Modal>
    )
  }
}

const CreateClusterModalWithForm = Form.create<CreateClusterModalProps>()(CreateClusterModal)