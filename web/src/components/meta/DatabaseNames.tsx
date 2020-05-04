import { Select } from 'antd'
import { autobind } from 'core-decorators'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import * as React from 'react'
import { getDatabaseNames } from 'service/Metadata'
import StoreManager from 'store/StoreManager'

const { Option } = Select

interface DatabaseNameSelectProps {
}

interface DatabaseNameSelectState{
}

@observer
export default class DatabaseNameSelect extends React.Component<DatabaseNameSelectProps, DatabaseNameSelectState> {
    @observable databaseNames: string[] | undefined = undefined

    componentDidMount(): void {
        //TODO need use store
        this.fetchDatabaseNames()
    }

    async fetchDatabaseNames() {
        this.databaseNames = await getDatabaseNames()
    }

    @autobind
    selectDatabase(value: string) {
        StoreManager.URLParamStore.changeURLParams({ db: value })
    }

    render() {
        const db = StoreManager.URLParamStore.getValue("db");
        const value = !db ? "Select Database" : db;
        return (
            <Select value={value} style={{ width: 150 }} onSelect={this.selectDatabase}>
                {this.databaseNames && this.databaseNames.map(item => (
                    <Option key={item} value={item}>{item}</Option>
                ))}
            </Select>
        )
    }
}