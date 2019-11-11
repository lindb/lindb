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

export interface StorageCluster {
  name: string,
  nodes: any,
  nodeStatus: { total: number, alive: number, suspect: number, dead: number },
  replicaStatus: ReplicaStatus,
  capacity: { usedPercent: number, total: number, used: number },
  databaseStatusList: Array<DatabaseStatus>,
}

export interface DatabaseStatus {
  config: { name: string, numOfShard: number, replicaFactor: number, desc: string },
  replicaStatus: ReplicaStatus,
}

export interface ReplicaStatus {
  total: number,
  underReplicated: number,
  unavailable: number,
}