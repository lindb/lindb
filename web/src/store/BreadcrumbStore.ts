/**
 * ---------------------------------------------------------------------------------------------------------------
 * the only data source of BreadcrumbHeader component
 * @example breadcrumbList => [{label: 'Setting',path: '/setting'}, {label:'Database', path: '/setting/database'}]
 * ---------------------------------------------------------------------------------------------------------------
 */
import {observable} from 'mobx'
import {BreadcrumbStatus} from '../model/Breadcrumb'

export class BreadcrumbStore {
    @observable public breadcrumbList: Array<BreadcrumbStatus> = []

    public setBreadcrumbs(breadcrumbs: Array<BreadcrumbStatus>) {
        this.breadcrumbList = breadcrumbs
    }

}