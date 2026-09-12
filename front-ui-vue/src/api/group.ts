import { request } from './client'
import type { GroupItem } from '@/types/api'

export interface CreateGroupParams {
  name: string
  member_ids?: number[]
}

export const create = (p: CreateGroupParams) =>
  request<GroupItem>({ url: '/api/v1/groups', method: 'POST', data: p })

export const join = (id: number) =>
  request<null>({ url: `/api/v1/groups/${id}/members`, method: 'POST' })

export const invite = (id: number, userIds: number[]) =>
  request<null>({ url: `/api/v1/groups/${id}/invite`, method: 'POST', data: { user_ids: userIds } })

export const list = () => request<GroupItem[]>({ url: '/api/v1/groups', method: 'GET' })

export const search = (keyword: string) =>
  request<GroupItem[]>({ url: '/api/v1/groups', method: 'GET', params: { search: keyword } })
