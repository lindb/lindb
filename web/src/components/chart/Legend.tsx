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
import classNames from "classnames";
import React, { useRef } from "react";
import { Chart } from "chart.js";
import * as _ from "lodash-es";
import { handleSeriesClick } from "./util";
import { Unit } from "@src/models";
import { FormatKit } from "@src/utils";

const LegendHeader: React.FC<{ values: string[] }> = (props) => {
  const { values } = props;
  return (
    <div className="chart-legend-table-header">
      <span className="chart-legend-th-content no-pointer">
        {/*search input*/}
      </span>
      {values.map((key) => {
        const headerClass = classNames("chart-legend-th-content", {
          order: false,
          desc: false,
        });

        return (
          <span
            key={key}
            className={headerClass}
            // onClick={() => handleSort(key)}
          >
            <span>{_.upperFirst(key)}</span>
          </span>
        );
      })}
    </div>
  );
};

const LegendItem: React.FC<{
  series: any;
  values: string[];
  chart: Chart | null;
}> = (props) => {
  const { series, values, chart } = props;
  const unit = _.get(chart, "lin.extend.unit", Unit.Short);
  const seriesDiv = useRef<HTMLDivElement>(null);
  const { borderColor, label, hidden, aggregateValues } = series;
  const seriesCls = classNames("chart-legend-series", {
    fade: hidden,
  });

  return (
    <div
      className={seriesCls}
      onClick={(e) => handleSeriesClick(chart, series, e)}
      ref={seriesDiv}
    >
      <span className="legend-series-key">
        <i
          className="legend-series-icon"
          style={{ backgroundColor: borderColor }}
        />
        <span className="legend-series-label">{label}</span>
      </span>
      {values.map((key: string) => (
        <span key={key} className="legend-series-value">
          {FormatKit.format(aggregateValues[key], unit)}
        </span>
      ))}
    </div>
  );
};

const Legend: React.FC<{
  chart: Chart | null;
}> = (props) => {
  const { chart } = props;
  if (!chart) {
    // if chart not exist, return
    return null;
  }

  const datasets = _.get(chart, "data.datasets", []);
  const asTable = _.get(chart, "options.legend.asTable", false);
  const values = _.get(chart, "options.legend.values", []);
  const legendCls = classNames("chart-metric-legend-container", {
    "as-table": asTable,
  });
  const legendContentCls = classNames("chart-legend-content", {
    table: asTable,
    active: _.find(datasets, { hidden: false }),
  });

  return (
    <div className={legendCls}>
      <div className={legendContentCls}>
        {asTable && !_.isEmpty(datasets) && <LegendHeader values={values} />}
        {datasets.map((series: any) => (
          <LegendItem
            chart={chart}
            series={series}
            key={series.label}
            values={values}
          />
        ))}
      </div>
    </div>
  );
};

export default Legend;
