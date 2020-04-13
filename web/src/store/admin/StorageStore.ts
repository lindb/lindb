import { observable } from 'mobx'
import { StorageCluster } from 'model/admin/Storage'
import { createStorageConfig, getStorageList } from 'service/admin/StorageService'

export class StorageStore {
    @observable public loading: boolean = false
    @observable public storageList: StorageCluster[] | undefined = []

    public async fetchStorageList() {
        this.loading = true
        this.storageList = await getStorageList()
        this.loading = false
    }

    public async createStorageConfig(config: StorageCluster) {
        await createStorageConfig(config)
    }
}