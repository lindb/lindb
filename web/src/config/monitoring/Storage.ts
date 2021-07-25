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

export const WriteDashboard = [
  // Row
  [
    metric(
      'Write Data Points',
      'select write_metrics from lindb.tsdb.memdb group by db',
      8,
      UnitEnum.None,
    ),
    metric(
        'Write Data Points duration',
        'select write_metric_time_total/write_metric_counter from lindb.tsdb.shard group by db, shard',
        8,
        UnitEnum.None,
    ),
    metric(
        'Write Fields',
        'select write_fields from lindb.tsdb.memdb group by db',
        8,
        UnitEnum.None,
    ),
    metric(
        'Shard Build Index Duration',
        'select quantile(0.99) from lindb.tsdb.shard.build_index_duration group by db, shard',
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
      'Unknown Field Type',
      'select unknown_field_type_counter from lindb.tsdb.memdb group by db',
      8,
      UnitEnum.None,
    ),
    metric(
      'Generate Field ID Fail',
      'select generate_field_id_fails from lindb.tsdb.memdb group by db',
      8,
      UnitEnum.None,
    ),
    metric(
      'Generate Metric ID',
      'select gen_metric_ids from lindb.tsdb.metadb group by db',
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
      'Generate Field ID',
      'select gen_field_ids from lindb.tsdb.metadb group by db',
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
export const IndexDashboard = [
  // Row
  [
    metric(
      'CPU Usage',
      'select 100-gauge*100 from system_cpu_stat  where type="idle" group by node',
      8,
      UnitEnum.Percent,
      "area",
    ),
    metric(
      'Memory Usage',
      'select gauge from system_mem_stat where type in ("used","total") group by node',
      8,
      UnitEnum.Bytes,
    ),
    metric(
      'Disk Usage',
      'select gauge from system_disk_usage where type in ("used","total") group by node',
      8,
      UnitEnum.Bytes,
    ),
  ],
]