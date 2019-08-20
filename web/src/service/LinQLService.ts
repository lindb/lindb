import { GET } from './http'

export async function explain(db: string, ql: string): Promise<any> {
  const url = 'http://adca-infra-etrace-lindb-1.vm.elenet.me:8080/api/search'
  return await GET(url, { db: db, q: ql })
}

// export async function query(params: any): Promise<ResultSet> {
//   const message = 'Query data using ql:'
//   const url = API_URL + '/search'
//   params.cluster = 'local'
//   try {
//     const resp = await GET(url, params)
//     return resp.data
//   } catch (err) {
//     httpCodeHandler(err.response, url, message)
//   }
//   return {}
// }
