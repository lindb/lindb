/*eslint no-template-curly-in-string: "off"*/
import { Metric, UnitEnum } from 'model/Metric';
import { uuid } from 'uuidv4';

function metric(title: string, ql: string, span: number = 24, unit: UnitEnum = UnitEnum.None, chartType: string = "line"): Metric {
  return {
    span,
    id: uuid(),
    chart: { unit, title, target: { ql, db: '_internal' }, type: chartType },
  }
}

export const KVDashboard = [
  // Row
  [
    metric(
        'KV Builder Added Keys',
        'select add_keys from lindb.kv.table.builder group by node',
        8,
        UnitEnum.None,
    ),
    metric(
        'KV Builder Added Bytes',
        'select add_bytes from lindb.kv.table.builder group by node',
        8,
        UnitEnum.Bytes,
    ),
      metric(
          'KV Builder Added Bad Keys',
          'select bad_keys from lindb.kv.table.builder group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Cache Hits',
          'select cache_hits from lindb.kv.table.cache group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Cache Misses',
          'select cache_misses from lindb.kv.table.cache group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Cache Evicts',
          'select evict_counts from lindb.kv.table.cache group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Cache Close Counts',
          'select close_counts from lindb.kv.table.cache group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Cache Close Errors',
          'select close_errors from lindb.kv.table.cache group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Reader Get Errors',
          'select get_errors from lindb.kv.table.reader group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Reader Get Counts',
          'select get_counts from lindb.kv.table.reader group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Reader Mmap Counts',
          'select mmap_counts from lindb.kv.table.reader group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Reader Unmmap Counts',
          'select unmmap_counts from lindb.kv.table.reader group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Reader Mmap Errors',
          'select mmap_errors from lindb.kv.table.reader group by node',
          8,
          UnitEnum.None,
      ),
      metric(
          'KV Table Reader Unmmap Errors',
          'select unmmap_errors from lindb.kv.table.reader group by node',
          8,
          UnitEnum.None,
      ),
  ],
]
