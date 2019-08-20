import { POST } from './http'
import { PATH } from '../config/config'

export function login(username: string, password: string) {
  const url = PATH.login
  return POST<string>(url, { username, password }, { success: 'Login Success' })
}