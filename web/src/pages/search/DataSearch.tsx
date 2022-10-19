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
  IconHelpCircleStroked,
  IconLineChartStroked,
  IconPlay,
} from "@douyinfe/semi-icons";
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
import {
  ExplainStatsView,
  MetadataSelect,
  SimpleStatusTip,
  StatusTip,
} from "@src/components";
import { SQL } from "@src/constants";
import { useMetric, useParams } from "@src/hooks";
import { URLStore } from "@src/stores";
import * as monaco from "monaco-editor";
import React, { MutableRefObject, useEffect, useRef } from "react";
import * as _ from "lodash-es";
import { ExecService } from "@src/services";
import { ChartType, Metadata, ResultSet } from "@src/models";
import { useQuery } from "@tanstack/react-query";
import CanvasChart from "@src/components/chart/CanvasChart";
import { ChartKit } from "@src/utils";
const { Text } = Typography;

const SearchForm: React.FC = () => {
  const sqlEditor = useRef() as MutableRefObject<any>;
  const sqlEditorRef = useRef() as MutableRefObject<HTMLDivElement | null>;
  const { sql } = useParams(["sql"]);

  useEffect(() => {
    if (sqlEditor.current) {
      sqlEditor.current.setValue(sql);
      return;
    }
    if (sqlEditorRef.current) {
      // if editor not init, create it
      monaco.languages.registerCompletionItemProvider("sql", {
        provideCompletionItems: function (_model, _position) {
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
          top: 4,
          bottom: 4,
        },
        automaticLayout: true,
        lineNumbers: "off",
        glyphMargin: false,
        folding: false,
        // Undocumented see https://github.com/Microsoft/vscode/issues/30795#issuecomment-410998882
        lineDecorationsWidth: 4,
        lineNumbersMinChars: 0,
        theme: "lindb",
        fontSize: 14,
        minimap: { enabled: false },
        wordWrap: "on",
        wrappingIndent: "same",
      });
    }
  }, [sql]);

  return (
    <Form style={{ paddingTop: 0, paddingBottom: 0 }}>
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
            style={{
              border: "1px solid var(--semi-color-primary)",
              marginTop: 12,
              marginBottom: 12,
            }}
          >
            <div ref={sqlEditorRef} style={{ height: 160 }} />
          </div>
        </Col>
      </Row>
      <Button
        style={{ marginRight: 12 }}
        icon={<IconPlay size="large" />}
        onClick={() => {
          const sql = _.trim(sqlEditor.current?.getValue());
          URLStore.changeURLParams({ params: { sql: sql }, forceChange: true });
        }}
      >
        Search
      </Button>
    </Form>
  );
};

const SearchMetadata: React.FC = () => {
  const { db, sql } = useParams(["db", "sql"]);

  const { isInitialLoading, isError, error, data } = useQuery(
    ["show_metadata", sql, db, URLStore.forceChanged],
    async () => {
      return ExecService.exec<Metadata>({
        sql: sql,
        db: db,
      });
    },
    { enabled: !_.isEmpty(db) && !_.isEmpty(sql) }
  );

  if (isInitialLoading || isError) {
    return (
      <StatusTip
        isLoading={isInitialLoading}
        isError={isError}
        error={error}
        style={{ marginTop: 50, marginBottom: 50 }}
      />
    );
  }
  return (
    <List
      header={
        <Typography.Title heading={6}>
          {_.startCase(data?.type)}
        </Typography.Title>
      }
      size="small"
      bordered
      dataSource={data?.values}
      renderItem={(item: any) => (
        <List.Item>
          {data?.type === "field" ? `${item.name}(${item.type})` : item}
        </List.Item>
      )}
    />
  );
};

const SearchData: React.FC = () => {
  const type = ChartType.Line;
  const { db, sql } = useParams(["db", "sql"]);

  const { isLoading, isError, data, error } = useMetric([
    {
      db: db,
      sql: sql,
    },
  ]);
  const content = () => {
    const explainState = _.get(data, "[0].stats", null);
    if (explainState) {
      return <ExplainStatsView state={explainState} />;
    }
    const datasets = ChartKit.createDatasets(data as ResultSet[], type);
    return (
      <div style={{ height: 400 }}>
        <CanvasChart
          type={type}
          datasets={datasets}
          config={{ options: { zoom: false } }}
        />
      </div>
    );
  };
  return (
    <Card
      bodyStyle={{ padding: 8 }}
      headerStyle={{ padding: 6 }}
      title={
        <Space align="center" className="lin-small-space">
          <IconLineChartStroked />
          <Text>{sql}</Text>
        </Space>
      }
      headerExtraContent={
        <SimpleStatusTip
          isLoading={isLoading}
          isError={isError}
          error={error}
        />
      }
    >
      {content()}
    </Card>
  );
};

const DataSearch: React.FC = () => {
  const { db, sql } = useParams(["sql", "db"]);
  const isDataSearch = (sql: string): boolean => {
    const sqlOfLowerCase = _.lowerCase(sql);
    return (
      _.startsWith(sqlOfLowerCase, "select") ||
      _.startsWith(sqlOfLowerCase, "explain")
    );
  };

  const content = () => {
    if (_.isEmpty(db) || _.isEmpty(sql)) {
      return null;
    }
    return isDataSearch(sql) ? <SearchData /> : <SearchMetadata />;
  };

  return (
    <>
      <Card style={{ marginBottom: 12 }} bodyStyle={{ padding: 12 }}>
        <SearchForm />
      </Card>
      {content()}
    </>
  );
};

export default DataSearch;
