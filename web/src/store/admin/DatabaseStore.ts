import { observable } from 'mobx'
import { DatabaseConfig } from 'model/admin/Database'
import { createDatabase, getDatabaseList } from 'service/admin/DatabaseService'

export class DatabaseStore {
    @observable public loading: boolean = false
    @observable public databaseList: DatabaseConfig[] | undefined = []

    public async fectchDatabaseList() {
        this.loading = true
        this.databaseList = await getDatabaseList()
        this.loading = false
    }

    public async createDatabase(config: DatabaseConfig) {
        await createDatabase(config)
    }
}