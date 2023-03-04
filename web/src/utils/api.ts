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
import axios from "axios";
import qs from "qs";
import * as _ from "lodash-es";
import JSONbig from "json-bigint";

// env control
switch (import.meta.env.MODE) {
  case "development":
    axios.defaults.baseURL = "http://localhost:9000/api/v1";
    break;
  case "production":
  default:
    axios.defaults.baseURL = "/api/v1";
    break;
}
axios.defaults.timeout = 60000; // set timeout
axios.defaults.headers.common["Content-Type"] = "application/json";
axios.defaults.transformResponse = [
  (data) => {
    if (typeof data === "string") {
      try {
        data = JSONbig.parse(data);
        return data;
      } catch (e) {
        /* Ignore */
      } // Added this Ignore as it's the same in the Axios
    }
    return data;
  },
];

async function GET<T>(
  url: string,
  params?: { [index: string]: any } | undefined
): Promise<T> {
  const target =
    url + (params ? `?${qs.stringify(params, { arrayFormat: "repeat" })}` : "");
  return axios
    .get<T>(target)
    .then((result) => {
      return Promise.resolve(result.data);
    })
    .catch((err) => {
      return Promise.reject(err);
    });
}

async function POST<T>(
  url: string,
  params?: { [index: string]: any } | undefined
): Promise<T> {
  return axios
    .post<T>(url, params)
    .then((result) => {
      return Promise.resolve(result.data);
    })
    .catch((err) => {
      return Promise.reject(err);
    });
}

const getErrorMsg = (err: any) => {
  if (_.has(err, "response.data")) {
    return _.get(err, "response.data");
  }
  if (_.has(err, "reason.response.data")) {
    return _.get(err, "reason.response.data");
  }
  const msg = _.get(err, "reason", "Unknown internal error");
  return `${msg}`;
};

const getErrorCode = (err: any) => {
  return _.get(err, "response.status", 0);
};

export default {
  getErrorMsg,
  getErrorCode,
  GET,
  POST,
};
