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
import Chart from "chart.js/auto";
import * as helpers from "chart.js/helpers";
import * as _ from "lodash-es";
/**
 * chart lazy plugin ref:
 * 1. https://github.com/chartjs/chartjs-plugin-deferred
 * 2. https://github.com/samber/chartjs-plugin-datasource-prometheus
 */
const STUB_KEY = "$chartjs_lazy";
const MODEL_KEY = "$lazy";

function computeOffset(value: any, base: number) {
  var number = parseInt(value, 10);
  if (isNaN(number)) {
    return 0;
  } else if (typeof value === "string" && value.indexOf("%") !== -1) {
    return (number / 100) * base;
  }
  return number;
}
function chartInViewport(chart: any) {
  var options = chart[MODEL_KEY].options;
  var canvas = _.get(chart, "canvas", null);

  // http://stackoverflow.com/a/21696585
  if (!canvas || canvas.offsetParent === null) {
    return false;
  }

  var rect = canvas.getBoundingClientRect();
  var dy = computeOffset(options.yOffset || 0, rect.height);
  var dx = computeOffset(options.xOffset || 0, rect.width);

  return (
    rect.right - dx >= 0 &&
    rect.bottom - dy >= 0 &&
    rect.left + dx <= window.innerWidth &&
    rect.top + dy <= window.innerHeight
  );
}
function isScrollable(node: any) {
  var type = node.nodeType;
  if (type === Node.ELEMENT_NODE) {
    var overflowX = helpers.getStyle(node, "overflow-x");
    var overflowY = helpers.getStyle(node, "overflow-y");
    return (
      overflowX === "auto" ||
      overflowX === "scroll" ||
      overflowY === "auto" ||
      overflowY === "scroll"
    );
  }

  return node.nodeType === Node.DOCUMENT_NODE;
}

function defer(fn: any, delay?: any) {
  if (delay) {
    window.setTimeout(fn, delay);
  } else {
    helpers.requestAnimFrame.call(window, fn);
  }
}
function onScroll(event: any) {
  var node = event.target;
  var stub = node[STUB_KEY];
  if (stub.ticking) {
    return;
  }
  console.log("scroll......");

  stub.ticking = true;
  defer(function () {
    var charts = stub.charts.slice();
    var ilen = charts.length;
    var chart, i;

    for (i = 0; i < ilen; ++i) {
      chart = charts[i];
      if (chartInViewport(chart)) {
          unwatch(chart); // eslint-disable-line
        chart[MODEL_KEY].appeared = true;
        chart.update();
      }
    }

    stub.ticking = false;
  });
}

function watch(chart: any) {
  console.log("watch ....", chart);
  var parent = _.get(chart, "canvas.parentElement", null);
  var stub, charts;

  while (parent) {
    console.log("xxxxxxx watch.....", parent);
    if (isScrollable(parent)) {
      stub = parent[STUB_KEY] || (parent[STUB_KEY] = {});
      charts = stub.charts || (stub.charts = []);
      if (charts.length === 0) {
        parent.addEventListener("scroll", onScroll);
        // helpers.addEvent(parent, "scroll", onScroll);
      }
      charts.push(chart);
      chart[MODEL_KEY].elements.push(parent);
      console.log("scrollable parent", parent);
    }

    parent = parent.parentElement || parent.ownerDocument;
  }
}
function unwatch(chart: any) {
  chart[MODEL_KEY].elements.forEach(function (element: any) {
    var charts = element[STUB_KEY].charts;
    charts.splice(charts.indexOf(chart), 1);
    if (!charts.length) {
      //   helpers.removeEvent(element, "scroll", onScroll);
      element.removeEventListener("scroll", onScroll);
      console.log("remove watch event");
      delete element[STUB_KEY];
    }
  });

  chart[MODEL_KEY].elements = [];
}

Chart.register({
  id: "lazy",
  beforeInit: function (chart: any, args: any, options: any) {
    console.log("lazy chart init");
    chart[MODEL_KEY] = {
      options: options,
      appeared: false,
      delayed: false,
      loaded: false,
      elements: [],
    };

    if (!chartInViewport(chart)) {
      // if chart not in view, need add scroll event watch
      watch(chart);
    }
  },
  beforeUpdate: function (chart: any, options: any) {
    console.log("before update", chart);
    const canvas = chart.canvas as HTMLCanvasElement;
    const ctx = chart.ctx;
    ctx.font = "13px Arial";
    ctx.fillStyle = "read";
    ctx.fillText("test", 200, 40);
    ctx.restore();
    if (!chartInViewport(chart)) {
      console.log("disable update");
      return false;
    }
    console.log("enable update");
    return true;
  },
  // afterDraw: function (chart: any) {
  //   const ctx = chart.ctx;
  //   ctx.font = "13px Arial";
  //   ctx.fillStyle = "red";
  //   ctx.fillText("test", 200, 40);
  //   ctx.restore();
  // },
  destroy: function (chart: any) {
    // unwatch(chart);
  },
});
