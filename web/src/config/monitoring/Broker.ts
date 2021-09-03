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
            'HTTP P99 Duration',
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
            'Native Ingestion IO',
            'select read_bytes from lindb.ingestion.native',
            8,
            UnitEnum.Bytes,
        ),
    ],
    [
        metric(
            'Native Ingestion',
            'select ingested_metrics, data_corrupted_count from lindb.ingestion.native',
            8,
            UnitEnum.None,
        ),
        metric(
            'InfluxDB Ingestion IO',
            'select read_bytes from lindb.ingestion.influx',
            8,
            UnitEnum.Bytes,
        ),
        metric(
            'InfluxDB Ingestion Failures',
            'select dropped_metrics, data_corrupted_count from lindb.ingestion.influx',
            8,
            UnitEnum.None,
        ),
    ],
    [
        metric(
            'InfluxDB Ingestion Count',
            'select ingested_metrics, ingested_fields from lindb.ingestion.influx',
            8,
            UnitEnum.None,
        ),
        metric(
            'Prometheus Ingestion Transformed',
            'select transformed_gauges, transformed_counters, transformed_histograms from lindb.ingestion.prometheus',
            8,
            UnitEnum.None,
        ),
        metric(
            'Prometheus Ingestion Failures',
            'select gzip_data_corrupted, bad_gauges, bad_counters, bad_histograms from lindb.ingestion.prometheus',
            8,
            UnitEnum.None,
        ),
    ]
]
