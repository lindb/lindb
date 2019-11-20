import { AxiosResponse } from 'axios';
import { QUERY_PATH } from '../config/config';
import { ResultSet } from '../model/Metric';
import { get } from './APIUtils';

export async function query(params: any): Promise<AxiosResponse<ResultSet|undefined>>  {
    const resp = await get<ResultSet>(QUERY_PATH.metric, params);
    return resp;
}