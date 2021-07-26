import { Select } from 'antd'
import { autobind } from 'core-decorators'
import { observer } from 'mobx-react'
import * as React from 'react'
import StoreManager from 'store/StoreManager'

const { Option } = Select

interface TagValuesSelectProps {
    tagKey: string
    metric: string
    watch?: string[]
    mode?: 'multiple' | 'tags'
}

interface TagValuesSelectState {
    values: string[]
}

@observer
export default class TagValuesSelect extends React.Component<TagValuesSelectProps, TagValuesSelectState> {
    constructor(props: TagValuesSelectProps) {
        super(props)
        this.state = { values: [] }
    }

    @autobind
    async fetchDatabaseNames() {
        const { metric, tagKey, watch } = this.props
        let showTagValuesSQL = `show tag values from '${metric}' with key='${tagKey}'`
        if (watch) {
            const tags: string[] = []
            watch.forEach(item => {
                const watchValues = StoreManager.URLParamStore.getValues(item)
                const values: string[] = []
                watchValues.forEach(v => {
                    values.push(`'${v}'`)
                })
                if (values.length > 0) {
                    const tv = values.join(",")
                    tags.push(`'${item}' in (${tv})`)
                }
            })

            if (tags.length > 0) {
                showTagValuesSQL += " where " + tags.join(" and ")
            }
        }
        const tagValues = await StoreManager.MetadataStore.getTagValues(showTagValuesSQL)
        if (tagValues!.values) {
            this.setState({ values: tagValues!.values })
        }
    }

    @autobind
    selectTagValue(value: string) {
        const { tagKey } = this.props
        StoreManager.URLParamStore.changeURLParams({ [tagKey]: value })
    }

    render() {
        const { tagKey, mode } = this.props
        const { values } = this.state
        let selected: any
        if (mode === 'multiple' || mode === 'tags') {
            selected = StoreManager.URLParamStore.getValues(tagKey)
        } else {
            selected = StoreManager.URLParamStore.getValue(tagKey);
        }
        return (
            <Select defaultValue={selected} placeholder="Please select" style={{ minWidth: 150 }} mode={mode} onChange={this.selectTagValue} onFocus={this.fetchDatabaseNames}>
                {values && values.map(item => (
                    <Option key={item} value={item}>{item}</Option>
                ))}
            </Select>
        )
    }
}