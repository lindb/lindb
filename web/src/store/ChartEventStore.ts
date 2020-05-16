import { observable } from 'mobx'
import { ChartTooltipData } from 'model/Metric'

export default class ChartEventStore {
  @observable public hiddenSeries: Map<string, number> = new Map()
  @observable public tooltipData: ChartTooltipData | null = null

  public changeSeriesByClick(uuid: string) {
    this.hiddenSeries.set(uuid, (new Date()).getTime())
  }

  public deleteSeriesById(uuid: string) {
    this.hiddenSeries.delete(uuid)
  }

  public setTooltipData(data: ChartTooltipData | null) {
    this.tooltipData = data
  }
}