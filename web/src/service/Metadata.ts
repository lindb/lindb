import { METADATA_PATH } from '../config/config'
import { GET } from './APIUtils'

/* get database names */
export function getDatabaseNames() {
    const url = METADATA_PATH.databaseNames
    return GET<string[]>(url)
}
