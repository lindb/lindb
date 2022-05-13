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
import { Card, Form, Empty } from "@douyinfe/semi-ui";
import {
  IllustrationIdle,
  IllustrationIdleDark,
} from "@douyinfe/semi-illustrations";
import { DatabaseView, ReplicaView } from "@src/components";
import { SQL } from "@src/constants";
import { StorageState } from "@src/models";
import { exec } from "@src/services";
import { getDatabaseList } from "@src/utils";
import * as _ from "lodash-es";
import React, {
  MutableRefObject,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";
import { URLStore } from "@src/stores";

export default function ReplicationView() {
  const formApi = useRef() as MutableRefObject<any>;
  const [databases, setDatabases] = useState<any[]>([]);
  const [selectDatabase, setSelectDatabase] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const getStorageList = useCallback(async () => {
    setLoading(true);
    try {
      const storages = await exec<StorageState[]>({
        sql: SQL.ShowStorageAliveNodes,
      });
      const databases = getDatabaseList(storages);
      setDatabases(databases);
      const databaseName = URLStore.params.get("db");
      if (databaseName) {
        const db = _.find(databases, (item: any) => {
          return item.name == databaseName;
        });
        setSelectDatabase(db);
      }
    } catch (err) {
      // setError(err?.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    getStorageList();
  }, [getStorageList]);

  return (
    <>
      <Card style={{ marginBottom: 12 }} bodyStyle={{ padding: 12 }}>
        <Form
          style={{ paddingTop: 0, paddingBottom: 0 }}
          wrapperCol={{ span: 20 }}
          getFormApi={(api) => (formApi.current = api)}
          layout="horizontal"
        >
          <Form.Select
            field="database"
            label="Database"
            labelPosition="inset"
            initValue={URLStore.params.get("db")}
            optionList={_.map(databases, (db: any) => {
              return { label: db.name, value: db.name };
            })}
            onChange={(value) => {
              const db = _.find(databases, (item: any) => {
                return item.name == value;
              });
              setSelectDatabase(db);
              formApi.current.setValue("database", db.name);
              URLStore.changeURLParams({
                params: { db: value },
              });
            }}
          />
        </Form>
      </Card>
      {!selectDatabase ? (
        <Card loading={loading} bodyStyle={{ padding: 12 }}>
          <Empty
            image={<IllustrationIdle style={{ width: 150, height: 150 }} />}
            darkModeImage={
              <IllustrationIdleDark style={{ width: 150, height: 150 }} />
            }
            description="Please Select Database"
            style={{ marginTop: 50, minHeight: 400 }}
          />
        </Card>
      ) : (
        <>
          <DatabaseView
            liveNodes={_.get(selectDatabase, "storage.liveNodes", {})}
            storage={_.get(selectDatabase, "storage", {})}
            loading={false}
            databaseName={selectDatabase.name}
          />
          <div style={{ marginTop: 12 }}>
            <ReplicaView
              liveNodes={_.get(selectDatabase, "storage.liveNodes", {})}
              db={selectDatabase.name as string}
              storage={selectDatabase.storage.name as string}
            />
          </div>
        </>
      )}
    </>
  );
}
