import { DatabaseStore } from 'store/admin/DatabaseStore';
import { StorageStore } from 'store/admin/StorageStore';
import { BreadcrumbStore } from 'store/BreadcrumbStore';
import ChartEventStore from 'store/ChartEventStore';
import { ChartStore } from 'store/ChartStore';
import { MetadataStore } from 'store/metadata/MetadataStore';
import { URLParamStore } from 'store/URLParamStore';

class StoreManager {
  public ChartEventStore: ChartEventStore = new ChartEventStore();
  public URLParamStore: URLParamStore = new URLParamStore();
  public ChartStore: ChartStore = new ChartStore(this.URLParamStore);
  public BreadcrumbStore: BreadcrumbStore = new BreadcrumbStore();
  public StorageStore: StorageStore = new StorageStore();
  public DatabaseStore: DatabaseStore = new DatabaseStore()
  public MetadataStore: MetadataStore = new MetadataStore()
}

export default new StoreManager()