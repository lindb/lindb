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
import * as React from "react";
import { createRoot } from "react-dom/client";
import App from "@src/App";
import "@src/styles/index.scss";
import { LocaleProvider } from "@douyinfe/semi-ui";
import { HashRouter as Router, Route, Switch } from "react-router-dom";
import { UIContextProvider } from "@src/context";
import { QueryClientProvider, QueryClient } from "@tanstack/react-query";
import * as _ from "lodash-es";
import { UIContext } from "./context/UIContextProvider";
import { useContext } from "react";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
      cacheTime: 0,
    },
  },
});

const container = document.getElementById("root");
const root = createRoot(container as any);

const AppPage: React.FC = () => {
  const { locale } = useContext(UIContext);
  return (
    <LocaleProvider locale={locale}>
      <QueryClientProvider client={queryClient}>
        <Router>
          <Switch>
            <Route path="/" component={App} />
          </Switch>
        </Router>
      </QueryClientProvider>
    </LocaleProvider>
  );
};

root.render(
  <UIContextProvider>
    <AppPage />
  </UIContextProvider>
);
