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
import { Form, Notification, useFormApi, Select } from "@douyinfe/semi-ui";
import { IconAppCenter } from "@douyinfe/semi-icons";
import { useWatchURLChange } from "@src/hooks";
import { Metadata, Variate } from "@src/models";
import { exec } from "@src/services";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { MutableRefObject, useRef, useState } from "react";

interface MetadataSelectProps {
  labelPosition?: "top" | "left" | "inset";
  variate: Variate;
  placeholder?: string;
  style?: React.CSSProperties;
}
const TagValueSelect: React.FC<MetadataSelectProps> = (
  props: MetadataSelectProps
) => {
  const { variate, placeholder, labelPosition, style } = props;
  const [optionList, setOptionList] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const oldVariate = useRef() as MutableRefObject<Variate>;
  const loaded = useRef() as MutableRefObject<boolean>;
  const where = useRef() as MutableRefObject<string>;
  const dropdownVisible = useRef() as MutableRefObject<boolean>;
  const formApi = useFormApi();

  useWatchURLChange(() => {
    // build where clause
    const tags: string[] = URLStore.getTagConditions(
      _.get(variate, "watch.cascade", [])
    );
    let whereClause = "";
    if (tags.length > 0) {
      whereClause = ` where ${tags.join(" and ")}`;
    }
    const field = variate.tagKey;
    if (where.current && whereClause != where.current) {
      where.current = whereClause;
      loaded.current = false; // if where cluase changed, need load tag values
      formApi.setValue(field, null);
    } else {
      where.current = whereClause;
      const value = URLStore.params.getAll(field);
      // set select value of url params changed
      formApi.setValue(field, value);
    }
  });
  const triggerRender = (props: { value: any }) => {
    const { value } = props;
    return (
      <div
        style={{
          minWidth: "112",
          backgroundColor: "var(--semi-color-fill-0)",
          height: 32,
          display: "flex",
          alignItems: "center",
          borderRadius: "var(--semi-border-radius-small)",
          border: "1px solid transparent",
          paddingLeft: 12,
          // borderRadius: 3,
          color: "var(--semi-color-text-2)",
        }}
      >
        <div
          style={{
            fontWeight: 600,
            flexShrink: 0,
            fontSize: 14,
            color: "var(--semi-color-secondary)",
          }}
        >
          {variate.label}
        </div>
        <div
          style={{
            margin: 4,
            color: "var(--semi-color-text-2)",
            whiteSpace: "nowrap",
            textOverflow: "ellipsis",
            flexGrow: 1,
            overflow: "hidden",
          }}
        >
          {value.map((item) => item.label).join(" , ")}
        </div>
        <IconAppCenter style={{ marginRight: 8, flexShrink: 0 }} />
      </div>
    );
  };

  const findMetadata = async () => {
    console.log("find...", variate.tagKey, loaded.current);
    if (loaded.current && _.isEqual(oldVariate.current, variate)) {
      // if data alread load, return it
      return;
    }
    setLoading(true);
    try {
      let showTagValuesSQL = variate.ql;

      if (where.current) {
        showTagValuesSQL += where.current;
      }

      const metadata = await exec<Metadata | string[]>({
        sql: showTagValuesSQL,
        db: variate.db,
      });
      const values = (metadata as Metadata).values;
      const optionList: any[] = [];
      (values || []).map((item: any) => {
        optionList.push({ value: item, label: item });
      });
      setOptionList(optionList);
      loaded.current = true; // set tag values alread loaded
      oldVariate.current = variate;
    } catch (err) {
      Notification.error({
        title: "Fetch tag values error",
        content: _.get(err, "response.data", "Unknown internal error"),
        position: "top",
        duration: 5,
      });
    } finally {
      setLoading(false);
    }
  };
  const handleAfterSelect = () => {
    const clear = _.get(variate, "watch.clear", []);
    clear.forEach((key: string) => {
      formApi.setValue(key, null);
    });
    // formApi.submitForm(); //trigger form submit, after use selected
  };
  return (
    <>
      <Form.Select
        style={style}
        multiple
        field={variate.tagKey}
        placeholder={placeholder}
        optionList={optionList}
        labelPosition={labelPosition}
        showClear
        triggerRender={triggerRender}
        filter
        onBlur={handleAfterSelect}
        onClear={handleAfterSelect}
        onDropdownVisibleChange={(val) => {
          dropdownVisible.current = val;
        }}
        onChange={(_val: any) => {
          if (!dropdownVisible.current) {
            handleAfterSelect();
          }
        }}
        loading={loading}
        onFocus={findMetadata}
      />
    </>
  );
};

export default TagValueSelect;
