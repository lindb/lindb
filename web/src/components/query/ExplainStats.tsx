import { CloudServerOutlined, HddOutlined, ShareAltOutlined } from '@ant-design/icons'
import { Tree } from 'antd'
import * as React from 'react'
import { UnitEnum } from '../../model/Metric'
import { DataFormatter } from '../../utils/DataFormatter'

interface ExplainStatsProps {
    stats?: any
}

interface ExplainStatsStatus {
}

export default class ExplainStats extends React.Component<ExplainStatsProps, ExplainStatsStatus> {
    buildStatsData(stats: any) {
        let root = {
            title: 'Root: ' + DataFormatter.formatter(stats.cost, UnitEnum.Nanoseconds),
            key: 'Root',
            children: [{
                title: <span>Expression: {this.renderCost(stats.expressCost, stats.cost)}</span>,
                key: 'Expression',
            }]
        }
        if (stats.storageNodes) {
            root.children.push(this.buildStorageNodes(stats.cost, stats.storageNodes))
        }
        return [root]
    }
    buildStorageNodes(total: any, storageNodes: any) {
        let children: any = []
        let storageNode = {
            title: (<span>Storage Nodes</span>),
            key: 'Storage-Nodes',
            children: children,
        }
        for (let key of Object.keys(storageNodes)) {
            const storageNodeStats = storageNodes[key]

            let nodeStats = {
                title: (<span>
                    <span className="identifier">{key} </span>
                    [Network: (Cost: {this.renderCost(storageNodeStats.netCost, total)}
                    , Payload: {DataFormatter.formatter(storageNodeStats.netPayload, UnitEnum.Bytes)})]
                </span>),
                key: key,
                icon: <CloudServerOutlined />,
                children: [
                    {
                        title: <span>Storage Execute: {this.renderCost(storageNodeStats.totalCost, total)}</span>,
                        key: key + 'storage-execute',
                    }, {
                        title: <span>Plan Execute: {this.renderCost(storageNodeStats.planCost, total)}</span>,
                        key: key + 'plan-execute',
                    }, {
                        title: <span>Tag Filtering: {this.renderCost(storageNodeStats.tagFilterCost, total)}</span>,
                        key: key + 'tag-filtering',
                    }
                ]
            }
            children.push(nodeStats)
            if (storageNodeStats.shards) {
                nodeStats.children.push(this.buildShardNodes(total, key, storageNodeStats.shards))
            }
        }
        return storageNode
    }
    buildShardNodes(total: any, key: string, shards: any) {
        let children: any = []
        const shardNodes = {
            title: <span>Shards</span>,
            key: key + 'shards',
            children: children,
        }
        for (let shardID of Object.keys(shards)) {
            const shardStats = shards[shardID]
            let nodeStats = {
                title: (<span>
                    <span className="identifier">{shardID} </span>
                    : Series Filtering: {this.renderCost(shardStats.seriesFilterCost, total)}
                    , Num. Of Series: {DataFormatter.formatter(shardStats.numOfSeries, UnitEnum.None)}
                </span>),
                key: key + shardID,
                icon: <ShareAltOutlined />,
                children: [
                    {
                        title: <span>Memory Filtering: {this.renderCost(shardStats.memFilterCost, total)}</span>,
                        key: key + shardID + 'memory fiiltering',
                    },
                    {
                        title: <span>KV Store Filtering: {this.renderCost(shardStats.kvFilterCost, total)}</span>,
                        key: key + shardID + 'kv store fiiltering',
                    },
                    {
                        title: <span>Grouping: {this.renderCost(shardStats.groupingCost, total)}</span>,
                        key: key + shardID + 'Grouping',
                    },
                    {
                        title: <span>Group Build: Count:  {DataFormatter.formatter(shardStats.groupBuildStats.count, UnitEnum.None)}
                            , Min: {this.renderCost(shardStats.groupBuildStats.min, total)}
                            , Max: {this.renderCost(shardStats.groupBuildStats.min, total)}</span>,
                        key: key + shardID + 'group build',
                    },
                ]
            }
            children.push(nodeStats)
            if (shardStats.scanStats) {
                nodeStats.children.push(this.buildLoadStats(total, key + shardID, shardStats.scanStats))
            }
        }
        return shardNodes
    }

    buildLoadStats(total: any, key: string, loadStats: any) {
        let children: any = []
        const loadNodes = {
            title: <span>Load Data</span>,
            key: key + 'load data',
            children: children,
        }
        for (let id of Object.keys(loadStats)) {
            const stats = loadStats[id]
            children.push({
                title: (<span>
                    <span className="identifier">{id} </span>
                    : Count: {DataFormatter.formatter(stats.count, UnitEnum.None)}, Min: {this.renderCost(stats.min, total)}
                    , Max: {this.renderCost(stats.min, total)},
                 </span>),
                key: key + id + 'cost',
                icon: <HddOutlined />,
            })
        }
        return loadNodes
    }

    renderCost(cost: any, total: any) {
        const percent = cost / total
        let clazz = 'green'
        if (percent > 50) {
            clazz = 'danger'
        } else if (percent > 30) {
            clazz = 'warn'
        }
        return (<span className={clazz}>{DataFormatter.formatter(cost, UnitEnum.Nanoseconds)}</span>)
    }
    render() {
        const { stats } = this.props
        if (!stats) {
            return null
        }
        return (
            <React.Fragment>
                <Tree
                    showIcon
                    defaultExpandAll
                    selectable={false}
                    treeData={this.buildStatsData(stats)}
                />
            </React.Fragment>
        )
    }
}