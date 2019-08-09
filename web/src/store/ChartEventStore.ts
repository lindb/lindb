import { observable } from 'mobx'
import { ChartTooltipData } from '../model/Metric'

export default class ChartEventStore {
  @observable public hiddenSeries: Map<string, boolean[]> = new Map()
  @observable public tooltipData: ChartTooltipData | null = null

  public changeSeriesByClick(uuid: string, seriesStatus: boolean[]) {
    this.hiddenSeries.set(uuid, seriesStatus)
  }

  public deleteSeriesById(uuid: string) {
    this.hiddenSeries.delete(uuid)
  }

  public setTooltipData(data: ChartTooltipData | null) {
    this.tooltipData = data
  }
}