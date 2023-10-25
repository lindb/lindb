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
  Tooltip,
  useFormApi,
} from "@douyinfe/semi-ui";
import { ClusterType, Route } from "@src/constants";
import { UIContext } from "@src/context/UIContextProvider";
import { ExecService } from "@src/services";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { useContext, useState } from "react";

export const ClusterConfig: React.FC<{ type: ClusterType }> = (props) => {
  const { type } = props;
  const [submiting, setSubmiting] = useState(false);
  const [error, setError] = useState("");
  const { locale } = useContext(UIContext);
  const { MetadataClusterView, Common } = locale;

  const Buttons = () => {
    const formApi = useFormApi();
    const create = async (values: any) => {
      try {
        setSubmiting(true);
        await ExecService.exec({
          sql: `create ${
            type === ClusterType.Storage ? "storage" : "broker"
          } ${JSON.stringify(values)}`,
        });
        URLStore.changeURLParams({
          path:
            type === ClusterType.Storage
              ? Route.MetadataStorage
              : Route.MetadataBroker,
        });
      } catch (err) {
        setError(_.get(err, "response.data", Common.unknownInternalError));
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
            {Common.submit}
          </Button>
          <Button
            type="tertiary"
            icon={<IconClose />}
            onClick={() =>
              URLStore.changeURLParams({
                path:
                  type === ClusterType.Storage
                    ? Route.MetadataStorage
                    : Route.MetadataBroker,
              })
            }
          >
            {Common.cancel}
          </Button>
        </Col>
      </Row>
    );
  };

  const labelWithHelp = (label: string) => {
    return (
      <Space align="center">
        <span>{label}</span>
        <Tooltip content={MetadataClusterView.timeoutTooltip}>
          <IconHelpCircleStroked />
        </Tooltip>
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
      <Card>
        <Form
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
            label={MetadataClusterView.endpoints}
          />
          <Form.Input
            field="config.timeout"
            label={labelWithHelp(MetadataClusterView.timeout)}
          />
          <Form.Input
            field="config.dialTimeout"
            label={MetadataClusterView.dialTimeout}
          />
          <Form.Input
            field="config.leaseTTL"
            label={MetadataClusterView.leaseTTL}
          />
          <Form.Input
            field="config.username"
            label={MetadataClusterView.username}
          />
          <Form.Input
            mode="password"
            field="config.password"
            label={MetadataClusterView.password}
          />
          <Buttons />
        </Form>
      </Card>
    </>
  );
};
