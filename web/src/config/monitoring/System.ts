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

export const SystemBoardForRole = [
    // Row
    [
        metric(
            'CPU Usage',
            'select 100-idle*100 as used_percent from lindb.monitor.system.cpu_stat',
            8,
            UnitEnum.Percent,
            "area",
        ),
        metric(
            'Memory Usage',
            'select used_percent from lindb.monitor.system.mem_stat',
            8,
            UnitEnum.Percent,
        ),
        metric(
            'Disk Usage',
            'select used_percent from lindb.monitor.system.disk_usage_stats',
            8,
            UnitEnum.Percent,
        ),
        metric(
            'Net In Speed',
            'select bytes_recv/1024 from lindb.monitor.system.net_stat',
            8,
            UnitEnum.KBytesPerSecond,
        ),
        metric(
            'Net Out Speed',
            'select bytes_sent/1024 from lindb.monitor.system.net_stat',
            8,
            UnitEnum.KBytesPerSecond,
        )
    ],
]
