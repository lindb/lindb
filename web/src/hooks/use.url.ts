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
import { useEffect } from "react";
import { reaction } from "mobx";
import { URLStore } from "@src/stores";

export function useWatchURLChange<T extends (...args: any[]) => any>(
  callback: T
): void {
  useEffect(() => {
    callback(); // init values

    const disposer = [
      reaction(
        () => URLStore.changed, // watch url change event
        () => {
          callback(); // set value after url changed
        }
      ),
    ];
    return () => {
      disposer.forEach((d) => d());
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
}
