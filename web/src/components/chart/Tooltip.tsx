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

import React, {
  useRef,
  MutableRefObject,
  useEffect,
  useState,
  useCallback,
} from "react";
import { createPortal } from "react-dom";
import { Chart } from "chart.js";
import { Button, Input } from "@douyinfe/semi-ui";
import classNames from "classnames";
import {
  getTooltipPositionAndSize,
  setPosition,
  handleSeriesClick,
  TOOLTIP_POSITION,
} from "./util";
import * as _ from "lodash-es";
import { MouseEventType } from "@src/stores/platform.store";
import {
  IconDescend2,
  IconOrderedList,
  IconSearchStroked,
} from "@douyinfe/semi-icons";
import { useChartEvent } from "@src/hooks";
import { Unit } from "@src/models";
import { CSSKit, FormatKit } from "@src/utils";

const TooltipTitle: React.FC<{ timestamp?: string }> = (props) => {
  const { timestamp } = props;
  const [searchVisible, setSearchVisible] = useState(false);
  return (
    <div className="tooltip-toolbar">
      <div className="toolbar-header">
        <div className="tooltip-curr-time">{timestamp}</div>
        <div className="tooltip-btn-group">
          <Button
            size="small"
            className="tooltip-toolbar-btn"
            icon={<IconOrderedList />}
          />
          <Button
            size="small"
            className="tooltip-toolbar-btn"
            icon={<IconDescend2 />}
          />
          <Button
            size="small"
            className="tooltip-toolbar-btn"
            icon={<IconSearchStroked />}
            onClick={() => setSearchVisible(!searchVisible)}
          />
        </div>
      </div>
      {searchVisible && (
        <Input
          className="tooltip-toolbar__search-input"
          size="small"
          // value={search}
          // onChange={(e) => onSearch && onSearch(e.target.value)}
          placeholder="Please input series"
        />
      )}
    </div>
  );
};

const TooltipItem: React.FC<{
  series: any;
  index: number;
  unit: Unit;
  chart: Chart | null;
}> = (props) => {
  const { series, index, unit, chart } = props;
  const { borderColor, label, hidden } = series;
  const itemCls = classNames("tooltip-content-list-item", {
    selected: !hidden,
  });
  return (
    <div
      className={itemCls}
      onClick={(e) => handleSeriesClick(chart, series, e)}
    >
      <span className="tooltip-series-key">
        <i
          className="tooltip-series-icon"
          style={{ background: borderColor }}
        />
        <span className="tooltip-series-label">{label}</span>
      </span>
      <span className="tooltip-series-value">
        {FormatKit.format(_.get(series, `data.[${index}]`, 0), unit)}
      </span>
    </div>
  );
};

const TooltipContent: React.FC<{
  datasets: any;
  index: number;
  unit: Unit;
  chart: Chart | null;
}> = (props) => {
  const { datasets, index, unit, chart } = props;

  return (
    <div className="tooltip-list">
      {datasets.map((item: any, idx: number) => (
        <TooltipItem
          key={idx}
          chart={chart}
          series={item}
          index={index}
          unit={unit}
        />
      ))}
    </div>
  );
};

const Tooltip: React.FC<{
  chart: Chart | null;
}> = (props) => {
  const { chart } = props;
  const canvas = chart?.canvas;

  const [visible, setVisible] = useState<boolean>(false);
  const container = useRef() as MutableRefObject<HTMLDivElement>;
  const kick = useRef() as MutableRefObject<HTMLDivElement>;
  const timer = useRef<number | null>();
  const [tooltipPosition, setTooltipPosition] = useState<any>(null);

  const { mouseEvent } = useChartEvent();
  const boundaryRect = useRef() as MutableRefObject<DOMRect | null>;
  const chartRect = useRef() as MutableRefObject<DOMRect | null>;
  const currentIndex = useRef() as MutableRefObject<number | null>;

  const removeTooltip = () => {
    if (timer.current) {
      return;
    }
    timer.current = +setTimeout(() => {
      setVisible(false);
      timer.current = null;
    }, 200);
  };

  useEffect(() => {
    if (tooltipPosition) {
      CSSKit.setStyle(container.current, {
        display: "flex",
      });
    }
  }, [tooltipPosition]);

  const keepTooltip = () => {
    if (!timer.current) {
      return;
    }
    clearTimeout(timer.current);
    timer.current = null;
  };

  const initRect = useCallback(() => {
    if (!canvas) {
      return;
    }
    if (!boundaryRect.current) {
      boundaryRect.current = document.body.getBoundingClientRect();
      chartRect.current = canvas.getBoundingClientRect();
    }
  }, [canvas]);

  const clearRect = () => {
    boundaryRect.current = null;
    chartRect.current = null;
  };

  const handleMouseMove = useCallback(
    (e: MouseEvent) => {
      // disable move
      if (e.metaKey || e.altKey || e.shiftKey || e.ctrlKey) {
        return;
      }

      initRect();

      if (!container.current || !chartRect.current || !boundaryRect.current) {
        return;
      }

      // 计算对应位置及尺寸
      const { position, size, direction } = getTooltipPositionAndSize(
        container.current,
        e.offsetX,
        chartRect.current,
        boundaryRect.current
      );

      // 设置位置
      setPosition(kick.current, position.kick);
      setPosition(container.current, position.container);

      // 设置尺寸
      const { height, width } = size;
      if (height) {
        container.current.style.maxHeight = `${height}px`;
      }
      if (width) {
        container.current.style.maxWidth = `${width}px`;
      }

      // 设置方位
      setTooltipPosition({ position: direction });
      keepTooltip();
    },
    [initRect]
  );

  const handleMouseOut = useCallback((_e: MouseEvent) => {
    clearRect();
  }, []);

  useEffect(() => {
    const {
      type,
      index,
      native,
      chart: chartOfMove,
    } = mouseEvent || ({} as any);
    if (!chart || !native) {
      return;
    }
    switch (type) {
      case MouseEventType.Move:
        if (_.get(chartOfMove, "id", 0) != _.get(chart, "id", 0)) {
          return;
        }
        setVisible(true);
        currentIndex.current = index;
        handleMouseMove(native);
        return;
      case MouseEventType.Out:
        removeTooltip();
        handleMouseOut(native);
        return;
    }
  }, [mouseEvent, chart, handleMouseOut, handleMouseMove]);

  if (!chart) {
    return null;
  }

  const tooltipCls = classNames("chart-metric-tooltip", {
    "in-chart":
      tooltipPosition?.position === TOOLTIP_POSITION.LEFT ||
      tooltipPosition?.position === TOOLTIP_POSITION.RIGHT,
  });

  const tooltip = (
    <div
      className={tooltipCls}
      ref={container}
      onMouseEnter={keepTooltip}
      onMouseLeave={removeTooltip}
    >
      <div ref={kick} className="tooltip-top-kick" />
      <div className="tooltip-title">
        <TooltipTitle
          timestamp={_.get(
            chart,
            `data.timeLabels[${currentIndex.current}]`,
            null
          )}
        />
      </div>
      <div className="tooltip-content-list">
        <TooltipContent
          unit={_.get(chart, "lin.extend.unit", Unit.Short)}
          datasets={_.get(chart, "data.datasets", [])}
          index={currentIndex.current || 0}
          chart={chart}
        />
      </div>
    </div>
  );

  if (!visible) {
    return null;
  }

  return createPortal(tooltip, document.body);
};

export default React.memo(Tooltip);
