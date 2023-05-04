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
import { Metadata, Variate } from "@src/models";
import { ExecService } from "@src/services";
import { URLStore } from "@src/stores";
import { LinSelect } from "@src/components";
import * as _ from "lodash-es";
import React from "react";
import { TemplateKit, URLKit } from "@src/utils";

const MetadataSelect: React.FC<{
  labelPosition?: "top" | "left" | "inset";
  multiple?: boolean;
  type: "db" | "namespace" | "metric" | "field" | "tagKey" | "tagValue";
  variate: Variate;
  placeholder?: string;
  rules?: any[];
  style?: React.CSSProperties;
}> = (props) => {
  const { variate, placeholder, labelPosition, multiple, type, rules, style } =
    props;

  const findMetadata = async (prefix?: string) => {
    const params = URLStore.getParams();
    const db = variate.db ? TemplateKit.template(variate.db, params || {}) : "";
    if (type != "db" && _.isEmpty(db)) {
      // not input db return it, except find db metadata
      return [];
    }
    let targetSQL = variate.sql;

    const whereClause = [];
    if (!_.isEmpty(prefix)) {
      switch (type) {
        case "namespace":
          whereClause.push(`namespace='${prefix}'`);
          break;
        case "metric":
          whereClause.push(`metric='${prefix}'`);
          break;
        case "tagValue":
          whereClause.push(`${variate.tagKey} like '${prefix}*'`);
          break;
      }
    }
    if (type == "tagValue") {
      // build tag where clause
      const tags: string[] = URLKit.getTagConditions(
        params,
        _.get(variate, "watch.cascade", [])
      );
      if (!_.isEmpty(tags)) {
        whereClause.push(...tags);
      }
    }
    if (!_.isEmpty(whereClause)) {
      targetSQL += ` where ${whereClause.join(" and ")}`;
    }

    return ExecService.exec<Metadata | string[]>({
      sql: targetSQL,
      db: db,
    }).then((metadata) => {
      var values: string[];
      if (type === "db") {
        values = metadata as string[];
      } else {
        values = (metadata as Metadata).values as string[];
      }
      const optionList: any[] = [];
      (values || []).map((item: any) => {
        if (type === "field") {
          optionList.push({
            label: `${item.name}(${item.type})`,
            value: item.name,
          });
        } else {
          optionList.push({ value: item, label: item });
        }
      });
      return optionList;
    });
  };

  return (
    <>
      <LinSelect
        style={_.merge({ minWidth: 200 }, style)}
        multiple={multiple}
        field={variate.tagKey}
        placeholder={placeholder}
        labelPosition={labelPosition}
        label={variate.label}
        rules={rules}
        showClear
        filter
        remote
        loader={findMetadata}
        reloadKeys={variate.watch?.cascade}
        clearKeys={variate.watch?.clear}
      />
    </>
  );
};

export default MetadataSelect;
