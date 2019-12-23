import { PATH } from '../config/config'
import { GET, POST } from './APIUtils'

export interface StorageClusterConfig {
  namespace: string
  endpoints: string[]
  dialTimeout?: number
}

export interface StorageClusterInfo {
  name: string
  config: StorageClusterConfig
}

export interface DatabaseCluster {
  name?: string
  cluster?: string
  numOfShard?: number
  replicaFactor?: number
  option: {
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
  }
}

export interface DatabaseInfo {
  name: string
  clusters: DatabaseCluster
}

/* create storage cluster */
export function createStorage(name: string, config: StorageClusterConfig) {
  const url = PATH.storageCluster
  return POST(url, { name, config })
}

/* get storage cluster info */
export function getStorageClusterList() {
  const url = PATH.storageClusterList
  return GET<StorageClusterInfo[]>(url)
}

/* Create Database */
export function createDatabase(payload: DatabaseCluster) {
  const url = PATH.database
  return POST(url, payload)
}

/* get storage cluster info */
export function getDatabaseList() {
  const url = PATH.databaseList
  return GET<DatabaseInfo[]>(url)
}