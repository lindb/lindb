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
import { Form, Notification, useFormApi } from "@douyinfe/semi-ui";
import { useWatchURLChange } from "@src/hooks";
import { Variate } from "@src/models";
import { fetchMetadata } from "@src/services";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { MutableRefObject, useRef, useState } from "react";

interface MetadataSelectProps {
  labelPosition?: "top" | "left" | "inset";
  multiple?: boolean;
  variate: Variate;
  placeholder?: string;
}
const MetadataSelect: React.FC<MetadataSelectProps> = (
  props: MetadataSelectProps
) => {
  const { variate, placeholder, labelPosition, multiple } = props;
  const [optionList, setOptionList] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
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
      const value = multiple
        ? URLStore.params.getAll(field)
        : URLStore.params.get(field);
      // set select value of url params changed
      formApi.setValue(field, value);
    }
  });

  const findMetadata = async () => {
    console.log("feind...", variate.tagKey, loaded.current);
    if (loaded.current) {
      // if data alread load, return it
      return;
    }
    setLoading(true);
    try {
      let showTagValuesSQL = variate.ql;

      if (where.current) {
        showTagValuesSQL += where.current;
      }

      const metadata = await fetchMetadata({
        sql: showTagValuesSQL,
        db: variate.db,
      });
      const { values } = metadata;
      const optionList: any[] = [];
      (values || []).map((item) => {
        optionList.push({ value: item, label: item });
      });
      setOptionList(optionList);
      loaded.current = true; // set tag values alread loaded
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
    formApi.submitForm(); //trigger form submit, after use selected
  };
  return (
    <>
      <Form.Select
        multiple={multiple}
        field={variate.tagKey}
        placeholder={placeholder}
        optionList={optionList}
        labelPosition={labelPosition}
        label={variate.label}
        showClear
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

export default MetadataSelect;
