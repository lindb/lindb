/**
 * Storage cluster config
 */
export interface StorageConfig {
  namespace: string
  endpoints: string[]
  dialTimeout?: number
}

export interface StorageCluster {
  name: string
  config: StorageConfig
}
