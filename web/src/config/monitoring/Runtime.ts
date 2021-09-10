/*eslint no-template-curly-in-string: "off"*/
import {Metric, UnitEnum} from 'model/Metric';
import {uuid} from 'uuidv4';

function metric(title: string, ql: string, span: number = 24, unit: UnitEnum = UnitEnum.None): Metric {
    return {
        span,
        id: uuid(),
        chart: {unit, title, target: {ql, db: '_internal'}},
    }
}

export const RuntimeStorageBoard = [
    // Row
    [
        metric(
            'Sys (number of heap bytes obtained from system)',
            'select heap_sys_bytes from lindb.runtime.mem',
            8,
            UnitEnum.Bytes,
        ),
        metric(
            'Frees (number of frees)',
            'select frees_total from lindb.runtime.mem',
            8,
            UnitEnum.Bytes,
        ),
        metric(
            'Total Alloc (bytes allocated even if freed)',
            'select alloc_bytes_total from lindb.runtime.mem',
            8,
            UnitEnum.Bytes,
        ),
    ],
    [
        metric(
            'HeapAlloc (bytes allocated and not yet freed)',
            'select heap_alloc_bytes from lindb.runtime.mem',
            8,
            UnitEnum.Bytes,
        ),
        metric(
            'Heap Objects (total number of allocated objects)',
            'select heap_objects from lindb.runtime.mem',
            8,
            UnitEnum.None,
        ),
        metric(
            'HeapInUsed (bytes in non-idle span)',
            'select heap_inuse_bytes from lindb.runtime.mem',
            8,
            UnitEnum.Bytes,
        ),
    ],
    [
        metric(
            'Number of goroutines',
            'select go_goroutines from lindb.runtime',
            8,
            UnitEnum.None,
        ),
        metric(
            'Number of Threads',
            'select go_threads from lindb.runtime',
            8,
            UnitEnum.None,
        ),
        metric(
            'Next GC Bytes',
            'select next_gc_bytes from lindb.runtime.mem',
            8,
            UnitEnum.Bytes,
        ),
    ],
    [
        metric(
            'GC CPU Fraction',
            'select gc_cpu_fraction from lindb.runtime.mem',
            8,
            UnitEnum.None,
        ),
        metric(
            'Lookups(number of pointer lookups)',
            'select lookups_total from lindb.runtime.mem',
            8,
            UnitEnum.None,
        ),
    ]
]