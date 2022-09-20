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
import * as monaco from "monaco-editor";
import { Theme } from "@src/constants";
import { getObject, setObjectValue } from "@src/utils";
import * as _ from "lodash-es";
import editorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker";
import jsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker";

//@ts-ignore
self.MonacoEnvironment = {
  getWorker(_: any, label: any) {
    if (label === "json") {
      return new jsonWorker();
    }
    return new editorWorker();
  },
};

monaco.editor.defineTheme("lindb-light", {
  base: "vs",
  inherit: true,
  rules: [
    { token: "string.sql", foreground: "ce9178" },
    // { token: "identifier.sql", foreground: "ce9178" },
  ],
  colors: {
    // "editor.foreground": "#f38518", #010f17 identifier
    // "editor.background": "#021627",
    // "editor.lineHighlight": "#f38518",
    // "editor.lineHighlightBackground": "#010f17",
  },
});

monaco.editor.defineTheme("lindb-dark", {
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

enum StorageType {
  ui = "LINDB_UI",
}

export const UIContext = React.createContext({
  theme: Theme.dark,
  collapsed: false,
  toggleTheme: () => {},
  toggleCollapse: () => {},
  isDark: (): boolean => {
    return true;
  },
});

const UIContextProvider: React.FC = (props) => {
  const { children } = props;
  const [theme, setTheme] = useState(Theme.dark);
  const [collapsed, setCollapsed] = useState(false);

  const isDark = (): boolean => {
    return theme === Theme.dark;
  };

  useEffect(() => {
    // init ui setting
    const localUISetting = getObject(StorageType.ui);
    setTheme(_.get(localUISetting, "theme", Theme.dark));
    setCollapsed(_.get(localUISetting, "sidebarCollapsed", false));
  }, []);

  useEffect(() => {
    setObjectValue(StorageType.ui, "theme", theme);
    if (theme === Theme.dark) {
      document.body.setAttribute("theme-mode", "dark");
    } else {
      document.body.removeAttribute("theme-mode");
    }
    // set editor theme
    monaco.editor.setTheme(isDark() ? "lindb-dark" : "lindb-light");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [theme]);

  useEffect(() => {
    setObjectValue(StorageType.ui, "sidebarCollapsed", collapsed);
  }, [collapsed]);

  const handleToggleTheme = () => {
    setTheme((t) => {
      switch (t) {
        case Theme.dark:
          return Theme.light;
        case Theme.light:
        default:
          return Theme.dark;
      }
    });
  };

  const handleToggleCollapsed = () => {
    setCollapsed(!collapsed);
  };

  return (
    <UIContext.Provider
      value={{
        theme,
        collapsed,
        isDark,
        toggleTheme: handleToggleTheme,
        toggleCollapse: handleToggleCollapsed,
      }}
    >
      {children}
    </UIContext.Provider>
  );
};

export default UIContextProvider;
