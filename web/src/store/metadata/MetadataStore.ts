import { observable } from 'mobx'
import { Metadata } from 'model/meta/Metadata'
import { fetchMetadata } from 'service/metadata/MetadataService'

export class MetadataStore {
    @observable public loading: boolean = false
    @observable public databaseNames: Metadata | undefined
    @observable public metadata: Metadata | undefined

    public async fetchMetadata(db: string, sql: string) {
        this.loading = true
        this.metadata = await fetchMetadata({ db: db, sql: sql })
        this.loading = false
    }

    public async fetchDatabaseNames(sql: string) {
        this.loading = true
        this.databaseNames = await fetchMetadata({ sql: sql })
        this.loading = false
    }
}