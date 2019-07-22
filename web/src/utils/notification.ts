import { notification } from 'antd'
import { AxiosResponse } from 'axios'

function errorHandler(props: any) {
  notification.error(props)
}

function successHandler(props: any) {
  notification.success(props)
}

function infoHandler(props: any) {
  notification.info(props)
}

function warningHandler(props: any) {
  notification.warning(props)
}

function httpCodeHandler(resp: AxiosResponse, url: string, message: string) {
  if (!resp) {
    return errorHandler({
      message: message,
      description: resp,
    })
  }
  if (!resp || (resp.status < 200 || resp.status >= 300)) {
    if (resp.status == 401) {
      let callbackUrl = window.location.hash
      window.location.hash = `/login?from=${encodeURIComponent(callbackUrl)}`
    } else {
      return errorHandler({
        message: message,
        description: resp.status + ':' + resp.data.error || (resp.data.stat && resp.data.stat.errorMsg),
      })
    }
  } else {
    successHandler({
      message: message,
      description: 'success',
    })
  }
}

export { errorHandler, successHandler, infoHandler, warningHandler, httpCodeHandler }
