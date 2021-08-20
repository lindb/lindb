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

export const TrafficDashboard = [
  // Row
  [
    metric(
        'GRPC Server Message Received',
        'select msg_received from lindb.traffic.grpc_server group by node, grpc_type, grpc_service, grpc_method',
        8,
        UnitEnum.None,
    ),
    metric(
        'GRPC Server Message Sent',
        'select msg_sent from lindb.traffic.grpc_server group by node, grpc_type, grpc_service, grpc_method',
        8,
        UnitEnum.None,
    ),
    metric(
        'GRPC Client Message Received',
        'select msg_received from lindb.traffic.grpc_client group by node, grpc_type, grpc_service, grpc_method',
        8,
        UnitEnum.None,
    ),
    metric(
        'GRPC Client Message Sent',
        'select msg_sent from lindb.traffic.grpc_client group by node, grpc_type, grpc_service, grpc_method',
        8,
        UnitEnum.None,
    ),
    metric(
        'GRPC Client Message Received P99 Duration',
        'select quantile(0.99) from lindb.traffic.grpc_client.msg_received_duration group by node, grpc_type, grpc_service, grpc_method',
        8,
        UnitEnum.Milliseconds,
    ),
    metric(
        'GRPC Client Message Sent P99 Duration',
        'select quantile(0.99) from lindb.traffic.grpc_client.msg_sent_duration group by node, grpc_type, grpc_service, grpc_method',
        8,
        UnitEnum.Milliseconds,
    ),
  ],
]
