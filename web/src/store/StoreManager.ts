import ChartEventStore from './ChartEventStore';
import { ChartStore } from "./ChartStore";
import { URLParamStore } from "./URLParamStore";

class StoreManager {
  public ChartEventStore: ChartEventStore = new ChartEventStore();
  public URLParamStore: URLParamStore = new URLParamStore();
  public ChartStore: ChartStore = new ChartStore(this.URLParamStore);
}

export default new StoreManager()