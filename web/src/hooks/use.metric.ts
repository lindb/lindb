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

import { QueryStatement, ResultSet, Query } from "@src/models";
import { URLStore } from "@src/stores";
import { useQuery } from "@tanstack/react-query";
import { reaction } from "mobx";
import { useEffect } from "react";
import * as _ from "lodash-es";
import { ExecService } from "@src/services";
import { URLKit } from "@src/utils";

const buildSQL = (ql: string | null, cascade: string[]) => {
  if (!ql) {
    return "";
  }
  const params = URLStore.getParams();
  const tags: string[] = URLKit.getTagConditions(params, cascade);
  const where: string[] = [];
  const timeRange = URLStore.getTimeRange();
  if (timeRange && ql.indexOf(timeRange) < 0) {
    where.push(timeRange);
  }
  if (tags.length > 0) {
    where.push(`${tags.join(" and ")}`);
  }
  const whereAt = ql.indexOf("where ");
  const whereClause = where.join(" and ");
  if (whereAt < 0) {
    if (where.length === 0) {
      return ql;
    } else {
      const groupByAt = ql.indexOf("group by");
      if (groupByAt < 0) {
        return `${ql} where ${whereClause}`;
      }
      return `${ql.slice(0, groupByAt)} where ${whereClause} ${ql.slice(
        groupByAt,
        ql.length
      )}`;
    }
  }
  if (whereClause.length > 0) {
    // has tag filter
    return `${ql.slice(0, whereAt + 6)}  ${whereClause} and ${ql.slice(
      whereAt + 6,
      ql.length
    )}`;
  }
  // no tag filter
  return `${ql.slice(0, whereAt + 6)}  ${whereClause} ${ql.slice(
    whereAt + 6,
    ql.length
  )}`;
};

export function useMetric(
  queries: Query[],
  options?: { disableBind?: boolean }
) {
  const {
    isInitialLoading,
    isLoading,
    isFetching,
    isError,
    data,
    error,
    refetch,
  } = useQuery(
    ["search_metric_data", queries],
    async () => {
      const requests: any[] = [];
      (queries || []).forEach((query: Query) => {
        const db = _.get(query, "db", "");
        const sql = _.get(query, "sql", "");
        let targetSQL = "";
        if (_.isString(sql)) {
          targetSQL = buildSQL(sql, _.get(query, "watch", []));
        } else {
          targetSQL = URLStore.bindSQL(sql as QueryStatement);
        }
        // console.log("loading........", queries, db, sql, targetSQL);
        if (targetSQL === "" || db === "") {
          return;
        }

        // add query request into batch
        requests.push(
          ExecService.exec<ResultSet>({
            db: db,
            sql: targetSQL,
          })
        );
      });
      return Promise.allSettled(requests).then((res) => {
        return res
          .map((item) => (item.status === "fulfilled" ? item.value : []))
          .flat();
      });
    },
    {
      // enabled:
    }
  );

  const doQuery = _.debounce(refetch, 200);

  useEffect(() => {
    let disposer: any;
    if (!options?.disableBind) {
      disposer = reaction(
        () => [URLStore.changed, URLStore.forceChanged], // watch params if changed
        () => {
          doQuery();
        }
      );
    }

    return () => {
      if (disposer) {
        disposer();
      }
    };
  }, [doQuery, options]);

  return {
    isLoading: isInitialLoading || isLoading || isFetching,
    isError,
    data,
    error,
  };
}
