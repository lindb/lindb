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
import { Theme, StorageType, Language } from "@src/constants";
import { LocalStorageKit } from "@src/utils";
import * as _ from "lodash-es";
import editorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker";
import jsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker";
import en_US from "@douyinfe/semi-ui/lib/es/locale/source/en_US";
import zh_CN from "@douyinfe/semi-ui/lib/es/locale/source/zh_CN";
import { en_US as lin_en_US, zh_CN as lin_zh_CN } from "@src/i18n";
import { useQuery } from "@tanstack/react-query";
import { PlatformService } from "@src/services";
import { Spin } from "@douyinfe/semi-ui";

const localeMap = {
  zh_CN: _.merge(zh_CN, lin_zh_CN),
  en_US: _.merge(en_US, lin_en_US),
};
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
    // "editor.foreground": "#f5f5f5", // #010f17 identifier
    "editor.background": "#f5f5f5",
    "editor.lineHighlight": "#f38518",
    "editor.lineHighlightBackground": "#eee",
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

export const UIContext = React.createContext({
  theme: Theme.dark,
  collapsed: false,
  locale: {} as any,
  env: {} as any,
  language: Language.en_US,
  toggleTheme: () => {},
  toggleCollapse: () => {},
  toggleLangauge: () => {},
  isDark: (): boolean => {
    return true;
  },
});

const UIContextProvider: React.FC<{ children: React.ReactNode }> = (props) => {
  const { children } = props;
  const localUISetting = LocalStorageKit.getObject(StorageType.ui);
  const [theme, setTheme] = useState(
    _.get(localUISetting, "theme", Theme.dark)
  );
  const defaultLanguage = _.get(localUISetting, "language", Language.en_US);
  const [language, setLanguage] = useState(defaultLanguage);
  const [locale, setLocale] = useState<any>(localeMap[defaultLanguage]);
  const [collapsed, setCollapsed] = useState(
    _.get(localUISetting, "sidebarCollapsed", false)
  );
  const [env, setEnv] = useState<any>({});
  const { isLoading } = useQuery(["load-env"], async () => {
    return PlatformService.fetchEnv().then((data) => {
      setEnv(data);
    });
  });
  const isDark = (): boolean => {
    return theme === Theme.dark;
  };

  useEffect(() => {
    LocalStorageKit.setObjectValue(StorageType.ui, "theme", theme);
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
    LocalStorageKit.setObjectValue(StorageType.ui, "language", language);
    setLocale(localeMap[language]);
  }, [language]);

  useEffect(() => {
    LocalStorageKit.setObjectValue(
      StorageType.ui,
      "sidebarCollapsed",
      collapsed
    );
  }, [collapsed]);

  const handleToggleTheme = () => {
    setTheme((t: Theme) => {
      switch (t) {
        case Theme.dark:
          return Theme.light;
        case Theme.light:
        default:
          return Theme.dark;
      }
    });
  };

  const handleToggleLanguage = () => {
    setLanguage(() => {
      switch (language) {
        case Language.en_US:
          return Language.zh_CN;
        case Language.zh_CN:
          return Language.en_US;
        default:
          return Language.en_US;
      }
    });
  };

  const handleToggleCollapsed = () => {
    setCollapsed(!collapsed);
  };

  const renderContent = () => {
    if (isLoading) {
      return (
        <div style={{ width: "100%", textAlign: "center", marginTop: 300 }}>
          <Spin size="large" />
        </div>
      );
    }
    return children;
  };

  return (
    <UIContext.Provider
      value={{
        theme,
        collapsed,
        isDark,
        locale,
        language,
        env,
        toggleLangauge: handleToggleLanguage,
        toggleTheme: handleToggleTheme,
        toggleCollapse: handleToggleCollapsed,
      }}
    >
      {renderContent()}
    </UIContext.Provider>
  );
};

export default UIContextProvider;
