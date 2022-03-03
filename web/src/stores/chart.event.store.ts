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
import { makeObservable, observable, action } from "mobx";
import { MouseMoveEvent } from "@src/models";

class ChartEventStore {
  public mouseMoveEvent: MouseMoveEvent | null = null;
  public mouseLeaveEvent: any = null;
  public showTooltip: boolean = false;

  constructor() {
    makeObservable(this, {
      mouseMoveEvent: observable,
      mouseLeaveEvent: observable,
      showTooltip: observable,
      mouseMove: action,
      mouseLeave: action,
      setShowTooltip: action,
    });
  }

  setShowTooltip(flag: boolean) {
    this.showTooltip = flag;
  }
  mouseMove(e: MouseMoveEvent) {
    this.mouseMoveEvent = e;
  }

  mouseLeave(e: any) {
    this.mouseLeaveEvent = e;
  }
}

export default new ChartEventStore();
