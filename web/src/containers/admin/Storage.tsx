import { CloseOutlined, PlusOutlined, SaveOutlined } from '@ant-design/icons'
import { Button, Card, Form, Input, message, Modal, Select, Table } from 'antd'
import { FormInstance } from 'antd/lib/form'
import { autobind } from 'core-decorators'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import { StorageCluster, StorageConfig } from 'model/admin/Storage'
import * as React from 'react'
import StoreManager from 'store/StoreManager'

const columns = [
  {
    title: 'Name',
    dataIndex: 'name',
    key: 'name',
  }, {
    title: 'Namespace',
    dataIndex: 'config',
    render: (config: StorageConfig) => (
      <span>{config.namespace}</span>
    ),
  }, {
    title: 'Endpoints',
    dataIndex: 'config',
    render: (config: StorageConfig) => (
      <span>{JSON.stringify(config.endpoints)} </span>
    ),
  }
]
interface StorageProps {
}

interface StorageStatus {
}

@observer
export default class Storage extends React.Component<StorageProps, StorageStatus> {
  @observable modalVisible: boolean = false
  @observable confirmLoading: boolean = false

  constructor(props: StorageProps) {
    super(props)
    this.state = {}
  }

  componentDidMount(): void {
    StoreManager.StorageStore.fetchStorageList()
  }

  @autobind
  handleToggleModalVisible() {
    this.modalVisible = !this.modalVisible
  }

  @autobind
  async handleCreateStorage(config: StorageCluster, form: FormInstance) {
    this.confirmLoading = true
    await StoreManager.StorageStore.createStorageConfig(config)
    this.confirmLoading = false

    // clear form fields
    form.resetFields()

    this.handleToggleModalVisible()
    // tip
    message.success('Create storage config successfully!')
    // fetch latest list
    StoreManager.StorageStore.fetchStorageList()
  }

  render() {
    const storageClusterList = StoreManager.StorageStore.storageList
    const loading = StoreManager.StorageStore.loading

    return (
      <React.Fragment>
        <Card size="small" className="align-right">
          <Button type="primary" icon={<PlusOutlined />} onClick={this.handleToggleModalVisible}>New Storage</Button>
          <StorageForm
            visible={this.modalVisible}
            confirmLoading={this.confirmLoading}
            onCreate={this.handleCreateStorage}
            onCancel={this.handleToggleModalVisible}
          />
        </Card>
        <Table
          dataSource={storageClusterList}
          size="small"
          columns={columns}
          pagination={false}
          loading={loading}
          rowKey="name"
        />
      </React.Fragment>
    )
  }
}

interface StorageFormProps {
  visible: boolean;
  confirmLoading: boolean;
  onCreate: (cluster: StorageCluster, form: FormInstance) => void;
  onCancel: () => void;
}

const StorageForm: React.FC<StorageFormProps> = ({
  visible,
  confirmLoading,
  onCreate,
  onCancel,
}) => {
  const [form] = Form.useForm();
  return (
    <Modal
      width={800}
      confirmLoading={confirmLoading}
      visible={visible}
      title="Storage configuration"
      okText="Save"
      okButtonProps={{ icon: <SaveOutlined /> }}
      cancelText="Cancel"
      cancelButtonProps={{ icon: <CloseOutlined /> }}
      onCancel={onCancel}
      onOk={() => {
        form
          .validateFields()
          .then((cluster: any) => {
            onCreate(cluster, form);
          })
      }}
    >
      <Form
        form={form}
        layout="horizontal" labelCol={{ span: 8 }} wrapperCol={{ span: 14 }}
      >
        <Form.Item
          name="name"
          label="Name"
          rules={[{ required: true, message: 'Please input the storage cluster name!' }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          name={['config', 'namespace']}
          label="Namespace"
          rules={[{ required: true, message: 'Please input the storage cluster namespace!' }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          name={['config', 'endpoints']}
          label="Endpoints"
          rules={[{ required: true, message: 'Please input the storage cluster endpoints!' }]}
        >
          <Select
            mode="tags"
            placeholder="Please input endpoint"
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};