/*eslint no-template-curly-in-string: "off"*/
import {Metric, UnitEnum} from 'model/Metric';
import {uuid} from 'uuidv4';

function metric(title: string, ql: string, span: number = 24, unit: UnitEnum = UnitEnum.None, chartType: string = "line"): Metric {
    return {
        span,
        id: uuid(),
        chart: {unit, title, target: {ql, db: '_internal'}, type: chartType},
    }
}

export const BrokerDashboard = [
    // Row
    [
        metric(
            'HTTP Duration',
            'select counter from http_handle_duration group by node',
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