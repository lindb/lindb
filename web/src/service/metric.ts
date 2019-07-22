import * as LinQLService from './LinQLService'

export async function search(database: string, sql: string) {
  const timeRange = 'time > now()-2h'
  if (!sql || sql.trim() == 'null') {
    return
  }
  const key = database + '$$' + sql
  let linQL = sql

  linQL = sql.replace('${time}', timeRange)
  linQL = linQL.replace('${node}', '*')

  return LinQLService.explain(database, linQL)
}