import { CloseOutlined, DeleteOutlined, EditOutlined, PlusOutlined, SaveOutlined } from '@ant-design/icons'
import { Button, Card, Form, Input, InputNumber, message, Modal, Select, Table } from 'antd'
import { FormInstance } from 'antd/lib/form'
import { autobind } from 'core-decorators'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import { DatabaseConfig, DefaultDatabaseConfig } from 'model/admin/Database'
import { StorageCluster } from 'model/admin/Storage'
import * as React from 'react'
import StoreManager from 'store/StoreManager'


interface DatabaseProps {
}

interface DatabaseStatus {
}

@observer
export default class Database extends React.Component<DatabaseProps, DatabaseStatus> {
  @observable modalVisible: boolean = false
  @observable confirmLoading: boolean = false
  config?: DatabaseConfig

  componentDidMount(): void {
    StoreManager.DatabaseStore.fectchDatabaseList()
  }

  @autobind
  handleToggleModalVisible(database?: DatabaseConfig) {
    this.modalVisible = !this.modalVisible
    this.config = database
    if (this.modalVisible) {
      // if modal visible, fectch newestly storage list
      StoreManager.StorageStore.fetchStorageList()
    }
  }
  @autobind
  async handleCreateDatabase(config: DatabaseConfig, form: FormInstance) {
    this.confirmLoading = true
    await StoreManager.DatabaseStore.createDatabase(config)
    this.confirmLoading = false

    // clear form fields
    form.resetFields()

    this.handleToggleModalVisible()
    // tip
    message.success('Save database configuration successfully!')
    // fetch latest list
    StoreManager.DatabaseStore.fectchDatabaseList()
  }
  render() {
    const databaseList = StoreManager.DatabaseStore.databaseList
    const loading = StoreManager.DatabaseStore.loading
    const storageList = StoreManager.StorageStore.storageList

    const columns = [
      {
        title: 'Name',
        dataIndex: 'name',
        key: 'name',
      }, {
        title: 'Description',
        dataIndex: 'desc',
      }, {
        title: 'Action',
        render: ((text: any, record: DatabaseConfig) => (
          <div>
            <Button type="primary" icon={<EditOutlined />} onClick={() => this.handleToggleModalVisible(record)}>Edit</Button>
            <Button type="primary" danger icon={<DeleteOutlined />}>Delete</Button>
          </div>
        ))
      }
    ]
    return (
      <React.Fragment>
        <Card size="small" className="align-right">
          <Button icon={<PlusOutlined />} type="primary" onClick={() => this.handleToggleModalVisible(DefaultDatabaseConfig)}>New Database</Button>
          <DatabaseForm
            visible={this.modalVisible}
            storageList={storageList}
            values={this.config}
            confirmLoading={this.confirmLoading}
            onCreate={this.handleCreateDatabase}
            onCancel={this.handleToggleModalVisible}
          />
        </Card>

        <Table
          dataSource={databaseList}
          size="small"
          loading={loading}
          columns={columns}
          pagination={false}
          rowKey="name"
        />
      </React.Fragment>
    )
  }
}

/* Database form Modal */
interface DatabaseFormProps {
  visible: boolean;
  confirmLoading: boolean;
  storageList: StorageCluster[] | undefined;
  values?: DatabaseConfig;
  onCreate: (database: DatabaseConfig, form: FormInstance) => void;
  onCancel: () => void;
}

const DatabaseForm: React.FC<DatabaseFormProps> = ({
  visible,
  confirmLoading,
  storageList,
  values,
  onCreate,
  onCancel,
}) => {
  const [form] = Form.useForm();
  return (
    <Modal
      width={1000}
      confirmLoading={confirmLoading}
      visible={visible}
      title="Database configuration"
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
        initialValues={values}
        layout="horizontal" labelCol={{ span: 6 }} wrapperCol={{ span: 16 }}
      >
        <Form.Item
          name="name"
          label="Name"
          rules={[{ required: true, message: 'Please input the database name!' }]}
        >
          <Input placeholder="Please input database name!" />
        </Form.Item>
        <Form.Item
          name="storage"
          label="Storage Cluster"
          rules={[{ required: true, message: 'Please select the storage cluster!' }]}
        >
          <Select
            style={{ width: '100%' }}
            placeholder="Please select a storage cluster"
            showSearch={true}
          >
            {storageList && storageList.map(item => (
              <Select.Option key={item.name} value={item.name}>{item.name}</Select.Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item
          label="Num. of Shard"
          name="numOfShard"
          rules={[{ required: true, message: 'Please input num. of shard!' }]}
        >
          <InputNumber style={{ width: '100%' }} min={1} />
        </Form.Item>
        <Form.Item
          label="Replica Factor"
          name="replicaFactor"
          rules={[{ required: true, message: 'Please input replica factor!' }]}
        >
          <InputNumber style={{ width: '100%' }} min={1} />
        </Form.Item>
        <Form.Item
          label="Interval"
          name={["option", "interval"]}
          rules={[{ required: true, message: 'Please input interval!' }]}
        >
          <Input placeholder="Please input interval(For example: 10s,1min,5min etc.)" />
        </Form.Item>
      </Form>
    </Modal>
  );
};