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
import { ChartType, LegendAggregateType, ResultSet } from "@src/models";
import { color } from "chart.js/helpers";
import { ColorKit } from "@src/utils";
import moment from "moment";
import { DateTimeFormat } from "@src/constants";
import * as _ from "lodash-es";

function getChartType(type: string): ChartType {
  switch (type) {
    case ChartType.Area:
      return ChartType.Line;
    case ChartType.Line:
      return ChartType.Line;
    default:
      return ChartType.Line;
  }
}

function getGroupByTags(tags: any) {
  if (!tags) {
    return "";
  }
  const tagKeys = Object.keys(tags);
  if (tagKeys.length === 1) {
    return tags[tagKeys[0]];
  }
  const result = [];
  for (let key of tagKeys) {
    result.push(`${key}:${tags[key]}`);
  }
  return result.join(",");
}

function createDatasets(resultSet: ResultSet[], chartType: ChartType) {
  const datasets: any[] = [];
  //TODO: calc min interval/max time range
  let timeCtx = {
    startTime: 0,
    endTime: 0,
    interval: 0,
  };
  let colorIdx = 0;
  (resultSet || []).forEach((rs: ResultSet) => {
    const { series, startTime, endTime, interval } = rs;

    if (_.isEmpty(series)) {
      return;
    }
    timeCtx = { startTime, endTime, interval };
    series.forEach((item) => {
      const { tags, fields } = item;

      if (!fields) {
        return;
      }

      const groupName = getGroupByTags(tags);

      for (let key of Object.keys(fields)) {
        const bgColor = ColorKit.getColor(colorIdx++);

        const fill = chartType === "area";
        const borderColor = bgColor;
        const backgroundColor =
          chartType === "area"
            ? color(bgColor).alpha(0.25).rgbString()
            : bgColor;
        const label = groupName ? `${groupName}:${key}` : key;
        const pointBackgroundColor = ColorKit.toRGBA(bgColor, 0.25);

        let data: any = [];
        const points: { [timestamp: string]: number } = fields![key];
        let i = 0;
        let timestamp = startTime! + i * interval!;
        let total = 0;
        let count = 0;
        let max = -Infinity;
        let min = Infinity;
        let current = null;
        for (; timestamp <= endTime!; ) {
          const value = points[`${timestamp}`];
          if (value !== null) {
            count++;
          }
          total += value ?? 0;
          current = value ?? current;
          max = Math.max(max, value ?? -Infinity);
          min = Math.min(min, value ?? Infinity);

          const v = value ? Math.floor(value * 1000) / 1000 : 0;
          data.push(v);
          i++;
          timestamp = startTime! + i * interval!;
        }

        datasets.push({
          label,
          data,
          fill,
          backgroundColor,
          borderColor,
          pointBackgroundColor,
          hidden: false,
          aggregateValues: {
            [LegendAggregateType.TOTAL]: total,
            [LegendAggregateType.MAX]: max === -Infinity ? 0 : max,
            [LegendAggregateType.MIN]: min === Infinity ? 0 : min,
            [LegendAggregateType.AVG]: count > 0 ? total / count : 0,
            [LegendAggregateType.CURRENT]: current === null ? 0 : current,
          },
        });
      }
    });
  });
  if (_.isEmpty(datasets)) {
    // no data in response
    return;
  }
  const labels = [];
  const { startTime, endTime, interval } = timeCtx;
  const start = new Date(startTime!);
  const end = new Date(endTime!);
  const showTimeLabel =
    start.getDate() !== end.getDate() ||
    start.getMonth() !== end.getMonth() ||
    start.getFullYear() !== end.getFullYear();
  const range = endTime! - startTime!;
  let i = 0;
  let timestamp = startTime! + i * interval!;
  const times = [];
  const timeLabels = [];
  for (; timestamp <= endTime!; ) {
    const dateTime = moment(timestamp);
    if (showTimeLabel) {
      labels.push(dateTime.format("MM/DD HH:mm"));
    } else if (range > 5 * 60 * 1000) {
      labels.push(dateTime.format("HH:mm"));
    } else {
      labels.push(dateTime.format("HH:mm:ss"));
    }
    timeLabels.push(dateTime.format(DateTimeFormat));
    times.push(timestamp);
    i++;
    timestamp = startTime! + i * interval!;
  }
  return { labels, datasets, interval, times, timeLabels };
}

export default {
  createDatasets,
  getChartType,
};
