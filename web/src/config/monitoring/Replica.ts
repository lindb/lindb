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

export const ReplicaDashboard = [
  // Row
  [
    metric(
        'Metric Out Of TimeRange',
        'select metrics_out_of_time_range from lindb.replica.database group by db',
        8,
        UnitEnum.None,
    ),
    metric(
        'Local Replica Count',
        'select replica_count from lindb.replica.local group by db, shard',
        8,
        UnitEnum.None,
    ),
    metric(
        'Local Replica Bytes',
        'select replica_bytes from lindb.replica.local group by db, shard',
        8,
        UnitEnum.Bytes,
    ),
      metric(
          'Local Replica Rows',
          'select replica_rows from lindb.replica.local group by db, shard',
          8,
          UnitEnum.None,
      ),
      metric(
          'Local Replica Sequence',
          'select replica_sequence from lindb.replica.local group by db, shard',
          8,
          UnitEnum.None,
      ),
    metric(
        'Local Replica Max Block Decoded',
        'select max_decoded_block from lindb.replica.local group by db, shard',
        8,
        UnitEnum.Bytes,
    ),
  ],
]
