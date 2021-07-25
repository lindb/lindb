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

export const ParallelBoard = [
    // Row
    [
        metric(
            'Worker created Counter',
            'select workers_created from lindb.parallel group by pool',
            8,
            UnitEnum.None,
        ),
        metric(
            'Workers killed Counter',
            'select workers_killed from lindb.parallel group by pool',
            8,
            UnitEnum.None,
        ),
        metric(
            'Workers Alive',
            'select workers_alive from lindb.parallel group by pool',
            8,
            UnitEnum.None,
        ),
    ],
    [
        metric(
            'Task consumed Counter',
            'select tasks_consumed from lindb.parallel group by pool',
            8,
            UnitEnum.None,
        ),
        metric(
            'Task waiting average Duration',
            'select tasks_waiting_duration_sum/tasks_consumed as avg_waiting_duration from lindb.parallel group by pool',
            8,
            UnitEnum.Milliseconds,
        ),
        metric(
            'Task executing average Duration',
            'select tasks_executing_duration_sum/tasks_consumed as avg_executing_duration from lindb.parallel group by pool',
            8,
            UnitEnum.Milliseconds,
        ),
    ],
]