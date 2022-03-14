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
import { Banner, Card, Form, Row, Col } from "@douyinfe/semi-ui";
import { Storage } from "@src/models";
import { exec } from "@src/services";
import React, { useEffect, useState } from "react";
import * as _ from "lodash-es";
export default function DatabaseConfig() {
  const [storageList, setStorageList] = useState([] as any[]);
  const [error, setError] = useState("");
  useEffect(() => {
    const getStorageList = async () => {
      try {
        const list = await exec<Storage[]>({ sql: "show storages" });
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
      <Card bordered>
        <Form
          labelPosition="left"
          labelAlign="right"
          labelCol={{ span: 4 }}
          wrapperCol={{ span: 12 }}
        >
          <Form.Input label="Name" field="name" rules={[{ required: true }]} />
          <Form.Select
            label="Storage"
            field="storeage"
            rules={[{ required: true }]}
            optionList={storageList}
            style={{ width: 200 }}
          />
          <Form.InputNumber label="Num. Of Shard" field="numOfShard" min="1" />
          <Form.InputNumber
            field="replicaFactor"
            label="Replica Factor"
            min="1"
          />
          <Form.Section
            text={
              <Row>
                <Col
                  span={4}
                  style={{ display: "flex", justifyContent: "flex-end" }}
                >
                  <div style={{ marginRight: 16 }}>TSDB Options</div>
                </Col>
              </Row>
            }
          >
            <Form.Switch
              label="Auto Create Namespace"
              field="option.autoCreateNS"
              initValue={true}
            />
          </Form.Section>
        </Form>
      </Card>
    </>
  );
}
