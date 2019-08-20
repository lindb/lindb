import { GET } from './http'
import { PATH } from '../config/config'
import { BrokersList } from '../model/service'

export function getMaster() {
  const url = PATH.master
  return GET(url)
}

export function getBrokersList() {
  const url = PATH.brokers
  return GET<BrokersList>(url)
}