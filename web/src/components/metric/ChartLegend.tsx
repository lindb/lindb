import * as React from 'react'
import classNames from 'classnames'
import { autobind } from 'core-decorators'
import { reaction } from 'mobx';
import StoreManager from 'store/StoreManager';
import { ChartStatusEnum } from 'model/Chart';

interface ChartLegendProps {
  uuid: string
  onLegendClick?: () => void
}

interface ChartLegendState {
  series: any
}

export default class ChartLegend extends React.Component<ChartLegendProps, ChartLegendState> {
  private chartLegendCls = 'lindb-chart-legend'
  disposers: any[]

  constructor(props: ChartLegendProps) {
    super(props)

    this.state = {
      series: {},
    }

    this.disposers = [
      reaction(
        () => StoreManager.ChartStore.chartStatusMap.get(props.uuid),
        chartStatus => {
          console.log("init", props.uuid)
          if (chartStatus!.status !== ChartStatusEnum.Loaded) {
            return
          }
          const series = StoreManager.ChartStore.seriesCache.get(props.uuid) || {};

          this.setState({ series: series })
        }
      )
    ]
  }

  componentWillUnmount(): void {
    this.disposers.map(handle => handle())
  }

  @autobind
  handleLegendItemClick(item: any) {
    const { onLegendClick, uuid } = this.props
    let selected = StoreManager.ChartStore.selectedSeries.get(uuid)
    if (!selected) {
      selected = new Map()
    }
    const label = item.label
    if (selected.has(label)) {
      selected.delete(label)
    } else {
      selected.clear()
      selected.set(label, label)
    }

    this.forceUpdate()

    StoreManager.ChartStore.selectedSeries.set(uuid, selected)
    onLegendClick && onLegendClick()
  }

  render() {
    const { uuid } = this.props
    const { series } = this.state
    const cls = this.chartLegendCls
    if (!series) {
      return null
    }
    const selected = StoreManager.ChartStore.selectedSeries.get(uuid)

    return (
      <div className={cls}>
        {/* render Legend Item */}
        {series.datasets && series.datasets.map((item: any, index: number) => {
          const { label, borderColor } = item
          return (
            <span
              key={label + borderColor}
              className={classNames(`${cls}__item`, { hidden: selected && selected.size > 0 && !selected.has(label) })}
              onClick={() => this.handleLegendItemClick(item)}
              title={label}
            >
              <i className={`${cls}__item__icon`} style={{ backgroundColor: borderColor }} />{label}
            </span>
          )
        })}
      </div>
    )
  }
}