import axios from 'axios'
export const API_URL = ''
export const cluster = 'local'

export async function Get(url: string, options?: any) {
  return await axios.get(url, {
    headers: { 'Content-Type': 'application/json', 'x-requested-with': 'ajax' },
    withCredentials: true,
    timeout: 60000,
    params: options,
  }).then(
    (response) => {
      return response
    },
  )
}

export async function Put(url: string, param?: any) {
  return await axios.put(url, param, {
    headers: { 'Content-Type': 'application/json', 'x-requested-with': 'ajax', timeout: 60000 },
    withCredentials: true,
  })
}

export async function Post(url: string, param?: any) {
  return await axios.post(url, param, {
    headers: { 'Content-Type': 'application/json', 'x-requested-with': 'ajax', timeout: 60000 },
    withCredentials: true,
  })
}

export async function Delete(url: string) {
  return await axios.delete(url, {
    headers: { 'Content-Type': 'application/json', 'x-requested-with': 'ajax', timeout: 60000 },
    withCredentials: true,
  })
}

axios.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response.status === 401) {
      let callbackUrl = window.location.hash
      window.location.hash = `/login?from=${callbackUrl}`
    }
    return error
  })
