/*
Licensed to LinDB under one or more contributor
license agreements. See the NOTICE file distributed with
this work for additional information regarding copyright
ownership. LinDB licenses this file to you under
the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
import { CapacityView } from "@src/components";
import { useStateMetric } from "@src/hooks";
import * as _ from "lodash-es";
import React, { useEffect, useState } from "react";

interface DiskUsageViewProps {
  sql: string;
}

export default function DiskUsageView(props: DiskUsageViewProps) {
  const { sql } = props;
  const { stateMetric } = useStateMetric(sql);
  const [stats, setStats] = useState({
    total: 0,
    used: 0,
    free: 0,
  });

  useEffect(() => {
    const getValue = (fields: any[], name: string): number => {
      return _.get(_.find(fields, { name: name }), "value", 0);
    };
    const stats = {
      total: 0,
      used: 0,
      free: 0,
    };
    _.forEach(stateMetric, (seriesList: any, _key: string) => {
      _.forEach(seriesList, (series: any) => {
        const fields = _.get(series, "fields", []);
        stats.total += getValue(fields, "total");
        stats.used += getValue(fields, "used");
        stats.free += getValue(fields, "free");
      });
    });
    setStats(stats);
  }, [stateMetric]);

  return (
    <CapacityView
      percent={stats.total ? (stats.used * 100) / stats.total : 0}
      total={stats.total}
      free={stats.free}
      used={stats.used}
    />
  );
}
