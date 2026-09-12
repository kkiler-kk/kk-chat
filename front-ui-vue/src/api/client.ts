import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { message } from 'ant-design-vue'
import { Local, StorageKeys } from '@/utils/storage'
import type { ApiResponse } from '@/types/api'
import { ErrCode } from '@/types/errorcode'

export class ApiError extends Error {
  code: number
  requestId: string
  status: number

  constructor(code: number, msg: string, requestId = '', status = 0) {
    super(msg)
    this.name = 'ApiError'
    this.code = code
    this.requestId = requestId
    this.status = status
  }
}

let unauthorizedHandler: () => void = () => {}

/** 注册 401/token 失效的全局处理器（auth store 调用） */
export function setUnauthorizedHandler(fn: () => void): void {
  unauthorizedHandler = fn
}

const http = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use((config) => {
  const token = Local.get<string>(StorageKeys.token)
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

http.interceptors.response.use(
  (resp: AxiosResponse<ApiResponse<unknown>>) => {
    const body = resp.data
    if (body.code === ErrCode.OK) return resp
    // 业务错误：统一 toast + 抛出 ApiError（token 类错误不 toast，交给 unauthorized 流程）
    const err = new ApiError(body.code, body.message, body.request_id, resp.status)
    if (err.code !== ErrCode.TokenInvalid) message.error(body.message)
    return Promise.reject(err)
  },
  (error: unknown) => {
    if (axios.isAxiosError(error) && error.response) {
      const body = error.response.data as ApiResponse<unknown> | undefined
      const err = new ApiError(
        body?.code ?? error.response.status * 100,
        body?.message ?? error.message,
        body?.request_id ?? '',
        error.response.status,
      )
      if (err.status === 401 || err.code === ErrCode.TokenInvalid) {
        unauthorizedHandler()
      } else {
        message.error(err.message)
      }
      return Promise.reject(err)
    }
    // 网络层错误（超时/断网）
    const msg =
      axios.isAxiosError(error) && error.code === 'ECONNABORTED' ? '网络超时' : '网络连接错误'
    message.error(msg)
    return Promise.reject(new ApiError(ErrCode.Internal, msg, '', 0))
  },
)

/** request 泛型封装：直接 resolve 信封中的 data 字段 */
export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  const resp = await http.request<ApiResponse<T>>(config)
  return resp.data.data
}
