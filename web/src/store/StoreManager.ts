import ChartEventStore from './ChartEventStore'

class StoreManager {
  public ChartEventStore: ChartEventStore = new ChartEventStore()
}

export default new StoreManager()