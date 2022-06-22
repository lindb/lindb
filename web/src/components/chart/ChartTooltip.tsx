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
  useState,
  useEffect,
  useLayoutEffect,
  MutableRefObject,
  useCallback,
  useMemo,
} from "react";
import classnames from "classnames";
import { reaction, toJS } from "mobx";
import { ChartEventStore } from "@src/stores";
import * as _ from "lodash-es";
import { Button, ButtonGroup, Badge } from "@douyinfe/semi-ui";
import { IconSearch, IconOrderedList } from "@douyinfe/semi-icons";
import { MouseMoveEvent, UnitEnum } from "@src/models";
import { formatter } from "@src/utils";
// 宏
const TOOLTIP_SPACING = 10; // Tooltip 与边界的距离
const TOOLTIP_MIN_HEIGHT = 100; // Chart Tooltip 最小高度
const CROSSHAIR_SPACING = 20; // 位于左右两侧时，鼠标与 tooltip 边界的间隙

enum TOOLTIP_POSITION {
  TOP = "top",
  RIGHT = "right",
  BOTTOM = "bottom",
  LEFT = "left",
}
// const setContentMaxSize = (
//   el: HTMLElement,
//   size: { height?: number; width?: number } = {}
// ) => {
//   if (!el) {
//     return;
//   }

//   const { height, width } = size;
//   height && (el.style.maxHeight = `${height}px`);
//   width && (el.style.maxWidth = `${width}px`);
// };

const getCanvasActualPixel = (val: number): number => {
  // canvas scale size in retina.
  // ref: https://developer.mozilla.org/en-US/docs/Web/API/Window/devicePixelRatio
  return val / window.devicePixelRatio;
};

const ChartTooltipToolBar = () => {
  return (
    <ButtonGroup>
      <Button
        style={{ paddingLeft: 0, paddingRight: 0 }}
        icon={<IconSearch />}
        size="small"
      />
      <Button
        style={{ paddingLeft: 0, paddingRight: 0 }}
        icon={<IconOrderedList />}
        size="small"
      />
    </ButtonGroup>
  );
};

const setPosition = (
  target: HTMLElement,
  position:
    | {
        top?: number | string;
        right?: number | string;
        bottom?: number | string;
        left?: number | string;
      }
    | undefined
) => {
  if (!target || !position) {
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
};

export default function ChartToltip() {
  const container = useRef() as MutableRefObject<HTMLDivElement>; // tooltip container
  const arrow = useRef() as MutableRefObject<HTMLDivElement>; // tooltip arrow
  const timer = useRef() as MutableRefObject<number | null>; // control mouse leave tooltip hide
  const [visible, setVisible] = useState<boolean>(false); // control tooltip if display
  const [event, setEvent] = useState<MouseMoveEvent>();
  const [dataSource, setDataSource] = useState<any[]>([]);
  // console.log("iniiiiiiiiiiiiiiiiiiiiiiiiii");
  /* 可视区域 */
  const boundaryRect = document.body.getBoundingClientRect();

  const clearTimer = () => {
    // make sure invoke it after setTimer
    setTimeout(() => {
      if (timer.current) {
        clearTimeout(timer.current);
        timer.current = null;
      }
    }, 0);
  };

  const setTimer = useCallback(() => {
    // console.log("eeeee");
    if (timer.current) {
      // console.log("hhhhhhhhh");
      return;
    }
    timer.current = +setTimeout(() => {
      setVisible(false);
      ChartEventStore.setShowTooltip(false);
      clearTimer();
    }, 200);
  }, []);

  useEffect(() => {
    const disposer = [
      reaction(
        () => ChartEventStore.mouseMoveEvent,
        (e) => {
          if (!e) {
            return;
          }
          setVisible(true);
          ChartEventStore.setShowTooltip(true);
          setEvent(e);
          // console.log("tooltip move ", e);
          if (timer.current) {
            // if previous timer not clear, need clear it
            clearTimer();
          }
        }
      ),
      reaction(
        () => ChartEventStore.mouseLeaveEvent,
        () => {
          console.log("settimer");
          setTimer();
        }
      ),
    ];
    return () => {
      disposer.forEach((d) => d());
    };
  }, [setTimer]);

  useEffect(() => {
    if (!event) {
      return;
    }
    setDataSource([event]);
  }, [event]);

  /**
   * 设置 Tooltip 位置
   * @param mouseX 相对于边界（boundary）的 offsetX
   * @param position 目标位置（上、右、下、左）
   * @param maxSize ToolTip 最大尺寸
   */
  const layoutTooltip = (
    mouseX: number,
    position: TOOLTIP_POSITION,
    maxSize: { height: number; width: number }
  ) => {
    if (!event) {
      return;
    }
    // 当前 Tooltip 高度
    const tooltipHeight = container.current.offsetHeight;
    // 当前 Tooltip 宽度
    const tooltipWidth = container.current.offsetWidth;

    // console.log({mouseX, position, tooltipHeight, tooltipWidth});
    // console.log(maxSize);

    const reachLeftBoundary = mouseX - tooltipWidth / 2 - TOOLTIP_SPACING < 0;
    const reachRightBoundary =
      mouseX + tooltipWidth / 2 + TOOLTIP_SPACING > boundaryRect.width;

    const reachBoundary = reachLeftBoundary || reachRightBoundary;
    const chartCanvasRect = event.chartCanvasRect;
    const chartCanvas = event.chartCanvas;
    let kickPosition, containerPosition;
    switch (position) {
      case TOOLTIP_POSITION.BOTTOM:
        containerPosition = {
          left: !reachBoundary
            ? chartCanvas.offsetLeft + mouseX - tooltipWidth / 2 //move
            : reachLeftBoundary
            ? TOOLTIP_SPACING
            : undefined,
          top: chartCanvasRect.bottom - 14,
          // (getCanvasActualPixel(chartCanvas.height) - chartArea.height) +
          // chartArea.top,
          right: reachRightBoundary ? TOOLTIP_SPACING : undefined,
        };
        kickPosition = {
          left: chartCanvasRect.left + event.mouseX,
        };
        break;
      case TOOLTIP_POSITION.TOP:
        containerPosition = {
          left: !reachBoundary
            ? mouseX - tooltipWidth / 2
            : reachLeftBoundary
            ? TOOLTIP_SPACING
            : undefined,
          top: chartCanvasRect.y - tooltipHeight + 5,
          right: reachRightBoundary ? TOOLTIP_SPACING : undefined,
        };
        kickPosition = {
          top: chartCanvasRect.y + 4,
          left: chartCanvasRect.left + event.mouseX,
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

    // debugger;

    // console.log("kickPosition", kickPosition);
    // console.log("containerPosition", containerPosition);

    setPosition(arrow.current, kickPosition);
    setPosition(container.current, containerPosition);
  };

  // calc tooltip move position
  useLayoutEffect(() => {
    if (!event) {
      // when component init, event is null
      return;
    }
    const chartCanvasRect = event.chartCanvasRect;
    const offsetLeft = event ? event.mouseX : 0;
    // 当前点鼠标对于边界的 offsetX 值
    const mouseClientX = offsetLeft + chartCanvasRect.left - boundaryRect.left;
    // console.log(
    //   "do move...",
    //   chartCanvasRect,
    //   boundaryRect,
    //   offsetLeft,
    //   mouseClientX
    // );

    // 上方可视区域
    const topAreaHeight = chartCanvasRect.top - boundaryRect.top;
    // 下方可视区域
    const bottomAreaHeight = boundaryRect.bottom - chartCanvasRect.bottom;
    // 左侧可视区域
    const leftAreaWidth = mouseClientX - CROSSHAIR_SPACING;
    // 右侧可视区域
    const rightAreaWidth =
      boundaryRect.width - mouseClientX - CROSSHAIR_SPACING;

    // 当前 Tooltip 高度
    const tooltipHeight = container.current.clientHeight;
    // 当前 Tooltip 宽度
    const tooltipWidth = container.current.clientWidth;

    // 判断所处位置
    let targetPosition: TOOLTIP_POSITION;

    // 上下足以容纳
    // console.log(
    //   "bottomAreaHeight",
    //   bottomAreaHeight,
    //   "topAreaHeight",
    //   topAreaHeight
    // );
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

    let maxSize = { height: 0, width: 0 };

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

    // setContentMaxSize(container.current, maxSize);

    // 设置位置
    // let tooltipPos;
    // setTimeout(() => {
    // }, 0);
    layoutTooltip(mouseClientX, targetPosition, maxSize);
  }, [dataSource, event]);

  const datasets = _.get(event, "chart.config.data.datasets", []);

  return (
    <div
      ref={container}
      className={classnames("lin-chart-tooltip", {
        fixed: true,
        // "in-chart": true,
        hide: !visible,
      })}
      onMouseEnter={clearTimer}
      onMouseLeave={setTimer}
    >
      <div ref={arrow} className="arrow" />
      <div className="title">
        <div className="lindb-chart-tooltip__timestamp">
          {_.get(event, `chart.config.data.timeLabels[${event?.index}]`, null)}
        </div>
        {/* <ChartTooltipToolBar /> */}
      </div>
      <div className="content-wrapper">
        <div className="content">
          <ul className="list">
            {datasets.map((item: any) => (
              <li className="list-item" key={item.label + item.borderColor}>
                <div className="icon">
                  <Badge style={{ backgroundColor: item.borderColor }} dot />
                </div>
                <span className="key">{item.label}</span>
                <span className="value">
                  {formatter(
                    _.get(item, `data[${event?.index}]`, 0),
                    _.get(event, "chart.config._config.unit", UnitEnum.None)
                  )}
                </span>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
}
