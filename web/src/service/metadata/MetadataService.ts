import { QUERY_PATH } from 'config/config';
import { Metadata } from 'model/meta/Metadata';
import { GET } from 'service/APIUtils';

/**
 * fetch metadata by sql 
 * @param params 
 */
export function fetchMetadata(params: any) {
    return GET<Metadata>(QUERY_PATH.metadata, params)
}