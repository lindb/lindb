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
import * as _ from "lodash-es";

const getTagConditions = (params: object, tagKeys: string[]): string[] => {
  const values = _.pick(params, tagKeys || []);
  const tags: string[] = [];
  _.forOwn(values, (objV, objK) => {
    if (_.isEmpty(objV)) {
      return;
    }
    const watchValues = _.concat([], objV);
    const tagValues: string[] = _.map(watchValues, (v: string) => `'${v}'`);
    if (!_.isEmpty(tagValues)) {
      tags.push(`'${objK}' in (${tagValues.join(",")})`);
    }
  });
  return tags;
};

export default { getTagConditions };
