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
  IconClose,
  IconHelpCircleStroked,
  IconSaveStroked,
} from "@douyinfe/semi-icons";
import {
  Banner,
  Button,
  Card,
  Col,
  Form,
  Row,
  Space,
  useFormApi,
} from "@douyinfe/semi-ui";
import { Route } from "@src/constants";
import { createStorage } from "@src/services";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { useState } from "react";

export default function StorageConfig() {
  const [submiting, setSubmiting] = useState(false);
  const [error, setError] = useState("");

  const Buttons = () => {
    const formApi = useFormApi();
    const create = async (values: any) => {
      try {
        setSubmiting(true);
        await createStorage(values);
      } catch (err) {
        setError(_.get(err, "response.data", "Unknown internal error"));
      } finally {
        setSubmiting(false);
      }
    };
    return (
      <Row style={{ paddingTop: 12 }}>
        <Col offset={4}>
          <Button
            type="secondary"
            icon={<IconSaveStroked />}
            style={{ marginRight: 8 }}
            loading={submiting}
            onClick={() => {
              create(formApi.getValues());
            }}
          >
            Submit
          </Button>
          <Button
            type="tertiary"
            icon={<IconClose />}
            onClick={() =>
              URLStore.changeURLParams({ path: Route.MetadataStorage })
            }
          >
            Cancel
          </Button>
        </Col>
      </Row>
    );
  };

  const labelWithHelp = (label: string) => {
    return (
      <Space align="center">
        <span>{label}</span>
        <IconHelpCircleStroked />
      </Space>
    );
  };

  return (
    <>
      {error && (
        <Banner
          description={error}
          type="danger"
          closeIcon={null}
          style={{ marginBottom: 12, justifyContent: "left" }}
        />
      )}
      <Card bordered={false}>
        <Form
          // layout="horizontal"
          labelPosition="left"
          labelAlign="right"
          labelCol={{ span: 4 }}
          wrapperCol={{ span: 12 }}
        >
          <Form.Input
            field="config.namespace"
            rules={[{ required: true }]}
            label="Namespace"
            helpText="ETCD namespace"
          />
          <Form.TagInput
            field="config.endpoints"
            rules={[{ required: true }]}
            label="Endpoints"
          />
          <Form.Input field="config.timeout" label={labelWithHelp("Timeout")} />
          <Form.Input field="config.dialTimeout" label="DialTimeout" />
          <Form.Input field="config.leaseTTL" label="LeaseTTL" />
          <Form.Input field="config.username" label="Username" />
          <Form.Input
            mode="password"
            field="config.password"
            label="Password"
          />
          <Buttons />
        </Form>
      </Card>
    </>
  );
}
