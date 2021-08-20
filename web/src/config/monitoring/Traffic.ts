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
    metric(
        'TCP Server Accepted Connections',
        'select accept_conns from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.None,
    ),
    metric(
        'TCP Server Accepted Errors',
        'select accept_errors from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.None,
    ),
    metric(
        'TCP Server Closed Connections',
        'select close_conns from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.None,
    ),
    metric(
        'TCP Server Closed Errors',
        'select close_errors from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.None,
    ),
    metric(
        'TCP Server Connections',
        'select conns_num from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.None,
    ),
    metric(
        'TCP Server Read Counts',
        'select read_count from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.None,
    ),
    metric(
        'TCP Server Write Counts',
        'select write_count from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.None,
    ),
    metric(
        'TCP Server Read Bytes',
        'select read_bytes from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.Bytes,
    ),
    metric(
        'TCP Server Write Bytes',
        'select write_bytes from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.Bytes,
    ),
    metric(
        'TCP Server Read Errors',
        'select read_errors from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.None,
    ),
    metric(
        'TCP Server Write Errors',
        'select write_errors from lindb.traffic.tcp group by node, addr',
        8,
        UnitEnum.None,
    ),
  ],
]
