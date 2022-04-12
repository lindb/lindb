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
import {
  ChartConfig,
  ChartStatus,
  ResultSet,
  Target,
  QueryStatement,
} from "@src/models";
import { exec } from "@src/services";
import { URLStore } from "@src/stores";
import { buildLineChart } from "@src/utils";
import * as _ from "lodash-es";
import { action, makeObservable, observable, reaction } from "mobx";

class ChartStore {
  chartStatusMap: Map<string, ChartStatus> = new Map<string, ChartStatus>(); // observe chart status
  chartErrMap: Map<string, string> = new Map<string, string>(); // for chart load err msg
  charts: Map<string, ChartConfig> = new Map<string, ChartConfig>(); // for chart confi
  seriesCache: Map<string, any> = new Map<string, any>(); // for chart series data
  stateCache: Map<string, any> = new Map<string, any>(); // for metric explain state data(only store one target in chart config)

  constructor() {
    makeObservable(this, {
      chartStatusMap: observable,
      register: action,
      unRegister: action,
      setChartStatus: action,
    });

    // listen url params if changed
    reaction(
      () => URLStore.changed,
      () => {
        this.load();
      }
    );

    reaction(
      () => URLStore.forceChanged,
      () => {
        this.load(true);
      }
    );
  }

  register(chartUniqueId: string, chart: ChartConfig) {
    if (chart) {
      // for react component register too many times, when state change
      if (this.charts.has(chartUniqueId)) {
        return;
      }
      this.charts.set(chartUniqueId, chart);
      this.chartStatusMap.set(chartUniqueId, ChartStatus.Init);
    }
  }

  reRegister(chartUniqueId: string, chart: ChartConfig) {
    if (chart) {
      this.charts.set(chartUniqueId, chart);
    }
  }

  unRegister(chartUniqueId: string) {
    this.charts.delete(chartUniqueId);
    this.seriesCache.delete(chartUniqueId);
    this.stateCache.delete(chartUniqueId);
    // this.chartTrackingMap.delete(uniqueId);
    this.chartStatusMap.delete(chartUniqueId);
    this.chartErrMap.delete(chartUniqueId);
    // this.statsCache.delete(uniqueId);
    // this.selectedSeries.delete(uniqueId);
  }

  setChartStatus(chartUniqueId: string, status: ChartStatus) {
    this.chartStatusMap.set(chartUniqueId, status);
  }

  public getChartConfig(chartUniqueId: string) {
    return this.charts.get(chartUniqueId);
  }

  private load(forceLoad?: boolean) {
    setTimeout(() => {
      this.chartStatusMap.forEach((_v: ChartStatus, uniqueId: string) => {
        this.loadChartData(uniqueId, forceLoad);
      });
    }, 0);
  }

  private loadChartData(chartUniqueId: string, forceLoad?: boolean) {
    const status: ChartStatus | undefined =
      this.chartStatusMap.get(chartUniqueId);
    if (status && status === ChartStatus.Loading) {
      return;
    }
    this.setChartStatus(chartUniqueId, ChartStatus.Loading);

    const chart = this.charts.get(chartUniqueId);
    // console.log("chart", toJS(chart));
    _.get(chart, "targets", []).forEach((target: Target, _index: number) => {
      console.log("chart store", target);
      const db = target.bind
        ? URLStore.params.get("db") || ""
        : _.get(target, "db", "");
      const sql = _.get(target, "sql", "");
      let targetSQL = "";
      if (_.isString(sql)) {
        targetSQL = this.buildSQL(sql, _.get(target, "watch", []));
      } else {
        targetSQL = URLStore.bindSQL(sql as QueryStatement);
      }
      if (targetSQL === "" || db === "") {
        console.log("sql or db is empty");
        this.seriesCache.delete(chartUniqueId);
        this.setChartStatus(chartUniqueId, ChartStatus.Empty);
        return;
      }
      exec<ResultSet>({
        db: db,
        sql: targetSQL,
      })
        .then((response) => {
          const series: ResultSet | undefined = response;
          const reportData: any = buildLineChart(series!, chart!.config);
          if (reportData) {
            // console.log("series", series, reportData);
            this.seriesCache.set(chartUniqueId, reportData);
            this.stateCache.set(chartUniqueId, series.stats);
            this.setChartStatus(chartUniqueId, ChartStatus.OK);
          } else {
            // no data in response
            this.seriesCache.delete(chartUniqueId);
            this.stateCache.delete(chartUniqueId);
            this.setChartStatus(chartUniqueId, ChartStatus.Empty);
          }
        })
        .catch((err) => {
          console.log("err", err);
          this.seriesCache.delete(chartUniqueId);
          this.chartErrMap.set(
            chartUniqueId,
            _.get(err, "response.data", "Unknown internal error")
          );
          this.setChartStatus(chartUniqueId, ChartStatus.Error);
        });
    });
  }

  private buildSQL(ql: string | undefined, cascade: string[]) {
    if (ql === undefined) {
      return "";
    }
    const tags: string[] = URLStore.getTagConditions(cascade);
    const where: string[] = [];
    const timeRange = URLStore.getTimeRange();
    if (timeRange) {
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
  }
}

export default new ChartStore();
