import {get} from "lodash";
import {action, observable, reaction, toJS} from "mobx";
import {Chart, ChartStatus, ChartStatusEnum, Target} from "model/Chart";
import {ResultSet} from "model/Metric";
import * as R from 'ramda';
import * as LinDBService from "service/LinQLService";
import {URLParamStore} from "store/URLParamStore";
import * as ProcessChartData from "utils/ProcessChartData";

export class ChartStore {
    urlParamStore: URLParamStore;

    charts: Map<string, Chart> = new Map(); // for chart config
    chartTrackingMap: Map<string, Chart> = new Map(); // tracking chart config change 
    seriesCache: Map<string, any> = new Map(); // for chart series data
    statsCache: Map<string, any> = new Map(); // for explain stats data
    selectedSeries: Map<string, Map<string, string>> = new Map(); // for chart series selected

    @observable chartStatusMap: Map<string, ChartStatus> = new Map(); // observe chart status

    constructor(urlParamStore: URLParamStore) {
        this.urlParamStore = urlParamStore

        // listen url params if changed
        reaction(
            () => this.urlParamStore.changed,
            () => {
                this.load(true)
            });
        reaction(
            () => this.urlParamStore.forceChanged,
            () => {
                this.load(true)
            });
    }

    public load(forceLoad?: boolean) {
        setTimeout(
            () => {
                this.charts.forEach((v: ChartStatus, uniqueId: string) => {
                    this.loadChartData(uniqueId, forceLoad)
                });
            },
            0
        );
    }

    @action
    public register(uniqueId: string, chart: Chart) {
        if (chart) {
            // for react component register too many times, when state change
            if (this.charts.has(uniqueId)) {
                return
            }
            this.add(uniqueId, chart)
        }
    }

    @action
    public reRegister(uniqueId: string, chart: Chart) {
        if (chart) {
            this.add(uniqueId, chart);
        }
    }

    @action
    public unRegister(uniqueId: string) {
        this.charts.delete(uniqueId);
        this.seriesCache.delete(uniqueId);
        this.chartTrackingMap.delete(uniqueId);
        this.chartStatusMap.delete(uniqueId);
        this.statsCache.delete(uniqueId);
        this.selectedSeries.delete(uniqueId);
    }

    public add(uniqueId: string, chart: Chart) {
        this.charts.set(uniqueId, chart)
        // copy chart data for tracking
        this.chartTrackingMap.set(uniqueId, R.clone(chart))
        this.chartStatusMap.set(uniqueId, {status: ChartStatusEnum.Init})
    }

    @action
    public setChartStatus(uniqueId: string, chartStatus: ChartStatus) {
        this.chartStatusMap.set(uniqueId, toJS(chartStatus))
    }

    public loadChartData(uniqueId: string, forceLoad: boolean = false) {
        const chartStatus: ChartStatus | undefined = this.chartStatusMap.get(uniqueId)
        if (chartStatus && chartStatus.status === ChartStatusEnum.Loading) {
            return
        }

        const chart = this.charts.get(uniqueId);
        const previousChart: Chart | undefined = this.chartTrackingMap.get(uniqueId)

        // create new targets for http request
        const target: Target | undefined = chart!.target
        chartStatus!.msg = ""

        if (forceLoad || !R.equals(target, previousChart!.target)) {
            chartStatus!.status = ChartStatusEnum.Loading
            this.setChartStatus(uniqueId, chartStatus!)
            LinDBService.query({db: target!.db, sql: this.buildQL(target!.ql)}).then(response => {
                const series: ResultSet | undefined = response.data
                const selectedSeries = this.selectedSeries.get(uniqueId);

                let reportData: any = ProcessChartData.LineChart(series!, chart!, selectedSeries!)
                this.seriesCache.set(uniqueId, reportData)
                const dataSet = get(reportData, "datasets", [])
                if (dataSet.length > 0) {
                    const limit = 50;
                    if (dataSet.length >= limit) {
                        chartStatus!.status = ChartStatusEnum.UnLimit
                    } else {
                        chartStatus!.status = ChartStatusEnum.Loaded
                    }
                } else {
                    chartStatus!.status = ChartStatusEnum.NoData
                }
                this.statsCache.set(uniqueId, series!.stats)
                this.setChartStatus(uniqueId, chartStatus!)
            }).catch((err) => {
                chartStatus!.status = ChartStatusEnum.LoadError
                chartStatus!.msg = (err.response && err.response.data) || err.message
                this.setChartStatus(uniqueId, chartStatus!)
            });
        } else {
            this.seriesCache.delete(uniqueId);
        }

        // set new target for chart config 
        previousChart!.target = target;
    }

    private buildQL(ql: string | undefined) {
        if (ql === undefined) {
            return ''
        }
        let queries: string[] = []
        this.urlParamStore.getHashSearch().forEach((value: string, key: string) => {
            queries.push(key + '="' + value + '" ')
        })
        let queryQL = queries.join(' and ')
        const whereAt = ql.indexOf('where ')
        if (whereAt < 0) {
            if (queries.length === 0) {
                return ql
            } else {
                return ql + ' where ' + queryQL
            }
        }
        if (queries.length > 0) {
            queryQL += ' and '
        }
        return ql.slice(0, whereAt + 6) + queryQL + ql.slice(whereAt + 6, ql.length)
    }

}