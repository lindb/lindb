/*eslint no-template-curly-in-string: "off"*/
import { Metric, UnitEnum } from 'model/Metric';
import { uuid } from 'uuidv4';

function metric(title: string, ql: string, span: number = 24, unit: UnitEnum = UnitEnum.None): Metric {
  return {
    span,
    id: uuid(),
    chart: { unit, title, target: { ql, db: '_internal' } },
  }
}

export const RuntimeStorageBoard = [
  // Row
  [
    metric(
      'Sys (number of heap bytes obtained from system)',
      'select gauge from go_memstats_heap_sys_bytes',
      8,
      UnitEnum.Bytes,
    ),
    metric(
      'Frees (number of frees)',
      'select counter from go_memstats_frees_total',
      8,
      UnitEnum.Bytes,
    ),
    metric(
      'Total Alloc (bytes allocated even if freed)',
      'select counter from go_memstats_alloc_bytes_total',
      8,
      UnitEnum.Bytes,
    ),
  ],
  [
    metric(
      'HeapAlloc (bytes allocated and not yet freed)',
      'select gauge from go_memstats_heap_alloc_bytes',
      8,
      UnitEnum.Bytes,
    ),
    metric(
      'Heap Objects (total number of allocated objects)',
      'select gauge from go_memstats_heap_objects',
      8,
      UnitEnum.None,
    ),
    metric(
      'HeapInused (bytes in non-idle span)',
      'select gauge from go_memstats_heap_inuse_bytes',
      8,
      UnitEnum.Bytes,
    ),
  ],
  [
    metric(
      'Number of goroutines',
      'select gauge from go_goroutines',
      8,
      UnitEnum.None,
    ),
    metric(
      'GC invocation duration',
      'select summary from go_gc_duration_seconds',
      8,
      UnitEnum.Seconds,
    ),
    metric(
      'GC invocation count',
      'select count(summary) from go_gc_duration_seconds',
      8,
      UnitEnum.None,
    ),
  ],
  [
    metric(
      'Lookups(number of pointer lookups)',
      'select counter from go_memstats_lookups_total',
      8,
      UnitEnum.None,
    ),
  ],
]