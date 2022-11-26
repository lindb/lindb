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
import { QueryStatement } from "@src/models";
import { FormatKit } from "@src/utils";
import { History } from "history";
import * as _ from "lodash-es";
import { action, makeObservable, observable } from "mobx";

class URLStore {
  public changed: boolean = false;
  public forceChanged: boolean = false;
  public path: string = "";

  private params: URLSearchParams = new URLSearchParams();
  private paramObj = {};
  private defaultParams = {}; // default params just save in store, don't put them into url params.
  private history: History | undefined = undefined;

  constructor() {
    makeObservable(this, {
      changed: observable,
      forceChanged: observable,
      applyURLChange: action,
      forceChange: action,
      changeDefaultParams: action,
    });
  }

  public setHistory(history: History | any) {
    if (this.history) {
      // history already register.
      return;
    }
    history.listen(() => {
      this.applyURLChange();
    });
    this.history = history;
    this.applyURLChange();
  }

  public changeDefaultParams(defaultParams: {}, change = true) {
    this.defaultParams = _.merge(
      _.cloneDeep(this.defaultParams),
      defaultParams
    );
    if (change) {
      this.changed = !this.changed;
    }
  }

  public deleteDefaultParams(keys: string[], change = true) {
    if (_.isEmpty(keys)) {
      return;
    }
    const newDefault = _.omit(this.defaultParams, keys);
    if (_.isEqual(this.defaultParams, newDefault)) {
      return;
    }
    this.defaultParams = newDefault;
    if (change) {
      this.changed = !this.changed;
    }
  }

  public forceChange() {
    this.forceChanged = !this.forceChanged;
  }

  public getPath(): string {
    return this.path;
  }

  changeURLParams(p: {
    path?: string;
    params?: { [key: string]: any };
    defaultParams?: { [key: string]: any };
    needDelete?: string[];
    clearAll?: boolean;
    clearTime?: boolean;
    forceChange?: boolean;
  }): void {
    const {
      params,
      defaultParams,
      needDelete,
      clearAll,
      clearTime,
      path,
      forceChange,
    } = p;
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

    if (!_.isEmpty(defaultParams)) {
      this.changeDefaultParams(defaultParams || {}, false);
    }
    this.updateSearchParams(searchParams, params || {});
    // Because of Hash history cannot PUSH the same path so delete the logic of checking path consistency
    const paramsStr = searchParams.toString();
    if (oldSearchParams.toString() !== paramsStr || path) {
      this.history?.push(
        `${path ? path : pathname}${paramsStr && `?${paramsStr}`}`
      );
    }

    if (forceChange) {
      this.forceChange();
    }
  }

  applyURLChange(): void {
    const groupParamsByKey = (params: any) =>
      [...params.entries()].reduce((acc, tuple) => {
        // getting the key and value from each tuple
        const [key, val] = tuple;
        const v = FormatKit.toObject(val);
        if (acc.hasOwnProperty(key)) {
          // if the current key is already an array, we'll add the value to it
          if (Array.isArray(acc[key])) {
            acc[key] = [...acc[key], v];
          } else {
            // if it's not an array, but contains a value, we'll convert it into an array
            // and add the current value to it
            acc[key] = [acc[key], v];
          }
        } else {
          // plain assignment if no special case is present
          acc[key] = v;
        }

        return acc;
      }, {});

    this.params = this.getSearchParams();
    this.paramObj = groupParamsByKey(this.params);
    const newPath = _.get(this.history, "location.pathname", "");
    if (newPath != this.path) {
      // if change path need clear default params
      this.defaultParams = {};
    }
    this.path = newPath;
    this.changed = !this.changed;
  }

  public getParams(): object {
    return _.merge(_.cloneDeep(this.defaultParams), this.paramObj);
  }

  private getSearchParams(): URLSearchParams {
    if (window.location.href.indexOf("?") > -1) {
      return new URLSearchParams(window.location.href.split("?")[1]);
    } else {
      return new URLSearchParams();
    }
  }

  private updateSearchParams(
    searchParams: URLSearchParams,
    params: { [key: string]: any }
  ) {
    _.forIn(params, (v, k) => {
      if (k) {
        if (!_.isUndefined(v)) {
          if (Array.isArray(v)) {
            searchParams.delete(k);
            v.forEach((oneValue) =>
              searchParams.append(
                k,
                _.isString(oneValue) ? oneValue : `${oneValue}`
              )
            );
          } else {
            searchParams.set(k, _.isString(v) ? v : `${v}`);
          }
        } else {
          searchParams.delete(k);
        }
      }
    });
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

  bindSQL(stmt: QueryStatement): string {
    stmt.metric = this.params.get("metric") || "";
    stmt.namespace = this.params.get("namespace") || "";
    stmt.field = this.params.getAll("field") || [];
    stmt.groupBy = this.params.getAll("groupBy") || [];
    stmt.tags = JSON.parse(this.params.get("tags") || "{}");
    if (_.isEmpty(stmt.metric) || _.isEmpty(stmt.field)) {
      return "";
    }
    const fields = _.map(stmt.field, (item: string) => `'${item}'`);
    let nsCluase = "";
    if (!_.isEmpty(stmt.namespace)) {
      nsCluase = ` on '${stmt.namespace}'`;
    }
    const whereClause: string[] = [];
    _.mapKeys(stmt.tags, (value, key) => {
      if (_.isArray(value) && value.length > 0) {
        whereClause.push(
          `${key} in (${_.map(value, (item: string) => `'${item}'`)})`
        );
      }
    });
    const timeRange = this.getTimeRange();
    if (timeRange !== "") {
      whereClause.push(timeRange);
    }
    let whereClauseStr = "";
    if (whereClause.length > 0) {
      whereClauseStr = ` where ${_.join(whereClause, " and ")}`;
    }

    let groupByStr = "";
    if (!_.isEmpty(stmt.groupBy)) {
      groupByStr = ` group by ${_.map(
        stmt.groupBy,
        (item: string) => `'${item}'`
      )}`;
    }

    return `select ${fields} from '${stmt.metric}'${nsCluase}${whereClauseStr} ${groupByStr}`;
  }

  getParamKeys(): string[] {
    const rs: string[] = [];
    for (var key of this.params.keys()) {
      rs.push(key);
    }
    return rs;
  }
}

export default new URLStore();
