import { API_URL, Get } from './index'
import { httpCodeHandler } from '../utils/notification'
import { Group, Result, ResultSet } from '../model/Metric'

export async function explain(db: string, ql: string): Promise<any> {
  const url = API_URL + '/api/search'
  return await Get(url, { db: db, q: ql })
}

export async function query(params: any): Promise<ResultSet> {
  const message = 'Query data using ql:'
  const url = API_URL + '/search'
  params.cluster = 'local'
  try {
    const resp = await Get(url, params)
    return resp.data
  } catch (err) {
    httpCodeHandler(err.response, url, message)
  }
  return {}
}
