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
import React, { useState } from "react";
export const PlatformStateContext = React.createContext({
  chartMouseEvent: null,
  setChartMouseEvent: (_e: any) => {},
});

const PlatformStateContextProvider: React.FC<{ children: React.ReactNode }> = (
  props
) => {
  const { children } = props;
  const [chartMouseEvent, setChartMouseEvent] = useState(null);
  return (
    <PlatformStateContext.Provider
      value={{ chartMouseEvent, setChartMouseEvent }}
    >
      {children}
    </PlatformStateContext.Provider>
  );
};

export default PlatformStateContextProvider;
