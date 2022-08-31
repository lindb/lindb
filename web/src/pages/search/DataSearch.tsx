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
import { IconHelpCircleStroked, IconPlay } from "@douyinfe/semi-icons";
import {
  Button,
  Card,
  Col,
  Form,
  Row,
  Space,
  List,
  Typography,
} from "@douyinfe/semi-ui";
import { ExplainStatsView, MetadataSelect } from "@src/components";
import { SQL } from "@src/constants";
import { useWatchURLChange } from "@src/hooks";
import { ChartStore, URLStore } from "@src/stores";
import * as monaco from "monaco-editor";
import React, { MutableRefObject, useEffect, useRef, useState } from "react";
import * as _ from "lodash-es";
import { exec } from "@src/services";
import { Metadata } from "@src/models";

const chartID = "9999999999999999";

export default function DataSearch() {
  const formApi = useRef() as MutableRefObject<any>;
  const sqlEditor = useRef() as MutableRefObject<any>;
  const sqlEditorRef = useRef() as MutableRefObject<HTMLDivElement | null>;
  const [metadata, setMetadata] = useState<Metadata | null>(null);
  const [loading, setLoading] = useState(false);
  const [isMetadata, setIsMetadata] = useState(false);
  const [error, setError] = useState("");

  const isDataSearch = (sql: string): boolean => {
    const sqlOfLowerCase = _.lowerCase(sql);
    return (
      _.startsWith(sqlOfLowerCase, "select") ||
      _.startsWith(sqlOfLowerCase, "explain")
    );
  };

  const fetchMetadata = async (target: any) => {
    try {
      setIsMetadata(true);
      setLoading(true);
      const metadata = await exec<Metadata>(target);
      setMetadata(metadata);
    } catch (err) {
      setError(_.get(err, "response.data", "Unknown internal error"));
    } finally {
      setLoading(false);
    }
    ChartStore.unRegister(chartID);
    URLStore.changeURLParams({ params: target, forceChange: true });
  };

  const query = async () => {
    const target = formApi.current.getValues();
    const sql = _.trim(sqlEditor.current?.getValue());
    target.sql = sql;
    if (isDataSearch(sql)) {
      setIsMetadata(false);
      setMetadata(null);
      ChartStore.reRegister(chartID, { targets: [target] });
      URLStore.changeURLParams({ params: target, forceChange: true });
    } else {
      fetchMetadata(target);
    }
  };

  useWatchURLChange(() => {
    if (formApi.current) {
      formApi.current.setValues({
        db: URLStore.params.get("db"),
      });
    }
    if (sqlEditor.current) {
      sqlEditor.current.setValue(URLStore.params.get("sql"));
    }
  });

  useEffect(() => {
    const sql = URLStore.params.get("sql") || "";
    if (sqlEditorRef.current && !sqlEditor.current) {
      // if editor not init, create it
      monaco.languages.registerCompletionItemProvider("sql", {
        provideCompletionItems: function(_model, _position) {
          // find out if we are completing a property in the 'dependencies' object.
          let suggestions: any[] = [];
          // let word = model.getWordUntilPosition(position);
          suggestions.push({
            label: "select",
            kind: monaco.languages.CompletionItemKind.Property,
            insertText: "select ",
            insertTextRules:
              monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            detail: "select query metric fields/functions",
          });
          return { suggestions: suggestions };
        },
      });
      sqlEditor.current = monaco.editor.create(sqlEditorRef.current, {
        value: sql,
        language: "sql",
        padding: {
          top: 8,
          bottom: 8,
        },
        automaticLayout: true,
        // lineNumbers: "off",
        theme: "lindb",
        fontSize: 14,
        lineNumbers: "off",
        minimap: { enabled: false },
        wordWrap: "on",
        wrappingIndent: "same",
      });
    }
    if (isDataSearch(sql)) {
      setIsMetadata(false);
      ChartStore.register(chartID, {
        targets: [
          {
            db: URLStore.params.get("db") || "",
            sql: sql,
          },
        ],
      });
    } else {
      fetchMetadata({
        db: URLStore.params.get("db") || "",
        sql: sql,
      });
    }

    return () => {
      // unRegister chart config when component destroy.
      ChartStore.unRegister(chartID);
    };
  }, []);
  return (
    <>
      <Card style={{ marginBottom: 12 }} bodyStyle={{ padding: 12 }}>
        <Form
          style={{ paddingTop: 0, paddingBottom: 0 }}
          getFormApi={(api) => (formApi.current = api)}
        >
          <MetadataSelect
            type="db"
            variate={{
              tagKey: "db",
              label: "Database",
              sql: SQL.ShowDatabases,
            }}
            labelPosition="left"
          />
          <Row>
            <Col span={24}>
              <Space align="center">
                <span>LinDB Query Language</span>
                <IconHelpCircleStroked />
              </Space>
              <div
                ref={sqlEditorRef}
                style={{ height: 160, marginTop: 12, marginBottom: 12 }}
              />
            </Col>
          </Row>
          <Button
            style={{ marginRight: 12 }}
            icon={<IconPlay size="large" />}
            onClick={query}
          >
            Search
          </Button>
        </Form>
      </Card>
      <Card style={{ marginTop: 12 }}>
        {isMetadata && (
          <List
            loading={loading}
            emptyContent={
              error ? (
                <Typography.Text type="danger">{error}</Typography.Text>
              ) : (
                "No Result"
              )
            }
            header={
              <Typography.Title heading={6}>
                {_.startCase(metadata?.type)}
              </Typography.Title>
            }
            dataSource={metadata?.values}
            renderItem={(item: any) => (
              <List.Item>
                {metadata?.type === "field"
                  ? `${item.name}(${item.type})`
                  : item}
              </List.Item>
            )}
          />
        )}
        <div style={{ display: isMetadata ? "none" : "block" }}>
          <ExplainStatsView chartId={chartID} />
        </div>
      </Card>
    </>
  );
}
