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
import { Unit } from "@src/models";
import convert from "convert-units";

/**
 * format percent value for 0-100.
 * @param input percent value
 */
function transformPercent(input: number): string {
  if (!input) {
    return "0%";
  } else {
    return `${input.toFixed(2).toString()}%`;
  }
}
/**
 * format percent value for 0-1.
 * @param input percent value
 */
function transformPercent2(input: number): string {
  if (!input) {
    return "0%";
  } else {
    return `${(100 * input).toFixed(2).toString()}%`;
  }
}

function transformBytes(input: number): string {
  if (input > 1024 * 1024 * 1024 * 1024 * 1024) {
    return `${(input / (1024 * 1024 * 1024 * 1024 * 1024)).toFixed(2)} PB`;
  } else if (input > 1024 * 1024 * 1024 * 1024) {
    return `${(input / (1024 * 1024 * 1024 * 1024)).toFixed(2)} TB`;
  } else if (input > 1024 * 1024 * 1024) {
    return `${(input / (1024 * 1024 * 1024)).toFixed(2)} GB`;
  } else if (input > 1024 * 1024) {
    return `${(input / (1024 * 1024)).toFixed(2)} MB`;
  } else if (input > 1024) {
    return `${(input / 1024).toFixed(2)} KB`;
  } else if (!input) {
    return "0 Byte";
  } else {
    return `${input.toFixed(2)} Byte`;
  }
}

function transformSeconds(input: number): string {
  if (input > 365 * 24 * 3600) {
    return `${(input / (365 * 24 * 3600)).toFixed(2)} years`;
  } else if (input > 24 * 3600) {
    return `${(input / (24 * 3600)).toFixed(2)} days`;
  } else if (input > 3600) {
    return `${(input / 3600).toFixed(2)} hours`;
  } else if (input > 60) {
    return `${(input / 60).toFixed(2)} minutes`;
  } else if (!input) {
    return "0 sec";
  } else {
    return `${input.toFixed(2)} sec`;
  }
}

function transformNanoSeconds(input: number, decimals?: number): string {
  if (!input) {
    return "0ns";
  }
  const best = convert(input).from("ns").toBest();
  const value = convert(input)
    .from("ns")
    .to(best.unit as any);
  if (decimals !== undefined) {
    return value.toFixed(decimals) + best.unit;
  } else {
    return Math.floor(value * 100) / 100 + best.unit;
  }
}

function transformMilliseconds(input: number): string {
  if (input > 24 * 3600 * 1000) {
    return `${(input / (24 * 3600 * 1000)).toFixed(2)} days`;
  } else if (input > 3600 * 1000) {
    return `${(input / (3600 * 1000)).toFixed(2)} hours`;
  } else if (input > 10 * 60 * 1000) {
    return `${(input / 60000).toFixed(2)} min`;
  } else if (input > 1000) {
    return `${(input / 1000).toFixed(2)} s`;
  } else if (!input) {
    return "0 ms";
  } else {
    return `${input.toFixed(2)} ms`;
  }
}

function transformNone(input: number): string {
  if (input > 1000 * 1000 * 1000) {
    return `${(input / (1000 * 1000 * 1000)).toFixed(2)} B`;
  } else if (input > 1000 * 1000) {
    return `${(input / (1000 * 1000)).toFixed(2)} M`;
  } else if (input > 1000) {
    return `${(input / 1000).toFixed(2)} K`;
  } else if (!input) {
    return "0";
  } else {
    return `${input.toFixed(2).toString()}`;
  }
}

function format(point: number, unit: Unit): string {
  switch (unit) {
    case Unit.Nanoseconds:
      return transformNanoSeconds(point);
    case Unit.Milliseconds:
      return transformMilliseconds(point);
    case Unit.Seconds:
      return transformSeconds(point);
    case Unit.Bytes:
      return transformBytes(point);
    case Unit.Percent:
      return transformPercent(point);
    case Unit.Percent2:
      return transformPercent2(point);
    default:
      return transformNone(point);
  }
}

function toObject(input: string): any {
  switch (input) {
    case "false":
      return false;
    case "true":
      return true;
    default:
      return input;
  }
}

export default {
  format,
  toObject,
};
