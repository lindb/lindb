import { GET } from './http'
import { PATH } from '../config/config'
import { NodeList } from '../model/Monitoring'

export function getMaster() {
  const url = PATH.master
  return GET(url)
}

export function getBrokerCluster() {
  const url = PATH.brokerClusterState
  return GET<NodeList>(url)
}