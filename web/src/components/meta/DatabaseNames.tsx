import { Select } from 'antd'
import { autobind } from 'core-decorators'
import { observer } from 'mobx-react'
import * as React from 'react'
import StoreManager from 'store/StoreManager'

const { Option } = Select

interface DatabaseNameSelectProps {
}

interface DatabaseNameSelectState{
}

@observer
export default class DatabaseNameSelect extends React.Component<DatabaseNameSelectProps, DatabaseNameSelectState> {

    componentDidMount(): void {
        this.fetchDatabaseNames()
    }

    async fetchDatabaseNames() {
       await StoreManager.MetadataStore.fetchDatabaseNames("show databases")
    }

    @autobind
    selectDatabase(value: string) {
        StoreManager.URLParamStore.changeURLParams({ db: value })
    }

    render() {
        const databaseNames = StoreManager.MetadataStore.databaseNames
        const db = StoreManager.URLParamStore.getValue("db");
        const value = !db ? "Select Database" : db;
        return (
            <Select value={value} style={{ width: 150 }} onSelect={this.selectDatabase}>
                {databaseNames && databaseNames!.values.map((item:any) => (
                    <Option key={item} value={item}>{item}</Option>
                ))}
            </Select>
        )
    }
}