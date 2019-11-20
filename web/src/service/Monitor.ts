import { GET } from './APIUtils'
import { PATH } from '../config/config'
import { NodeList, StorageCluster } from '../model/Monitoring'

export function getMaster() {
  const url = PATH.master
  return GET(url)
}

export function getBrokerCluster() {
  const url = PATH.brokerClusterState
  return GET<NodeList>(url)
}

export function listStorageCluster() {
  const url = PATH.listStorageClusterState
  return GET<Array<StorageCluster>>(url)
}

export function getStorageCluster(name: string) {
  const url = PATH.getStorageClusterState
  return GET<StorageCluster>(url, { name: name })
}