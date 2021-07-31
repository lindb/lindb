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

export const QueryBoard = [
    // Row
    [
        metric(
            'Broker Created Tasks',
            'select created_tasks from lindb.broker.query group by node',
            8,
            UnitEnum.None,
        ),
        metric(
            'Broker Alive Tasks',
            'select alive_tasks from lindb.broker.query group by node',
            8,
            UnitEnum.None,
        ),
        metric(
            'Broker Sent Task Requests',
            'select sent_requests from lindb.broker.query group by node',
            8,
            UnitEnum.None,
        ),
    ],
    [
        metric(
            'Broker Sent Tasks Responses',
            'select sent_responses from lindb.broker.query group by node',
            8,
            UnitEnum.None,
        ),
        metric(
            'Broker Sent Task Request Failures',
            'select sent_requests_failures from lindb.broker.query group by node',
            8,
            UnitEnum.None,
        ),
        metric(
            'Broker Sent Task Response Failures',
            'select sent_responses_failures from lindb.broker.query group by node',
            8,
            UnitEnum.None,
        ),
    ],
    [
        metric(
            'Storage Metric Queries Counter',
            'select metric_queries from lindb.storage.query group by node',
            8,
            UnitEnum.None,
        ),
        metric(
            'Storage Meta Task Queries Counter',
            'select meta_queries from lindb.storage.query group by node',
            8,
            UnitEnum.None,
        ),
        metric(
            'Storage Omitted Response',
            'select omitted_responses from lindb.storage.query group by node',
            8,
            UnitEnum.None,
        ),
    ],
]