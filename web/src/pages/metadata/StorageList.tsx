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
import { IconRedoStroked } from "@douyinfe/semi-icons";
import { Button, Popconfirm, Notification } from "@douyinfe/semi-ui";
import { Route } from "@src/constants";
import { ExecService } from "@src/services";
import * as _ from "lodash-es";
import React, { useContext } from "react";
import { UIContext } from "@src/context/UIContextProvider";
import { ClusterList } from "./ClusterList";

export default function StorageList() {
  const { locale } = useContext(UIContext);
  const { MetadataStorageView, Common } = locale;

  const recover = (storageName: string) => {
    return new Promise(async (resolve: any) => {
      try {
        const msg: any = await ExecService.exec({
          sql: `recover storage '${storageName}'`,
        });
        Notification.success({
          title: MetadataStorageView.recoverSuccessTitle,
          content: _.join(msg || [], ","),
          position: "top",
          theme: "light",
          duration: 5,
        });
        resolve();
      } catch (err) {
        Notification.error({
          title: MetadataStorageView.recoverErrorTitle,
          content: _.get(err, "response.data", Common.unknownInternalError),
          position: "top",
          theme: "light",
          duration: 5,
        });
        resolve();
      }
    });
  };
  return (
    <ClusterList
      sql="show storages"
      registerPath={Route.MetadataStorageConfig}
      actions={{
        title: Common.actions,
        key: "actions",
        width: 100,
        render: (_text: any, r: any) => {
          return (
            <Popconfirm
              title={Common.pleaseConfirm}
              content={MetadataStorageView.recoverConfirmMessage}
              onConfirm={() => recover(_.get(r, "config.namespace"))}
            >
              <Button icon={<IconRedoStroked />} type="primary" />
            </Popconfirm>
          );
        },
      }}
    />
  );
}
