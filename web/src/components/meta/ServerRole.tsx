import { Select } from 'antd'
import { autobind } from 'core-decorators'
import { observer } from 'mobx-react'
import * as React from 'react'
import StoreManager from 'store/StoreManager'

const { Option } = Select

interface ServerRoleSelectProps {
}

interface ServerRoleSelectState {
}

@observer
export default class ServerRoleSelect extends React.Component<ServerRoleSelectProps, ServerRoleSelectState> {

    @autobind
    selectRole(value: string) {
        StoreManager.URLParamStore.changeURLParams({ role: value })
    }

    render() {
        const role = StoreManager.URLParamStore.getValue("role");
        const value = !role ? "broker" : role;
        return (
            <Select value={value} style={{ width: 150 }} onSelect={this.selectRole}>
                <Option value="broker">Broker</Option>
                <Option value="storage">Storage</Option>
            </Select>
        )
    }
}