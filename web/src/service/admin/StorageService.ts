import { ADMIN_PATH } from 'config/config'
import { StorageCluster } from 'model/admin/Storage'
import { GET, POST } from 'service/APIUtils'

/**
 *  get storage cluster list
 */
export function getStorageList() {
  const url = ADMIN_PATH.storageClusterList
  return GET<StorageCluster[]>(url)
}

/**
 * create storage cluster
 * @param config  storage config
 */
export function createStorageConfig(config: StorageCluster) {
  const url = ADMIN_PATH.storageCluster
  return POST(url, config)
}
