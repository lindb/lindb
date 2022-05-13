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
import { Button, Card, Col, Form, Row, Space } from "@douyinfe/semi-ui";
import { ExplainStatsView, MetadataSelect } from "@src/components";
import { SQL } from "@src/constants";
import { useWatchURLChange } from "@src/hooks";
import { ChartStore, URLStore } from "@src/stores";
import * as monaco from "monaco-editor";
import React, { MutableRefObject, useEffect, useRef } from "react";

const chartID = "9999999999999999";

export default function DataSearch() {
  const formApi = useRef() as MutableRefObject<any>;
  const sqlEditor = useRef() as MutableRefObject<any>;
  const sqlEditorRef = useRef() as MutableRefObject<HTMLDivElement | null>;

  const query = () => {
    const target = formApi.current.getValues();
    target.sql = sqlEditor.current?.getValue();
    ChartStore.reRegister(chartID, { targets: [target] });
    URLStore.changeURLParams({ params: target, forceChange: true });
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
        provideCompletionItems: function (model, position) {
          // find out if we are completing a property in the 'dependencies' object.
          let suggestions: any[] = [];
          let word = model.getWordUntilPosition(position);
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
        // lineNumbers: "off",
        theme: "lindb",
        fontSize: 14,
      });
    }

    ChartStore.register(chartID, {
      targets: [
        {
          db: URLStore.params.get("db") || "",
          sql: sql,
        },
      ],
    });
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
        <ExplainStatsView chartId={chartID} />
      </Card>
    </>
  );
}
