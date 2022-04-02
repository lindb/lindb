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
import { ChartTooltip } from "@src/components";
import { Console } from "@src/pages";
import { URLStore } from "@src/stores";
import React from "react";
import { Route, Switch, useHistory } from "react-router-dom";
import * as monaco from "monaco-editor";

monaco.editor.defineTheme("lindb", {
  base: "vs-dark",
  inherit: true,
  rules: [
    { token: "string.sql", foreground: "ce9178" },
    // { token: "identifier.sql", foreground: "ce9178" },
  ],
  colors: {
    // "editor.foreground": "#f38518", #010f17 identifier
    "editor.background": "#021627",
    "editor.lineHighlight": "#f38518",
    "editor.lineHighlightBackground": "#010f17",
  },
});

export default function App() {
  const history = useHistory();
  // register global history in URLStore, all history operators need use URLStore.
  URLStore.setHistory(history);

  return (
    <>
      <Switch>
        <Route path="/" component={Console} />
      </Switch>
      <ChartTooltip />
    </>
  );
}
