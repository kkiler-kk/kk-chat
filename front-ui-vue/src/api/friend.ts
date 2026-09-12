import { request } from './client'
import type { FriendItem } from '@/types/api'

export const add = (userId: number) =>
  request<null>({ url: '/api/v1/friends', method: 'POST', data: { user_id: userId } })

export const list = () => request<FriendItem[]>({ url: '/api/v1/friends', method: 'GET' })
