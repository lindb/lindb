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

export const StorageDashboard = [
  // Row
  [
    metric(
        'MemDB total Size',
        'select memdb_total_size from lindb.tsdb.shard group by db, shard',
        8,
        UnitEnum.Bytes,
    ),
    metric(
        'MemDB Number(mutable and immutable)',
        'select memdb_number from lindb.tsdb.shard group by db, shard',
        8,
        UnitEnum.None,
    ),
    metric(
      'Write Data Points',
      'select write_metrics from lindb.tsdb.shard group by db, shard',
      8,
      UnitEnum.None,
    ),
    metric(
        'Bad Metric Data Points',
        'select bad_metrics from lindb.tsdb.shard group by db, shard',
        8,
        UnitEnum.None,
    ),
    metric(
        'Write Data Points Failures',
        'select write_metric_failures from lindb.tsdb.shard group by db, shard',
        8,
        UnitEnum.None,
    ),
    metric(
        'Write Fields',
        'select write_fields from lindb.tsdb.shard group by db, shard',
        8,
        UnitEnum.None,
    ),
    metric(
        'Escaped Fields',
        'select escaped_fields from lindb.tsdb.shard group by db, shard',
        8,
        UnitEnum.None,
    ),
    metric(
        'Metrics Out of Time Range',
        'select metrics_out_of_range from lindb.tsdb.shard group by db, shard',
        8,
        UnitEnum.Milliseconds,
    ),
    metric(
        'Shard Flush Duration',
        'select quantile(0.99) from lindb.tsdb.shard.memdb_flush_duration group by db, shard',
        8,
        UnitEnum.Milliseconds,
    ),
    metric(
      'Generate Metric ID',
      'select gen_metric_ids from lindb.tsdb.metadb group by db',
      8,
      UnitEnum.None,
    ),
    metric(
        'Get Metric ID',
        'select get_metric_ids from lindb.tsdb.metadb group by db',
        8,
        UnitEnum.None,
    ),
    metric(
      'Generate Tag Key ID',
      'select gen_tag_key_ids from lindb.tsdb.metadb group by db',
      8,
      UnitEnum.None,
    ),
    metric(
        'Get Tag Key ID',
        'select get_tag_key_ids from lindb.tsdb.metadb group by db',
        8,
        UnitEnum.None,
    ),
    metric(
      'Generate Field ID',
      'select gen_field_ids from lindb.tsdb.metadb group by db',
      8,
      UnitEnum.None,
    ),
    metric(
        'Get Field ID',
        'select get_field_ids from lindb.tsdb.metadb group by db',
        8,
        UnitEnum.None,
    ),
    metric(
      'Build Inverted Index',
      'select build_inverted_index_counter from lindb.tsdb.indexdb group by db',
      8,
      UnitEnum.None,
    ),
    metric(
        'Build Tag Key ID Failures',
        'select gen_tag_key_id_fails from lindb.tsdb.indexdb  group by db',
        8,
        UnitEnum.None,
    ),
    metric(
      'Generate Tag Value ID Failures',
      'select gen_tag_value_id_fails from lindb.tsdb.indexdb group by db',
      8,
      UnitEnum.None,
    ),
  ],
]
