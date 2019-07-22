import { Metric, UnitEnum } from '../../model/Metric'

let SEQ = 910000000

function metric(title: string, ql: string, span: number = 24, unit: UnitEnum = UnitEnum.None): Metric {
  return {
    span,
    chart: { id: SEQ++, ql, unit, title },
  }
}

export const SystemStorageBoard = [
  // Row
  [
    metric(
      'System CPU',
      'select max(gauge)*100.0 from cpu.load.system where ${time} and node=\'${node}\' group by node',
      12,
      UnitEnum.Percent,
    ),
    metric(
      'LinDB CPU',
      'select max(gauge)*100.0 from cpu.load.process where ${time} and node=\'${node}\' group by node',
      12,
      UnitEnum.Percent,
    ),
  ],
  // Row
  [
    metric(
      'Heap Memory Total',
      'select gauge from memory.heap where ${time} and type in (\'max\',\'totalUsed\')and node=\'${node}\' group by node,type',
      8,
      UnitEnum.Bytes,
    ),
    metric(
      'Heap Memory Table',
      'select gauge from memory.heap where ${time} and type=\'tableUsed\' and node=\'${node}\' group by node',
      8,
      UnitEnum.Bytes,
    ),
    metric(
      'Direct Memory',
      'select gauge from memory.direct where ${time} and node=\'${node}\' group by node,type',
      8,
      UnitEnum.Bytes,
    ),
  ],
  // Row
  [
    metric(
      'Chunk Alloc Count',
      'select count from direct.alloc.chunk where ${time} and node=\'${node}\' group by node',
      12,
    ),
    metric(
      'Chunk Free Count',
      'select count from direct.free.chunk where ${time} and node=\'${node}\' group by node',
      12,
    ),
  ],
  // Row
  [
    metric(
      'Storage Capacity',
      'select gauge from disk.usage where ${time} and node=\'${node}\' group by node,type,time(10m)',
      12,
      UnitEnum.Bytes,
    ),
    metric(
      'File Opens',
      'select max(gauge) from os.file.descriptor where ${time} and type=\'open\' and node=\'${node}\' group by node',
      12,
    ),
  ],
]