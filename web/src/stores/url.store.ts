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
import { makeObservable, observable, action } from "mobx";
import * as _ from "lodash-es";
import { History } from "history";

class URLStore {
  public changed: boolean = false;
  public forceChanged: boolean = false;
  public params: URLSearchParams = new URLSearchParams();
  public path: string = "";

  private history: History | undefined = undefined;

  constructor() {
    makeObservable(this, {
      changed: observable,
      forceChanged: observable,
      params: observable,
      applyURLChange: action,
      forceChange: action,
    });
  }

  public setHistory(history: History) {
    if (this.history) {
      // history already register.
      return;
    }
    history.listen(() => {
      console.log("linten.........................");
      this.applyURLChange();
    });
    this.history = history;
    this.applyURLChange();
  }

  getTagConditions(tagKeys: string[]): string[] {
    const tags: string[] = [];
    (tagKeys || []).forEach((item: string) => {
      const watchValues = this.params.getAll(item);
      const tagValues: string[] = _.map(watchValues, (v: string) => `'${v}'`);
      if (tagValues.length > 0) {
        tags.push(`'${item}' in (${tagValues.join(",")})`);
      }
    });
    return tags;
  }

  getTimeRange(): string {
    const times: string[] = [];
    ["from", "to"].forEach((item: string) => {
      const time = this.params.get(item);
      if (time) {
        const val = time.indexOf("now()") >= 0 ? time : `'${time}'`;
        switch (item) {
          case "from":
            times.push(`time>${val}`);
            break;
          case "to":
            times.push(`time<${val}`);
            break;
        }
      }
    });
    return times.join(" and ");
  }

  getParamKeys(): string[] {
    const rs: string[] = [];
    for (var key of this.params.keys()) {
      rs.push(key);
    }
    return rs;
  }

  forceChange() {
    this.forceChanged = !this.forceChanged;
  }

  changeURLParams(p: {
    path?: string;
    params?: { [key: string]: any };
    needDelete?: string[];
    clearAll?: boolean;
    clearTime?: boolean;
  }): void {
    const { params, needDelete, clearAll, clearTime, path } = p;
    const { hash } = window.location;
    const oldSearchParams = this.getSearchParams();
    const searchParams = clearAll
      ? new URLSearchParams()
      : this.getSearchParams();
    let pathname = hash;
    if (hash.startsWith("#")) {
      pathname = hash.substring(1, hash.length);
    }
    if (pathname.indexOf("?") > -1) {
      pathname = pathname.split("?")[0];
    }

    if (!clearAll) {
      (needDelete || []).map((key) => {
        searchParams.delete(key);
      });
    } else {
      if (!clearTime) {
        if (oldSearchParams.has("from")) {
          searchParams.set("from", oldSearchParams.get("from")!);
        }
        if (oldSearchParams.has("to")) {
          searchParams.set("to", oldSearchParams.get("to")!);
        }
      }
    }
    this.updateSearchParams(searchParams, params || {});
    // Because of Hash history cannot PUSH the same path so delete the logic of checking path consistency
    const paramsStr = searchParams.toString();
    if (oldSearchParams.toString() !== paramsStr || path) {
      this.history?.push(
        `${path ? path : pathname}${paramsStr && `?${paramsStr}`}`
      );
    }
  }
  applyURLChange(): void {
    this.params = this.getSearchParams();
    this.path = _.get(this.history, "location.pathname", "");
    this.changed = !this.changed;
    console.log("apply........", this.history?.location.pathname);
  }

  getSearchParams(): URLSearchParams {
    // console.log(
    //   window.location.href,
    //   window.location.href.indexOf("?"),
    //   window.location.href.split("?")[1]
    // );
    if (window.location.href.indexOf("?") > -1) {
      return new URLSearchParams(window.location.href.split("?")[1]);
    } else {
      return new URLSearchParams();
    }
  }
  updateSearchParams(
    searchParams: URLSearchParams,
    params: { [key: string]: any }
  ) {
    for (let k of Object.keys(params)) {
      const v = params[`${k}`];
      if (k) {
        if (!_.isEmpty(v)) {
          if (Array.isArray(v)) {
            searchParams.delete(k);
            v.forEach((oneValue) => searchParams.append(k, oneValue));
          } else {
            searchParams.set(k, v);
          }
        } else {
          searchParams.delete(k);
        }
      }
    }
  }
}

export default new URLStore();
