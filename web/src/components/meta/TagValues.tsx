import { Select } from 'antd'
import { autobind } from 'core-decorators'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import * as React from 'react'
import StoreManager from 'store/StoreManager'

const { Option } = Select

interface TagValuesSelectProps {
    tagKey: string
}

interface TagValuesSelectState {
}

@observer
export default class TagValuesSelect extends React.Component<TagValuesSelectProps, TagValuesSelectState> {
    constructor(props: TagValuesSelectProps) {
        super(props)
    }

    @autobind
    selectTagValue(value: string) {
        StoreManager.URLParamStore.changeURLParams({ db: value })
    }

    render() {
        const db = StoreManager.URLParamStore.getValue("db");
        const value = !db ? "Select Database" : db;
        return (
            <Select value={value} style={{ width: 150 }} onSelect={this.selectTagValue}>
                <Option key="test" value="test">test</Option>
            </Select>
        )
    }
}