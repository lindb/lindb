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
  ArrayField,
  Banner,
  Button,
  Card,
  Col,
  Form,
  Row,
} from "@douyinfe/semi-ui";
import {
  IconClose,
  IconMinusCircle,
  IconPlusCircle,
  IconSaveStroked,
} from "@douyinfe/semi-icons";
import { UIContext } from "@src/context/UIContextProvider";
import { Broker } from "@src/models";
import { ExecService } from "@src/services";
import React, {
  MutableRefObject,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import * as _ from "lodash-es";
import { URLStore } from "@src/stores";
import { Route } from "@src/constants";

export const LogicDatabaseConfig: React.FC = () => {
  const formApi = useRef() as MutableRefObject<any>;
  const { locale } = useContext(UIContext);
  const { MetadataLogicDatabaseView, Common } = locale;
  const [submiting, setSubmiting] = useState(false);
  const [brokerList, setBrokerList] = useState([] as any[]);
  const [error, setError] = useState("");

  useEffect(() => {
    const getBrokerList = async () => {
      try {
        const list = await ExecService.exec<Broker[]>({
          sql: "show brokers",
        });
        const selectList: any[] = [];
        _.forEach(list || [], (s) => {
          const ns = _.get(s, "config.namespace");
          selectList.push({ value: ns, label: ns });
        });
        setBrokerList(selectList);
      } catch (err: any) {
        setError(err?.message);
        setBrokerList([]);
      }
    };
    getBrokerList();
  }, []);

  const create = async () => {
    if (!formApi.current) {
      return;
    }
    const createDatabase = async (values: any) => {
      try {
        setSubmiting(true);
        await ExecService.exec({
          sql: `create database ${JSON.stringify(values)}`,
        });
        URLStore.changeURLParams({ path: Route.MetadataLogicDatabase });
      } catch (err) {
        setError(_.get(err, "response.data", Common.unknownInternalError));
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
          getFormApi={(api: any) => (formApi.current = api)}
          className="lin-db-form"
          labelPosition="left"
          labelAlign="right"
          labelCol={{ span: 4 }}
          wrapperCol={{ span: 12 }}
        >
          <Form.Input
            label={MetadataLogicDatabaseView.name}
            field="name"
            rules={[
              {
                required: true,
                message: MetadataLogicDatabaseView.nameRequired,
              },
            ]}
          />
          <Row>
            <Col
              span={4}
              style={{ display: "flex", justifyContent: "flex-end" }}
            >
              <Form.Label style={{ paddingRight: 16 }} required>
                {MetadataLogicDatabaseView.router}
              </Form.Label>
            </Col>
            <Col>
              <Form.Label style={{ width: 166 }} required>
                {MetadataLogicDatabaseView.tagKey}
              </Form.Label>
              <Form.Label style={{ width: 266 }} required>
                {MetadataLogicDatabaseView.tagValues}
              </Form.Label>
              <Form.Label style={{ width: 204 }} required>
                {MetadataLogicDatabaseView.brokers}
              </Form.Label>
              <Form.Label style={{ width: 167 }}>
                {MetadataLogicDatabaseView.targetDatabase}
              </Form.Label>
            </Col>
          </Row>
          <ArrayField
            field="routers"
            initValue={[{ key: "", values: [], brokers: [] }]}
          >
            {({ add, arrayFields }) => (
              <>
                {arrayFields.map((f: any, idx) => (
                  <Row key={f.key}>
                    <Col offset={4} className="lin-form-input-group">
                      <Form.InputGroup>
                        <Form.Input
                          field={`${f.field}[key]`}
                          style={{ width: 162, marginRight: 4 }}
                          noLabel
                        />
                        <Form.TagInput
                          style={{ width: 262, marginRight: 4 }}
                          field={`${f.field}[values]`}
                          noLabel
                        />
                        <Form.Select
                          style={{ width: 202, marginRight: 4 }}
                          field={`${f.field}[broker]`}
                          optionList={brokerList}
                          noLabel
                        />
                        <Form.Input
                          field={`${f.field}[database]`}
                          style={{ width: 162, marginRight: 4 }}
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
                {Common.submit}
              </Button>
              <Button
                type="tertiary"
                icon={<IconClose />}
                onClick={() =>
                  URLStore.changeURLParams({
                    path: Route.MetadataLogicDatabase,
                  })
                }
              >
                {Common.cancel}
              </Button>
            </Col>
          </Row>
        </Form>
      </Card>
    </>
  );
};

export default LogicDatabaseConfig;
