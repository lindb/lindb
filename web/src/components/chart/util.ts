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
import * as _ from "lodash-es";
import { Chart } from "chart.js";

export enum TOOLTIP_POSITION {
  TOP = "top",
  RIGHT = "right",
  BOTTOM = "bottom",
  LEFT = "left",
}
export const TOOLTIP_SPACING = 12; // Tooltip 与边界的距离
export const TOOLTIP_MIN_HEIGHT = 100; // Chart Tooltip 最小高度
export const CROSSHAIR_SPACING = 20; // 位于左右两侧时，鼠标与 tooltip 边界的间隙

const calcTooltipPosition = (
  tooltipContainer: HTMLElement,
  position: TOOLTIP_POSITION,
  // 当前鼠标相对于边界（boundary）的 offsetX
  mouseX: number,
  // Chart Canvas
  chartCanvasRect: DOMRect,
  // 边界元素
  boundaryRect: DOMRect
) => {
  // 当前 Tooltip 高度
  const tooltipHeight = tooltipContainer.offsetHeight;
  // 当前 Tooltip 宽度
  const tooltipWidth = tooltipContainer.offsetWidth;

  const reachLeftBoundary = mouseX - tooltipWidth / 2 - TOOLTIP_SPACING < 0;
  const reachRightBoundary =
    mouseX + tooltipWidth / 2 + TOOLTIP_SPACING > boundaryRect.width;

  const reachBoundary = reachLeftBoundary || reachRightBoundary;

  let kickPosition, containerPosition;

  switch (position) {
    case TOOLTIP_POSITION.BOTTOM:
      kickPosition = {
        left: mouseX,
      };
      containerPosition = {
        left: !reachBoundary
          ? mouseX - tooltipWidth / 2 // moving
          : reachLeftBoundary
          ? TOOLTIP_SPACING
          : undefined,
        top: chartCanvasRect.bottom - 8,
        right: reachRightBoundary ? 0 : undefined,
      };
      break;
    case TOOLTIP_POSITION.TOP:
      kickPosition = {
        top: chartCanvasRect.y - 4,
        left: mouseX,
      };
      containerPosition = {
        left: !reachBoundary
          ? mouseX - tooltipWidth / 2
          : reachLeftBoundary
          ? TOOLTIP_SPACING
          : undefined,
        top: chartCanvasRect.top - tooltipHeight,
        right: reachRightBoundary ? TOOLTIP_SPACING : undefined,
      };
      break;
    case TOOLTIP_POSITION.LEFT:
      kickPosition = {
        right: -15,
        top: "50%",
      };
      containerPosition = {
        left: mouseX - CROSSHAIR_SPACING - tooltipWidth,
        // left: isRight ? mouseClientX + CROSSHAIR_SPACING : mouseClientX - CROSSHAIR_SPACING - toolTipReact.width - 15,
        top: "50%",
      };
      break;
    case TOOLTIP_POSITION.RIGHT:
      kickPosition = {
        left: 2,
        top: "50%",
      };
      containerPosition = {
        left: mouseX + CROSSHAIR_SPACING,
        top: "50%",
      };
      break;
    default:
      break;
  }

  return {
    container: containerPosition,
    kick: kickPosition,
  };
};

export const getTooltipPositionAndSize = (
  tooltipContainer: HTMLElement,
  offsetX: number,
  chartCanvasRect: DOMRect,
  boundaryRect: DOMRect
): any => {
  // 当前点鼠标对于「边界」的 offsetX 值
  const mouseClientX = offsetX + chartCanvasRect.left - boundaryRect.left;
  // 上方可视区域
  const topAreaHeight = chartCanvasRect.top - boundaryRect.top;
  // 下方可视区域
  const bottomAreaHeight = boundaryRect.bottom - chartCanvasRect.bottom;
  // 左侧可视区域
  const leftAreaWidth = mouseClientX - CROSSHAIR_SPACING;
  // 右侧可视区域
  const rightAreaWidth = boundaryRect.width - mouseClientX - CROSSHAIR_SPACING;

  // 当前 Tooltip 高度
  const tooltipHeight = tooltipContainer.clientHeight;
  // 当前 Tooltip 宽度
  const tooltipWidth = tooltipContainer.clientWidth;
  // 判断所处位置
  let targetPosition: TOOLTIP_POSITION;

  // 上下足以容纳
  if (
    bottomAreaHeight > TOOLTIP_MIN_HEIGHT ||
    topAreaHeight > TOOLTIP_MIN_HEIGHT
  ) {
    if (bottomAreaHeight - TOOLTIP_SPACING > tooltipHeight) {
      targetPosition = TOOLTIP_POSITION.BOTTOM;
    } else if (topAreaHeight - TOOLTIP_SPACING > tooltipHeight) {
      targetPosition = TOOLTIP_POSITION.TOP;
    } else {
      targetPosition =
        bottomAreaHeight > topAreaHeight
          ? TOOLTIP_POSITION.BOTTOM
          : TOOLTIP_POSITION.TOP;
    }
  } else {
    if (leftAreaWidth - TOOLTIP_SPACING > tooltipWidth) {
      targetPosition = TOOLTIP_POSITION.LEFT;
    } else if (rightAreaWidth - TOOLTIP_SPACING > tooltipWidth) {
      targetPosition = TOOLTIP_POSITION.RIGHT;
    } else {
      targetPosition =
        leftAreaWidth > rightAreaWidth
          ? TOOLTIP_POSITION.LEFT
          : TOOLTIP_POSITION.RIGHT;
    }
  }

  let maxSize;

  // 设置高度
  switch (targetPosition) {
    case TOOLTIP_POSITION.BOTTOM:
      maxSize = {
        height: bottomAreaHeight - TOOLTIP_SPACING,
        width: boundaryRect.width - TOOLTIP_SPACING * 2,
      };
      break;
    case TOOLTIP_POSITION.TOP:
      maxSize = {
        height: topAreaHeight - TOOLTIP_SPACING,
        width: boundaryRect.width - TOOLTIP_SPACING * 2,
      };
      break;
    case TOOLTIP_POSITION.LEFT:
      maxSize = {
        height: chartCanvasRect.height - TOOLTIP_SPACING * 2,
        width: leftAreaWidth - TOOLTIP_SPACING,
      };
      break;
    case TOOLTIP_POSITION.RIGHT:
      maxSize = {
        height: chartCanvasRect.height - TOOLTIP_SPACING * 2,
        width: rightAreaWidth - TOOLTIP_SPACING,
      };
      break;
    default:
      break;
  }

  const { container, kick } = calcTooltipPosition(
    tooltipContainer,
    targetPosition,
    mouseClientX,
    chartCanvasRect,
    boundaryRect
  );

  return {
    position: { container, kick },
    size: maxSize,
    direction: targetPosition,
  };
};

export function setPosition(
  target: HTMLElement,
  position: {
    top?: number | string;
    right?: number | string;
    bottom?: number | string;
    left?: number | string;
  }
) {
  if (!target) {
    return;
  }

  const { top, right, bottom, left } = position;
  const topValue =
    top === undefined ? "auto" : typeof top === "number" ? `${top}px` : top;

  const rightValue =
    right === undefined
      ? "auto"
      : typeof right === "number"
      ? `${right}px`
      : right;

  const bottomValue =
    bottom === undefined
      ? "auto"
      : typeof bottom === "number"
      ? `${bottom}px`
      : bottom;

  const leftValue =
    left === undefined ? "auto" : typeof left === "number" ? `${left}px` : left;

  target.style.top = topValue;
  target.style.right = rightValue;
  target.style.bottom = bottomValue;
  target.style.left = leftValue;
}

export function handleSeriesClick(
  chart: Chart | null,
  series: any,
  event: React.MouseEvent
) {
  const onSelect = _.get(chart, "lin.extend.onSeriesSelect", null);
  if (onSelect) {
    onSelect(series, event);
  }
}
