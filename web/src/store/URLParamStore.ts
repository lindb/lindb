import { action, observable } from 'mobx'
import * as R from 'ramda'
import { isEmpty } from '../utils/URLUtil'

const createHashHistory = require('history').createHashHistory
const hashHistory = createHashHistory()

export class URLParamStore {
    @observable public changed: boolean = false;
    @observable public forceChanged: boolean = false;
    @observable private params: URLSearchParams = new URLSearchParams();

    constructor() { // tslint:disable-line
        this.applyURLChange();
        window.addEventListener('hashchange', () => {
            this.applyURLChange();
        });
        window.addEventListener('popstate', () => {
            this.applyURLChange();
        });
    }

    public getValue(key: string): string {
        const value = this.params.get(key);
        if (value) {
            return value
        }
        return ''
    }

    public getValues(key: string): string[] {
        return this.params.getAll(key);
    }

    @action
    public forceChange(): void {
        this.forceChanged = !this.forceChanged;
    }

    public updateSearchParams(searchParams: URLSearchParams, params: { [key: string]: any }) {
        for (let k of Object.keys(params)) {
            const v = params[`${k}`];
            if (k) {
                if (!isEmpty(v)) {
                    if (Array.isArray(v)) {
                        searchParams.delete(k);
                        v.forEach(oneValue => searchParams.append(k, oneValue));
                    } else {
                        searchParams.set(k, v);
                    }
                }
            }
        }
    }

    public getHashSearch() {
        let search = hashHistory.location.search;
        return new URLSearchParams(search);
    }

    public getHashPathName() {
        return hashHistory.location.pathname;
    }

    public clearAll() {
        this.changeURLParams({}, [], true)
    }

    @action
    public changeURLParams(params: { [key: string]: any }, needDelete: Array<string> = [],
        clearAll: boolean = false, clearTime: boolean = false): void {
        const { hash } = window.location;
        const oldSearchParams = this.getSearchParams();
        const searchParams = clearAll ? new URLSearchParams() : this.getSearchParams();
        let pathname = hash
        if (hash.startsWith('#')) {
            pathname = hash.substring(1, hash.length)
        }
        if (pathname.indexOf('?') > -1) {
            pathname = pathname.split('?')[0]
        }

        if (!clearAll) {
            for (const key of needDelete) {
                searchParams.delete(key);
                searchParams.getAll(key);
            }
        } else {
            if (!clearTime) {
                if (oldSearchParams.has('from')) {
                    searchParams.set('from', oldSearchParams.get('from')!);
                }
                if (oldSearchParams.has('to')) {
                    searchParams.set('to', oldSearchParams.get('to')!);
                }
            }
        }
        this.updateSearchParams(searchParams, params);
        // Because of Hash history cannot PUSH the same path so delete the logic of checking path consistency
        if (!R.equals(oldSearchParams.toString(), searchParams.toString())) {
            hashHistory.push({ pathname: pathname, search: searchParams.toString() });
        }
    }

    @action
    public applyURLChange(): void {
        this.params = this.getSearchParams()
        this.changed = !this.changed;
    }

    private getSearchParams(): URLSearchParams {
        if (window.location.href.indexOf('?') > -1) {
            return new URLSearchParams(window.location.href.split('?')[1]);
        } else {
            return new URLSearchParams();
        }
    }

}