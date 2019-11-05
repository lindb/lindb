export type NodeList = Node[]

export interface Node {
  node: {
    node: { ip: string, port: number, hostName: string },
    onlineTime: number,
  },
  system: {
    cpus: number,
    memoryStat: { usedPercent: number, total: number, used: number },
    diskStat: { usedPercent: number, total: number, used: number },
  }
}
