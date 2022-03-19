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
import { Form } from "@douyinfe/semi-ui";
import { MetadataSelect } from "@src/components";
import { Variate } from "@src/models";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { MutableRefObject, useRef } from "react";

/**
 *  VariatesSelect based on LinDB tag values metadata.
 * @param props variates config defination
 */
export default function VariatesSelect(props: { variates: Variate[] }) {
  const { variates } = props;
  const formApi = useRef() as MutableRefObject<any>;

  return (
    <Form
      style={{ paddingBottom: 0, paddingTop: 0, display: "inline-flex" }}
      wrapperCol={{ span: 20 }}
      layout="horizontal"
      getFormApi={(api: object) => {
        formApi.current = api;
      }}
      onSubmit={(values: object) => {
        console.log("valuesss...", values);
        URLStore.changeURLParams({ params: values });
      }}
    >
      {_.map(variates, (v: any) => (
        <MetadataSelect
          variate={v}
          key={v.tagKey}
          labelPosition="inset"
          multiple={v.multiple}
        />
      ))}
    </Form>
  );
}
