import { request } from './client'
import type { UserDetail, UserInfo, UserSearchItem } from '@/types/api'

export interface UpdateMeParams {
  name?: string
  phone?: string
  email?: string
  email_code?: string
  avatar?: string
  signature?: string
  birth_date?: string // YYYY-MM-DD
}

export const getMe = () => request<UserInfo>({ url: '/api/v1/users/me', method: 'GET' })

export const updateMe = (p: UpdateMeParams) =>
  request<null>({ url: '/api/v1/users/me', method: 'PATCH', data: p })

export const getDetail = (id: number) =>
  request<UserDetail>({ url: `/api/v1/users/${id}`, method: 'GET' })

export const search = (keyword: string) =>
  request<UserSearchItem[]>({ url: '/api/v1/users', method: 'GET', params: { search: keyword } })
