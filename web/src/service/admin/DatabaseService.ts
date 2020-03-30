import { ADMIN_PATH } from '../../config/config'
import { DatabaseConfig } from '../../model/admin/Database'
import { GET, POST } from '../APIUtils'

/**
 *  create database config
 * @param config database config
 */
export function createDatabase(config: DatabaseConfig) {
    const url = ADMIN_PATH.database
    return POST(url, config)
}

/**
 * get all database list
 */
export function getDatabaseList() {
    const url = ADMIN_PATH.databaseList
    return GET<DatabaseConfig[]>(url)
}