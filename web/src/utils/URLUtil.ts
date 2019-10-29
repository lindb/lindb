import * as H from 'history'
import * as R from "ramda";

export const history = H.createHashHistory()

export function redirectTo(pathname: string, method: string = 'push') {
  const search = (new URLSearchParams(window.location.search)).toString()

  method === 'push' && history.push({ search, pathname })
  method === 'replace' && history.replace({ search, pathname })
}

export function getQueryValueOf(key: string) {
  const search = new URLSearchParams(history.location.search.split('?')[1])
  return search.get(key)
}

export function isEmpty(value: any) { // tslint:disable-line
  return R.isNil(value) || R.isEmpty(value);
}