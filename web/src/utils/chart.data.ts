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
import { ChartType, ResultSet } from "@src/models";
import { ChartDataset } from "chart.js/auto";
import { color } from "chart.js/helpers";
import { getColor, toRGBA } from "@src/utils";
import moment from "moment";
import { DateTimeFormat } from "@src/constants";

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

export function buildLineChart(resultSet: ResultSet | null, chart: any) {
  if (!resultSet) {
    return;
  }

  const { series, startTime, endTime, interval } = resultSet;

  if (!series || series.length === 0) {
    return;
  }

  const datasets: ChartDataset[] = [];
  let colorIdx = 0;
  let leftMax = 0;
  const chartType = chart?.type || ChartType.Line;

  series.forEach((item) => {
    const { tags, fields } = item;

    if (!fields) {
      return;
    }

    const groupName = getGroupByTags(tags);

    for (let key of Object.keys(fields)) {
      const bgColor = getColor(colorIdx++);

      const fill = chartType === "area";
      const borderColor = bgColor;
      const backgroundColor =
        chartType === "area" ? color(bgColor).alpha(0.25).rgbString() : bgColor;
      const label = groupName ? `${groupName}:${key}` : key;
      const pointBackgroundColor = toRGBA(bgColor, 0.25);

      let data: any = [];
      const points: { [timestamp: string]: number } = fields[key];
      let i = 0;
      let timestamp = startTime! + i * interval!;
      for (; timestamp <= endTime!; ) {
        const value = points[`${timestamp}`];
        const v = value ? Math.floor(value * 1000) / 1000 : 0;
        if (leftMax < v) {
          leftMax = v;
        }
        data.push(v);
        i++;
        timestamp = startTime! + i * interval!;
      }

      let hidden = false;
      datasets.push({
        label,
        data,
        fill,
        backgroundColor,
        borderColor,
        pointBackgroundColor,
        hidden,
      });
    }
  });
  if (datasets.length == 0) {
    // no data in response
    return;
  }
  const labels = [];
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
  return { labels, datasets, interval, times, timeLabels, leftMax };
}
