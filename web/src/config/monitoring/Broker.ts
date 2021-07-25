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
            'HTTP p99 Duration',
            'select quantile(0.99) from lindb.broker.http_handle_duration group by path',
            8,
            UnitEnum.Milliseconds,
        ),
        metric(
            'HTTP Count',
            'select HistogramCount from lindb.broker.http_handle_duration group by path',
            8,
            UnitEnum.None,
        ),
        metric(
            'Native ingestion IO',
            'select read_bytes_count from lindb.ingestion.native',
            8,
            UnitEnum.Bytes,
        ),
    ],
    [
        metric(
            'Native ingestion',
            'select unmarshal_metric_count, data_corrupted_count from lindb.ingestion.native',
            8,
            UnitEnum.None,
        ),
        metric(
            'Prometheus ingestion transformed',
            'select transformed_gauges, transformed_counters, transformed_histograms from lindb.ingestion.prometheus',
            8,
            UnitEnum.None,
        ),
        metric(
            'Prometheus ingestion failures',
            'select gzip_data_corrupted, bad_gauges, bad_counters, bad_histograms from lindb.ingestion.prometheus',
            8,
            UnitEnum.None,
        ),
    ],
]
