export type BrokersList = BrokerInfo[]

export interface BrokerInfo {
  node: { ip: string, port: number },
  onlineTime: number
}