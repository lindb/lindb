import { message } from 'antd'
import axios, { AxiosResponse } from 'axios'
import { API_URL, LOCALSTORAGE_TOKEN } from 'config/config'
import { redirectTo } from 'utils/URLUtil'

// env control
const env = process.env.NODE_ENV
const qs = require('qs')
switch (env) {
  case 'development':
  case 'production':
  default:
    axios.defaults.baseURL = API_URL
    break
}
axios.defaults.timeout = 60000 // set timeout

/* Global Message Setting */
message.config({
  top: 50,
  duration: 3,
  maxCount: 3,
})

export function get<T>(url: string, params?: { [index: string]: any } | undefined): Promise<AxiosResponse<T | undefined>> {
  const target = url + (params ? '?' + qs.stringify(params) : '')
  return axios.get<T>(target)
}

/* package */
export function GET<T>(url: string, params?: { [index: string]: any } | undefined, msg?: { success?: string, error?: string }): Promise<T | undefined> {
  const target = url + (params ? '?' + qs.stringify(params) : '')
  return axios
    .get<T>(target)
    .then((result) => {
      return Promise.resolve(processData(result, msg && msg.success))
    })
    .catch((err) => {
      message.error((msg && msg.error) || (err.response && err.response.data) || err.message)
      return Promise.reject(err)
    })
}

export function POST<T>(url: string, params?: object | undefined, msg?: { success?: string, error?: string }): Promise<T | undefined> {
  return axios
    .post<T>(url, params)
    .then((result) => {
      return Promise.resolve(processData(result, msg && msg.success))
    })
    .catch((err) => {
      message.error((msg && msg.error) || (err.response && err.response.data) || err.message)
      return Promise.reject(err)
    })
}

export function PUT<T>(url: string, params?: object | undefined, msg?: { success?: string, error?: string }): Promise<T | undefined> {
  return axios
    .put<T>(url, params)
    .then((result) => {
      return Promise.resolve(processData(result, msg && msg.success))
    })
    .catch((err) => {
      message.error((msg && msg.error) || (err.response && err.response.data) || err.message)
      return Promise.reject(err)
    })
}

export async function DELETE<T>(url: string, msg?: { success?: string, error?: string }): Promise<T | undefined> {
  return axios
    .delete(url)
    .then((result) => {
      return Promise.resolve(processData(result, msg && msg.success))
    })
    .catch((err) => {
      message.error((msg && msg.error) || (err.response && err.response.data) || err.message)
      return Promise.reject(err)
    })
}

function processData(result: any, msg?: string) {
  msg && message.success(msg)

  return result.data
}

// Request interceptor
axios.interceptors.request.use(
  async (config) => {
    /* Set request header */
    // config.headers[ 'Content-Type' ] = 'application/json'
    // config.headers[ 'x-requested-with' ] = 'ajax'

    /* Set Token */
    const token = localStorage.getItem(LOCALSTORAGE_TOKEN)
    token && (config.headers.Authorization = token)

    return config
  },
  error => {
    message.warning('Request Timeout, please try again.')
    return Promise.resolve(error)
  },
)

// Response interceptor
axios.interceptors.response.use(
  (data: any) => {
    return data
  },
  err => {
    // Error handler
    if (err && err.response) {
      switch (err.response.status) {
        case 400:
          err.message = 'Bad Request (400)'
          break
        case 401:
          err.message = 'Permission denied, please login (401)'
          const path = window.location.hash.split('#')[1]
          redirectTo(`/login${path ? `?from=${path}` : ''}`)
          break
        case 403:
          err.message = 'Forbidden (403)'
          break
        case 404:
          err.message = 'Not Found (404)'
          break
        case 408:
          err.message = 'Request Timeout (408)'
          break
        case 500:
          err.message = 'Internal Server Error (500)'
          break
        case 501:
          err.message = 'Not Implemented (501)'
          break
        case 502:
          err.message = 'Bad Gateway (502)'
          break
        case 503:
          err.message = 'Service Unavailable (503)'
          break
        case 504:
          err.message = 'Gateway Timeout (504)'
          break
        case 505:
          err.message = 'HTTP Version Not Supported (505)'
          break
        default:
          err.message = `Error Connection (${err.response.status})!`
      }
    } else {
      err.message = 'Network Error！'
    }
    return Promise.reject(err)
    // return Promise.resolve(err)
  },
)