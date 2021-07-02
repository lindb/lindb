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
      'select counter from mem_write_data_points group by node',
      8,
      UnitEnum.None,
    ),
    metric(
      'Unknown Field Type',
      'select counter from mem_get_unknown_field_type group by node',
      8,
      UnitEnum.None,
    ),
    metric(
      'Generate Field Id Fail',
      'select counter from mem_generate_field_id_fail group by node',
      8,
      UnitEnum.None,
    ),
    metric(
      'Generate Metric Id',
      'select counter from meta_gen_metric_id group by node',
      8,
      UnitEnum.None,
    ),
    metric(
      'Generate Tag Key Id',
      'select counter from meta_gen_tag_key_id group by node',
      8,
      UnitEnum.None,
    ),
    metric(
      'Generate Field Id',
      'select counter from meta_gen_field_id group by node',
      8,
      UnitEnum.None,
    ),
    metric(
      'Build Inverted Index',
      'select counter from build_inverted_index_counter group by node',
      8,
      UnitEnum.None,
    ),
    metric(
      'Generate Tag Value Id',
      'select counter from meta_gen_tag_value_id group by node',
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