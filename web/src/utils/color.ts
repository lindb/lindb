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
export function getColor(idx: number): string {
  const colors = [
    "#7EB26D", // 0: pale green
    "#EAB839", // 1: mustard
    "#6ED0E0", // 2: light blue
    "#EF843C", // 3: orange
    "#E24D42", // 4: red
    "#1F78C1", // 5: ocean
    "#BA43A9", // 6: purple
    "#705DA0", // 7: violet
    "#508642", // 8: dark green
    "#CCA300", // 9: dark sand
    "#447EBC",
    "#C15C17",
    "#890F02",
    "#0A437C",
    "#6D1F62",
    "#584477",
    "#B7DBAB",
    "#F4D598",
    "#70DBED",
    "#F9BA8F",
    "#F29191",
    "#82B5D8",
    "#E5A8E2",
    "#AEA2E0",
    "#629E51",
    "#E5AC0E",
    "#64B0C8",
    "#E0752D",
    "#BF1B00",
    "#0A50A1",
    "#962D82",
    "#614D93",
    "#9AC48A",
    "#F2C96D",
    "#65C5DB",
    "#F9934E",
    "#EA6460",
    "#5195CE",
    "#D683CE",
    "#806EB7",
    "#3F6833",
    "#967302",
    "#2F575E",
    "#99440A",
    "#58140C",
    "#052B51",
    "#511749",
    "#3F2B5B",
    "#E0F9D7",
    "#FCEACA",
    "#CFFAFF",
    "#F9E2D2",
    "#FCE2DE",
    "#BADFF4",
    "#F9D9F9",
    "#DEDAF7",
  ]; // tslint:disable-line
  return colors[idx % colors.length];
}

/**
 * hex to rgba
 * @param {string} hex Color hex value
 * @param {number} alpha
 * @return {string} rgba
 */
export function toRGBA(hex: string, alpha: number) {
  const hexReg = /^#?([0-9a-fA-f]{3}|[0-9a-fA-f]{6})$/;

  if (hexReg.test(hex)) {
    const validAlpha = Math.max(0, Math.min(1, alpha));
    const validHex = hex.startsWith("#") ? hex.slice(1) : hex;

    const length = validHex.length === 3 ? 1 : 2;

    const color = [];
    for (let i = 0; i < validHex.length; i += length) {
      color.push(parseInt(`0x${validHex.slice(i, i + length)}`));
    }

    return `rgba(${color.join(", ")}, ${validAlpha})`;
  } else {
    return hex;
  }
}
