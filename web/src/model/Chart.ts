import { isEmpty } from "../utils/URLUtil";
import { ResultSet } from "./Metric";

export class Chart {
    loading?: boolean = false;
    title?: string;
    description?: string;
    config?: any;
    from?: string;
    to?: string;
    timeShift?: string;
    target?: Target;
    series?: Array<ResultSet>;
}


export class Target {
    db?: string;
    ql?: string;
    // check target is valid
    public static valid(target: Target): boolean {
        if (isEmpty(target.db) || isEmpty(target.ql)) {
            return false;
        }
        return true;
    }
}

export enum ChartStatusEnum {
    Init = "init",
    Loading = "loading",
    Loaded = "loaded",
    BadRequest = "badRequest",
    NoData = "noData",
    LoadError = "loadError",
    UnMount = "unMount",
    UnLimit = "unLimit"
}

export class ChartStatus {
    loading?: boolean;
    status?: ChartStatusEnum;
    msg?: string;

    constructor() {
        this.status = ChartStatusEnum.Init;
    }
}