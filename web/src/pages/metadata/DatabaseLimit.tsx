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
import { IconClose, IconSaveStroked } from "@douyinfe/semi-icons";
import { Banner, Button, Card, Col, Row } from "@douyinfe/semi-ui";
import { Route } from "@src/constants";
import { UIContext } from "@src/context/UIContextProvider";
import { useParams } from "@src/hooks";
import { ExecService } from "@src/services";
import { URLStore } from "@src/stores";
import { useQuery } from "@tanstack/react-query";
import * as _ from "lodash-es";
import * as monaco from "monaco-editor";
import React, {
  MutableRefObject,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";

export default function DatabaseLimits() {
  const editorRef = useRef() as MutableRefObject<HTMLDivElement>;
  const editor = useRef() as MutableRefObject<any>;
  const [saveError, setSaveError] = useState("");
  const { db } = useParams(["db"]);
  const { isError, error, data, isLoading } = useQuery(
    ["show_limits", db],
    async () => {
      return ExecService.exec<string>({
        sql: "show limit",
        db: db,
      });
    },
    {
      enabled: !_.isEmpty(db),
    }
  );
  const [submiting, setSubmiting] = useState(false);
  const { locale } = useContext(UIContext);
  const { Common } = locale;

  useEffect(() => {
    if (isLoading || isError || !editorRef.current) {
      return;
    }
    editor.current = monaco.editor.create(editorRef.current, {
      language: "ini",
      minimap: { enabled: false },
      value: data,
    });
  }, [isLoading, isError, editorRef, data]);

  const setLimit = async () => {
    try {
      setSaveError("");
      setSubmiting(true);
      await ExecService.exec({
        sql: `set limit '${editor.current.getValue()}'`,
        db: db,
      });
      URLStore.changeURLParams({ path: Route.MetadataDatabase });
    } catch (err) {
      setSaveError(_.get(err, "response.data", Common.unknownInternalError));
    } finally {
      setSubmiting(false);
    }
  };

  return (
    <>
      {(error || saveError) && (
        <Banner
          description={`${error}` || saveError}
          type="danger"
          closeIcon={null}
          style={{ marginBottom: 12, justifyContent: "left" }}
        />
      )}
      <Card style={{ height: "85vh" }} bodyStyle={{ height: "100%" }}>
        <Row style={{ paddingBottom: 12 }}>
          <Col>
            <Button
              type="secondary"
              icon={<IconSaveStroked />}
              style={{ marginRight: 8 }}
              loading={submiting}
              onClick={() => {
                setLimit();
              }}
            >
              {Common.save}
            </Button>
            <Button
              type="tertiary"
              icon={<IconClose />}
              onClick={() =>
                URLStore.changeURLParams({ path: Route.MetadataDatabase })
              }
            >
              {Common.cancel}
            </Button>
          </Col>
        </Row>
        <div ref={editorRef} style={{ height: "100%" }} />
      </Card>
    </>
  );
}
