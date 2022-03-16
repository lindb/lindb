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
import React, { useEffect, useState } from "react";
import { Popover, Button, SplitButtonGroup } from "@douyinfe/semi-ui";
import { IconFilter } from "@douyinfe/semi-icons";
import { TagValueSelect } from "@src/components";
import { exec } from "@src/services";
import { Metadata } from "@src/models";

export default function TagFilterSelect(props: {
  db: string;
  namespace?: string;
  metric: string;
}) {
  const { db, namespace, metric } = props;
  const [tagKeys, setTagKeys] = useState<string[]>([]);
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const fetchTagKeys = async () => {
      const metadata = await exec<Metadata>({
        db: db,
        sql: `show tag keys from '${metric}'`,
      });
      const tagKeys = (metadata as Metadata).values || [];
      setTagKeys(tagKeys);
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
          {tagKeys.map((tagKey: string) => (
            <div key={tagKey} style={{ marginBottom: 4, width: 400 }}>
              <TagValueSelect
                style={{ width: "100%" }}
                variate={{
                  db: db,
                  tagKey: tagKey,
                  label: tagKey,
                  ql: `show tag values from '${metric}' with key='${tagKey}'`,
                }}
                labelPosition="inset"
              />
            </div>
          ))}
          <SplitButtonGroup
            style={{ marginTop: 4, width: "100%", textAlign: "right" }}
          >
            <Button type="tertiary" onClick={() => setVisible(false)}>
              Cancel
            </Button>
            <Button type="secondary" onClick={() => setVisible(false)}>
              OK
            </Button>
          </SplitButtonGroup>
        </>
      }
    >
      <IconFilter
        style={{ marginLeft: 8, marginRight: 8, cursor: "pointer" }}
        onClick={() => setVisible(true)}
      />
    </Popover>
  );
}
