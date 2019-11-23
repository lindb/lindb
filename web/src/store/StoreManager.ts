import ChartEventStore from './ChartEventStore';
import { ChartStore } from './ChartStore';
import { URLParamStore } from './URLParamStore';
import { BreadcrumbStore } from './BreadcrumbStore';

class StoreManager {
  public ChartEventStore: ChartEventStore = new ChartEventStore();
  public URLParamStore: URLParamStore = new URLParamStore();
  public ChartStore: ChartStore = new ChartStore(this.URLParamStore);
  public BreadcrumbStore: BreadcrumbStore = new BreadcrumbStore();
}

export default new StoreManager()