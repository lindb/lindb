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

export const SystemBoardForRole = [
  // Row
  [
    metric(
      'CPU Usage',
      'select 100-gauge*100 from system_cpu_stat where type="idle" group by role, node',
      8,
      UnitEnum.Percent,
      "area",
    ),
    metric(
      'Memory Usage',
      'select gauge from system_mem_stat where type="used_percent" group by role, node',
      8,
      UnitEnum.Percent,
    ),
    metric(
      'Disk Usage',
      'select gauge from system_disk_usage_stat where type="used_percent" group by role, node',
      8,
      UnitEnum.Percent,
    ),
    metric(
        'Net In Speed',
        'select counter/1024 from net_stat where type="bytes_recv" group by role,node',
        8,
        UnitEnum.KBytesPerSecond,
    ),
    metric(
        'Net Out Speed',
        'select counter/1024 from net_stat where type="bytes_sent" group by role,node',
        8,
        UnitEnum.KBytesPerSecond,
    )
  ],
]
export const SystemBoardForNode = [
  // Row
  [
    metric(
      'CPU Usage',
      'select gauge*100 from system_cpu_stat group by node',
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
      'select gauge from system_disk_usage_stat where type in ("used","total") group by node',
      8,
      UnitEnum.Bytes,
    ),
  ],
]