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
  IconMinusCircle,
  IconPlusCircle,
  IconSaveStroked,
} from "@douyinfe/semi-icons";
import {
  ArrayField,
  Banner,
  Button,
  Card,
  Col,
  Form,
  Row,
  Typography,
} from "@douyinfe/semi-ui";
import { Route } from "@src/constants";
import { Storage } from "@src/models";
import { ExecService } from "@src/services";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";
import React, { MutableRefObject, useEffect, useRef, useState } from "react";

const Text = Typography.Text;

export default function DatabaseConfig() {
  const formApi = useRef() as MutableRefObject<any>;
  const [storageList, setStorageList] = useState([] as any[]);
  const [error, setError] = useState("");
  const [submiting, setSubmiting] = useState(false);

  useEffect(() => {
    const getStorageList = async () => {
      try {
        const list = await ExecService.exec<Storage[]>({
          sql: "show storages",
        });
        const selectList: any[] = [];
        _.forEach(list || [], (s) => {
          const ns = _.get(s, "config.namespace");
          selectList.push({ value: ns, label: ns });
        });
        setStorageList(selectList);
      } catch (err) {
        setError(err?.message);
        setStorageList([]);
      }
    };
    getStorageList();
  }, []);

  const create = async () => {
    if (!formApi.current) {
    }
    const createDatabase = async (values: any) => {
      try {
        setSubmiting(true);
        await ExecService.exec({
          sql: `create database ${JSON.stringify(values)}`,
        });
        URLStore.changeURLParams({ path: Route.MetadataDatabase });
      } catch (err) {
        setError(_.get(err, "response.data", "Unknown internal error"));
      } finally {
        setSubmiting(false);
      }
    };
    formApi.current
      .validate()
      .then((values: any) => {
        createDatabase(values);
      })
      .catch(() => {});
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
          className="lin-db-form"
          getFormApi={(api) => (formApi.current = api)}
          labelPosition="left"
          labelAlign="right"
          labelCol={{ span: 4 }}
          wrapperCol={{ span: 12 }}
        >
          <Form.Input
            label="Name"
            field="name"
            rules={[{ required: true, message: "Name is required" }]}
          />
          <Form.Select
            label="Storage"
            field="storage"
            rules={[{ required: true, message: "Storage is required" }]}
            optionList={storageList}
            style={{ width: 200 }}
          />
          <Form.InputNumber
            rules={[{ required: true, message: "Num. of Shard is required" }]}
            label="Num. Of Shard"
            field="numOfShard"
            min={1}
          />
          <Form.InputNumber
            field="replicaFactor"
            rules={[{ required: true, message: "Replica Factor is require" }]}
            label="Replica Factor"
            min={1}
          />
          <div
            style={{
              borderBottom: "1px solid var(--semi-color-border)",
              paddingTop: 12,
              paddingBottom: 12,
            }}
          >
            <Row>
              <Col
                span={4}
                style={{ display: "flex", justifyContent: "flex-end" }}
              >
                <Text strong style={{ paddingRight: 16 }}>
                  Storage Engine Options
                </Text>
              </Col>
            </Row>
          </div>
          <Form.Switch
            label="Auto Create Namespace"
            field="option.autoCreateNS"
            initValue={true}
          />
          <Row>
            <Col
              span={4}
              style={{ display: "flex", justifyContent: "flex-end" }}
            >
              <Form.Label style={{ paddingRight: 16 }} required>
                Intervals
              </Form.Label>
            </Col>
            <Col>
              <Form.Label style={{ width: 220 }} required>
                Interval(Seconds)
              </Form.Label>
              <Form.Label style={{ width: 200 }} required>
                Retention(Days)
              </Form.Label>
            </Col>
          </Row>
          <ArrayField
            field="option.intervals"
            initValue={[{ interval: "10s", retention: "30d" }]}
          >
            {({ add, arrayFields }) => (
              <>
                {arrayFields.map((f: any, idx) => (
                  <Row key={f.key}>
                    <Col offset={4} className="lin-form-input-group">
                      <Form.InputGroup>
                        <Form.Input
                          field={`${f.field}[interval]`}
                          style={{ width: 202, marginRight: 16 }}
                          noLabel
                        />
                        <Form.Input
                          style={{ width: 202, marginRight: 4 }}
                          field={`${f.field}[retention]`}
                          noLabel
                        />
                      </Form.InputGroup>
                      {arrayFields.length > 1 && (
                        <Button
                          type="danger"
                          theme="borderless"
                          onClick={f.remove}
                          icon={<IconMinusCircle />}
                        />
                      )}
                      {idx == arrayFields.length - 1 && (
                        <Button
                          type="primary"
                          theme="borderless"
                          onClick={add}
                          icon={<IconPlusCircle />}
                        />
                      )}
                    </Col>
                  </Row>
                ))}
              </>
            )}
          </ArrayField>
          <Row style={{ paddingTop: 12 }}>
            <Col
              span={4}
              style={{ display: "flex", justifyContent: "flex-end" }}
            >
              <Form.Label style={{ paddingRight: 16, paddingTop: 10 }}>
                Writeable Time Range
              </Form.Label>
            </Col>
            <Col className="lin-form-input-group">
              <Form.InputGroup>
                <Form.Input
                  label="Behead"
                  labelPosition="inset"
                  field="option.behead"
                  style={{ width: 202, marginRight: 16 }}
                  placeholder="30m/1h"
                />
                <Form.Input
                  label="Ahead"
                  labelPosition="inset"
                  field="option.ahead"
                  placeholder="30m/1h"
                  style={{ width: 202, marginRight: 16 }}
                />
              </Form.InputGroup>
              <Text size="small" type="tertiary">
                For Example: [ now()-1h ~ now()+1h ]
              </Text>
            </Col>
          </Row>
          <Form.Slot style={{ padding: 0 }}></Form.Slot>
          <Row style={{ paddingTop: 12 }}>
            <Col offset={4}>
              <Button
                type="secondary"
                icon={<IconSaveStroked />}
                style={{ marginRight: 8 }}
                loading={submiting}
                onClick={() => {
                  create();
                }}
              >
                Submit
              </Button>
              <Button
                type="tertiary"
                icon={<IconClose />}
                onClick={() =>
                  URLStore.changeURLParams({ path: Route.MetadataDatabase })
                }
              >
                Cancel
              </Button>
            </Col>
          </Row>
        </Form>
      </Card>
    </>
  );
}
