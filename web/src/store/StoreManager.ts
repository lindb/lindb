import { DatabaseStore } from './admin/DatabaseStore';
import { StorageStore } from './admin/StorageStore';
import { BreadcrumbStore } from './BreadcrumbStore';
import ChartEventStore from './ChartEventStore';
import { ChartStore } from './ChartStore';
import { URLParamStore } from './URLParamStore';

class StoreManager {
  public ChartEventStore: ChartEventStore = new ChartEventStore();
  public URLParamStore: URLParamStore = new URLParamStore();
  public ChartStore: ChartStore = new ChartStore(this.URLParamStore);
  public BreadcrumbStore: BreadcrumbStore = new BreadcrumbStore();
  public StorageStore: StorageStore = new StorageStore();
  public DatabaseStore: DatabaseStore = new DatabaseStore()
}

export default new StoreManager()