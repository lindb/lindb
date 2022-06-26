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
import { IconHistogram } from "@douyinfe/semi-icons";
import {
  IllustrationIdle,
  IllustrationIdleDark,
} from "@douyinfe/semi-illustrations";
import {
  Card,
  Empty,
  Form,
  Space,
  Switch,
  Typography,
} from "@douyinfe/semi-ui";
import {
  CanvasChart,
  MetadataSelect,
  MetricStatus,
  TagFilterSelect,
} from "@src/components";
import { Route, SQL } from "@src/constants";
import { useWatchURLChange } from "@src/hooks";
import { QueryStatement } from "@src/models";
import { ChartStore, URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { MutableRefObject, useEffect, useRef, useState } from "react";
import { Link } from "react-router-dom";

const { Text } = Typography;
const chartId = "666666666666666666";

export default function DataExplore() {
  const formApi = useRef() as MutableRefObject<any>;
  const [params, setParams] = useState<any>(URLStore.params);
  const [showLQL, setShowLQL] = useState(false);
  const [sql, setSql] = useState("");
  const tagFilter = useRef() as MutableRefObject<Object>;

  useWatchURLChange(() => {
    if (formApi.current) {
      const tagsStr = URLStore.params.get("tags");
      if (tagsStr && tagsStr?.length > 0) {
        try {
          const tags = JSON.parse(tagsStr);
          const tagsSelected: string[] = [];
          _.mapKeys(tags, (value, key) => {
            if (_.isArray(value) && value.length > 0) {
              tagsSelected.push(`${key}(${value.length})`);
            }
          });
          tagFilter.current = tags;
          formApi.current.setValue("tag", tagsSelected);
        } catch (err) {
          formApi.current.setValue("tag", []);
        }
      }
    }
    setParams(URLStore.params);

    setSql(
      URLStore.bindSQL(
        _.get(
          ChartStore.getChartConfig(chartId),
          "targets[0].sql",
          {}
        ) as QueryStatement
      )
    );
  });

  useEffect(() => {
    // register chart config
    ChartStore.register(chartId, {
      targets: [{ sql: {}, bind: true }],
    });
    return () => {
      // unRegister chart config when component destroy.
      ChartStore.unRegister(chartId);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const renderContent = () => {
    const metric = params.get("metric");
    if (!metric) {
      return (
        <Card>
          <Empty
            image={<IllustrationIdle style={{ width: 150, height: 150 }} />}
            darkModeImage={
              <IllustrationIdleDark style={{ width: 150, height: 150 }} />
            }
            title="Please select metric name"
            style={{ marginTop: 50, minHeight: 400 }}
          />
        </Card>
      );
    }
    return (
      <>
        {showLQL && sql && (
          <Card
            headerStyle={{ padding: 12 }}
            bodyStyle={{ padding: 12 }}
            style={{ marginBottom: 12 }}
          >
            <Text>
              <Link
                target={"_blank"}
                to={`${Route.Search}?db=${params.get("db")}&sql=${sql}`}
              >
                Execute LQL:
              </Link>
            </Text>
            <Text style={{ marginLeft: 8 }}>{sql}</Text>
          </Card>
        )}
        <Card
          title={
            <Space align="center">
              <IconHistogram />
              <Text strong>{params.get("metric")}</Text>
            </Space>
          }
          headerStyle={{ padding: 12 }}
          style={{ marginBottom: 12 }}
          bodyStyle={{ padding: 12 }}
          headerExtraContent={<MetricStatus chartId={chartId} />}
        >
          <Form
            className="lin-tag-filter"
            style={{ marginBottom: 12 }}
            layout="horizontal"
            getFormApi={(api: object) => {
              formApi.current = api;
            }}
            onSubmit={(values: object) => {
              URLStore.changeURLParams({ params: values });
            }}
          >
            <MetadataSelect
              variate={{
                db: params.get("db"),
                namespace: params.get("namespace"),
                tagKey: "field",
                label: "Field",
                sql: `show fields from '${params.get("metric")}'`,
              }}
              type="field"
              labelPosition="inset"
              multiple
            />
            <Form.TagInput
              field="tag"
              prefix="Filter By"
              labelPosition="inset"
              style={{ minWidth: 0 }}
              onRemove={(removedValue: string, _idx: number) => {
                if (tagFilter.current) {
                  URLStore.changeURLParams({
                    params: {
                      tags: JSON.stringify(
                        _.omit(
                          tagFilter.current,
                          removedValue.substring(
                            0,
                            removedValue.lastIndexOf("(")
                          )
                        )
                      ),
                    },
                  });
                }
              }}
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
                sql: `show tag keys from '${params.get("metric")}'`,
              }}
              labelPosition="inset"
              multiple
            />
          </Form>
          <CanvasChart chartId={chartId} height={300} />
        </Card>
      </>
    );
  };

  return (
    <>
      <Card style={{ marginBottom: 12 }} bodyStyle={{ padding: 12 }}>
        <Form
          style={{ paddingBottom: 0, paddingTop: 0 }}
          wrapperCol={{ span: 20 }}
          layout="horizontal"
          onSubmit={(values: object) => {
            URLStore.changeURLParams({ params: values });
          }}
        >
          <MetadataSelect
            variate={{
              tagKey: "db",
              label: "Database",
              sql: SQL.ShowDatabases,
              watch: { clear: ["namespace", "metric"] },
            }}
            labelPosition="inset"
            type="db"
          />
          <MetadataSelect
            variate={{
              db: params.get("db"),
              tagKey: "namespace",
              label: "Namespace",
              sql: "show namespaces",
              watch: { clear: ["metric"] },
            }}
            labelPosition="inset"
          />
          <MetadataSelect
            variate={{
              db: params.get("db"),
              namespace: params.get("namespace"),
              tagKey: "metric",
              label: "Metrics",
              sql: "show metrics",
              watch: {
                clear: ["field", "groupBy"],
              },
            }}
            labelPosition="inset"
          />
          <Space>
            <Switch onChange={setShowLQL} />
            Show LQL
          </Space>
        </Form>
      </Card>
      {renderContent()}
    </>
  );
}
