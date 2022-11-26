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
import React, {
  useEffect,
  useState,
  useRef,
  MutableRefObject,
  useContext,
} from "react";
import { Popover, Form, Button, SplitButtonGroup } from "@douyinfe/semi-ui";
import { IconFilter } from "@douyinfe/semi-icons";
import { TagValueSelect } from "@src/components";
import { ExecService } from "@src/services";
import { Metadata } from "@src/models";
import { URLStore } from "@src/stores";
import { UIContext } from "@src/context/UIContextProvider";

export default function TagFilterSelect(props: {
  db: string;
  namespace?: string;
  metric: string;
}) {
  const formApi = useRef() as MutableRefObject<any>;
  const { db, namespace, metric } = props;
  const [tagKeys, setTagKeys] = useState<string[]>([]);
  const [visible, setVisible] = useState(false);
  const { locale } = useContext(UIContext);
  const { Common } = locale;

  useEffect(() => {
    const fetchTagKeys = async () => {
      const metadata = await ExecService.exec<Metadata>({
        db: db,
        sql: `show tag keys from '${metric}'`,
      });
      const tagKeys = (metadata as Metadata).values || [];
      setTagKeys(tagKeys as string[]);
    };
    fetchTagKeys();
  }, [db, metric]);

  return (
    <Popover
      trigger="custom"
      visible={visible}
      showArrow
      content={
        <>
          <Form
            getFormApi={(api: any) => (formApi.current = api)}
            className="lin-tag-filter"
          >
            {tagKeys.map((tagKey: string) => (
              <div key={tagKey} style={{ marginBottom: 4, width: 400 }}>
                <TagValueSelect
                  style={{ width: "100%" }}
                  variate={{
                    db: db,
                    tagKey: tagKey,
                    label: tagKey,
                    sql: `show tag values from '${metric}' with key='${tagKey}'`,
                  }}
                  labelPosition="inset"
                />
              </div>
            ))}
          </Form>
          <SplitButtonGroup
            style={{ marginTop: 4, width: "100%", textAlign: "right" }}
          >
            <Button type="tertiary" onClick={() => setVisible(false)}>
              {Common.cancel}
            </Button>
            <Button
              type="secondary"
              onClick={() => {
                URLStore.changeURLParams({
                  params: { tags: JSON.stringify(formApi.current.getValues()) },
                });
                setVisible(false);
              }}
            >
              {Common.ok}
            </Button>
          </SplitButtonGroup>
        </>
      }
    >
      <IconFilter
        style={{ cursor: "pointer", width: 32 }}
        onClick={() => setVisible(true)}
      />
    </Popover>
  );
}
