export const PREFIXCLS = 'lindb-base'
export const LOCALSTORAGE_TOKEN = 'TOKEN'
export const TIME_FORMAT: string = 'YYYY-MM-DD HH:mm:ss'

export const SPACING = 10

/* API URL */
// export const API_URL = 'http://localhost:9000/api/v1'
export const API_URL = '/api/v1'

export const PATH = {
  login: '/login',
  master: '/cluster/master',
  brokerClusterState: '/broker/cluster/state',
  listStorageClusterState: '/storage/cluster/state/list',
  getStorageClusterState: '/storage/cluster/state',
  storageStateList: '/storage/state/list',
}

export const ADMIN_PATH = {
  storageCluster: '/storage/cluster',
  storageClusterList: '/storage/cluster/list',
  database: '/database',
  databaseList: '/database/list',
}

export const QUERY_PATH = {
  metric: '/query/metric',
  metadata: '/query/metadata'
}