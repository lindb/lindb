import * as React from 'react'
import { autobind } from 'core-decorators'
import { Button, Card, Table } from 'antd'
import { redirectTo } from '../../utils/URLUtil'
import { DatabaseInfo, getDatabaseList, } from '../../service/storage'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import ReactJson from 'react-json-view'

const { Column } = Table

interface DatabaseProps {
}

interface DatabaseStatus {
}

@observer
export default class Database extends React.Component<DatabaseProps, DatabaseStatus> {
  @observable listLoading: boolean = false
  @observable createClusterModalVisible: boolean = false
  @observable storageClusterList: DatabaseInfo[] | undefined = undefined

  componentDidMount(): void {
    this.fetchDatabaseList()
  }

  @autobind
  handleCreateDatabase() {
    redirectTo('/setting/database/new')
  }

  async fetchDatabaseList() {
    this.listLoading = true
    this.storageClusterList = await getDatabaseList()
    this.listLoading = false
  }

  render() {
    return (
      <React.Fragment>
        <Card size="small" className="align-right">
          <Button icon="plus-circle" type="primary" onClick={this.handleCreateDatabase}>New Database</Button>
        </Card>

        <Table
          dataSource={this.storageClusterList}
          size="small"
          loading={this.listLoading}
          rowKey="name"
        >
          <Column title="Database Name" dataIndex="name" width={200}/>
          <Column
            title="Config"
            // render={(clusters: any) => JSON.stringify(clusters)}
            render={(config: any) => (
              <ReactJson
                name={false}
                src={config}
                enableClipboard={false}
                displayDataTypes={false}
                displayObjectSize={false}
              />
            )}
          />
        </Table>
      </React.Fragment>
    )
  }
}
/* Database Modal */