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
import { IconFilter, IconHistogram } from "@douyinfe/semi-icons";
import {
  Card,
  Form,
  Popover,
  Space,
  Typography,
  SplitButtonGroup,
  Button,
} from "@douyinfe/semi-ui";
import { CanvasChart, MetadataSelect, TagFilterSelect } from "@src/components";
import { useWatchURLChange } from "@src/hooks";
import { Target } from "@src/models";
import { ChartStore, URLStore } from "@src/stores";
import * as _ from "lodash-es";
import queryString from "query-string";
import React, { MutableRefObject, useEffect, useRef, useState } from "react";

const { Text } = Typography;
const chartId = "666666666666666666";

export default function DataExplore() {
  const formApi = useRef() as MutableRefObject<any>;
  const [params, setParams] = useState<any>(URLStore.params);
  const init = useRef() as MutableRefObject<boolean>;
  const target = useRef() as MutableRefObject<Target>;
  target.current = {
    db: "",
    ql: "",
  };

  const buildTarget = (paramsObj: any, init?: boolean) => {
    target.current.db = _.get(paramsObj, "db", "");
    target.current.ql = "";
    const metric = _.get(paramsObj, "metric", "");
    let fields = _.get(paramsObj, "field", []);
    if (_.isArray(fields)) {
      fields = _.map(fields, (item: string) => `'${item}'`);
    } else {
      fields = [fields];
    }
    if (!metric || !fields) {
      return;
    }
    let groupBy = _.get(paramsObj, "groupBy", []);
    let groupByClause = "";
    if (_.isArray(groupByClause)) {
      groupBy = _.map(groupBy, (item: string) => `'${item}'`);
    } else {
      groupBy = [groupBy];
    }
    const groupByStr = _.join(groupBy, ",");

    if (groupByStr) {
      groupByClause = ` group by ${groupByStr}`;
    }

    const ns = _.get(paramsObj, "namespace", null);
    let nsCluase = "";
    if (ns) {
      nsCluase = ` on '${ns}'`;
    }
    target.current.ql = `select ${_.join(
      fields,
      ","
    )} from '${metric}'${nsCluase} ${groupByClause}`;

    if (init) {
      // init target using url params
      // register chart config
      ChartStore.register(chartId, {
        targets: [target.current],
      });
    } else {
      ChartStore.reRegister(chartId, {
        targets: [target.current],
      });
    }

    console.log("buildTarget", paramsObj, target.current);
  };

  useWatchURLChange(() => {
    setParams(URLStore.params);
    if (!init.current) {
      init.current = true;
    }
    console.log("metric", URLStore.params.get("db"));
  });

  // useLayoutEffect(()=>{

  // })

  useEffect(() => {
    console.log(
      "queryString.parse(URLStore.params.toString())",
      URLStore.params.get("db"),
      queryString.parse(URLStore.params.toString())
    );
    buildTarget(queryString.parse(URLStore.params.toString()), true);
    URLStore.forceChange();
    // ChartStore.register(chartId, {
    //   targets: [target.current],
    // });
    return () => {
      // unRegister chart config when component destroy.
      ChartStore.unRegister(chartId);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <Card
        bordered={false}
        style={{ marginBottom: 12 }}
        bodyStyle={{ padding: 12 }}
      >
        <Form
          style={{ paddingBottom: 0, paddingTop: 0 }}
          wrapperCol={{ span: 20 }}
          layout="horizontal"
          getFormApi={(api: object) => {
            formApi.current = api;
          }}
          onSubmit={(values: object) => {
            console.log("set valuesssss...", values);
            URLStore.changeURLParams({ params: values });
            console.log("after", URLStore.params.get("db"));
          }}
        >
          <MetadataSelect
            variate={{ tagKey: "db", label: "Database", ql: "show databases" }}
            labelPosition="inset"
            type="db"
          />
          <MetadataSelect
            variate={{
              db: params.get("db"),
              tagKey: "namespace",
              label: "Namespace",
              ql: "show namespaces",
            }}
            labelPosition="inset"
          />
          <MetadataSelect
            variate={{
              db: params.get("db"),
              namespace: params.get("namespace"),
              tagKey: "metric",
              label: "Metrics",
              ql: "show metrics",
              watch: {
                clear: ["field", "groupBy"],
              },
            }}
            labelPosition="inset"
          />
        </Form>
      </Card>
      <Card
        title={
          <Space align="center">
            <IconHistogram />
            <Text strong>{params.get("metric")}</Text>
          </Space>
        }
        headerStyle={{ padding: 12 }}
        bordered={false}
        style={{ marginBottom: 12 }}
        bodyStyle={{ padding: 12 }}
      >
        <Form
          style={{ marginBottom: 4 }}
          layout="horizontal"
          onSubmit={(values: object) => {
            buildTarget(_.merge(formApi.current.getValues(), values));
            URLStore.changeURLParams({ params: values });
          }}
        >
          <MetadataSelect
            variate={{
              db: params.get("db"),
              namespace: params.get("namespace"),
              tagKey: "field",
              label: "Field",
              ql: `show fields from '${params.get("metric")}'`,
            }}
            type="field"
            labelPosition="inset"
            multiple
          />
          <Form.Input
            field="name"
            label="Filter By"
            trigger="blur"
            style={{ width: 250 }}
            initValue="Filter(0)"
            labelPosition="inset"
            suffix={
              <TagFilterSelect
                db={params.get("db")}
                metric={params.get("metric")}
              />
            }
          />
          <MetadataSelect
            variate={{
              db: params.get("db"),
              namespace: params.get("namespace"),
              tagKey: "groupBy",
              label: "Group By",
              ql: `show tag keys from '${params.get("metric")}'`,
            }}
            labelPosition="inset"
            multiple
          />
        </Form>
        <CanvasChart chartId={chartId} />
      </Card>
    </>
  );
}
