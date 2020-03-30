export interface DatabaseConfig {
    name?: string
    cluster?: string
    numOfShard?: number
    replicaFactor?: number
    option?: {
        interval?: string
        timeWindow?: number
        autoCreateNS?: boolean
        behind?: string
        ahead?: string
        index: {
            timeThreshold?: number
            sizeThreshold?: number
        },
        data: {
            timeThreshold?: number
            sizeThreshold?: number
        }
    },
    desc?: string
}

export const DefaultDatabaseConfig:DatabaseConfig = {
    numOfShard: 1,
    replicaFactor: 1,
}