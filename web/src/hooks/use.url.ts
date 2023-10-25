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
  MutableRefObject,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";
import { reaction } from "mobx";
import { URLStore } from "@src/stores";
import * as _ from "lodash-es";

function useWatchURLChange<T extends (...args: any[]) => any>(
  callback: T
): void {
  useEffect(() => {
    const disposer = [
      reaction(
        () => URLStore.changed, // watch url change event
        () => {
          callback(); // set value after url changed
        },
        {
          fireImmediately: true,
        }
      ),
    ];
    return () => {
      disposer.forEach((d) => d());
    };
    // NOTE: Don't add callback into deps
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
}

const getParamValues = (keys: string[]) => {
  let newParams: any = URLStore.getParams();
  if (!_.isEmpty(keys)) {
    newParams = _.pick(newParams, keys);
  }
  return newParams || {};
};

function useParams(keys?: string[]) {
  const previous = useRef({}) as MutableRefObject<any>;
  const [params, setParams] = useState<any>(() => {
    const values = getParamValues(keys || []);
    previous.current = values;
    return values;
  });

  const buildParams = useCallback(() => {
    const values = getParamValues(keys || []);
    if (!_.isEqual(values, previous.current)) {
      previous.current = values;
      setParams(values);
    }
    // NOTE: Don't add keys into deps
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    const disposer = reaction(
      () => URLStore.changed, // watch params change
      () => {
        buildParams();
      },
      { fireImmediately: true }
    );
    return () => {
      disposer();
    };
  }, [buildParams]);

  // NOTE: don't use {...params}
  return params;
}

export { useWatchURLChange, useParams };
