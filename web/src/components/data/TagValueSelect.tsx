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
import {
  Form,
  Notification,
  Tag,
  Typography,
  useFormApi,
} from "@douyinfe/semi-ui";
import { useWatchURLChange } from "@src/hooks";
import { Variate, Metadata } from "@src/models";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { MutableRefObject, useRef, useState } from "react";
import { ExecService } from "@src/services";

const Text = Typography.Text;

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
  const dropdownVisible = useRef() as MutableRefObject<boolean>;
  const formApi = useFormApi();

  useWatchURLChange(() => {
    const tagsStr = URLStore.params.get("tags");
    if (tagsStr && tagsStr?.length > 0) {
      try {
        const tags = JSON.parse(tagsStr);
        formApi.setValue(variate.tagKey, _.get(tags, variate.tagKey, []));
      } catch (err) {
        formApi.setValue(variate.tagKey, []);
      }
    }
  });

  const renderMultipleWithCustomTag = (optionNode: any, { onClose }: any) => {
    const content = (
      <Tag avatarShape="square" closable={true} onClose={onClose} size="large">
        {optionNode.label}
      </Tag>
    );
    return {
      isRenderInTag: false,
      content,
    };
  };

  const fetchTagValues = async (prefix?: string) => {
    setLoading(true);
    try {
      let showTagValuesSQL = variate.sql;

      if (!_.isEmpty(prefix)) {
        showTagValuesSQL += ` where ${variate.tagKey} like '${prefix}*'`;
      }

      const metadata = await ExecService.exec<Metadata | string[]>({
        sql: showTagValuesSQL,
        db: variate.db,
      });
      const values = (metadata as Metadata).values;
      const optionList: any[] = [];
      (values || []).map((item: any) => {
        optionList.push({ value: item, label: item });
      });
      setOptionList(optionList);
    } catch (err) {
      Notification.error({
        title: "Fetch tag values error",
        content: _.get(err, "response.data", "Unknown internal error"),
        position: "top",
        theme: "light",
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
  };

  // lazy find metadata when user input.
  const search = _.debounce(fetchTagValues, 200);

  return (
    <>
      <Form.Select
        style={style}
        multiple
        field={variate.tagKey}
        label={
          <Text link strong style={{ marginLeft: 4 }}>
            {variate.tagKey}
          </Text>
        }
        placeholder={placeholder}
        optionList={optionList}
        labelPosition={labelPosition}
        showClear
        renderSelectedItem={renderMultipleWithCustomTag}
        filter
        remote
        onSearch={(input: string) => {
          search(input);
        }}
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
        onFocus={() => fetchTagValues()}
      />
    </>
  );
};

export default TagValueSelect;
