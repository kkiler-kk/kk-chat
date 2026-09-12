import { request } from './client'
import type { OutgoingMessage, RecentConversation } from '@/types/api'

export interface SendMessageParams {
  conversation_id: string
  content: string
  content_type: 'text' | 'image'
}

export const sendMessage = (p: SendMessageParams) =>
  request<OutgoingMessage>({ url: '/api/v1/messages', method: 'POST', data: p })

export const conversations = () =>
  request<RecentConversation[]>({ url: '/api/v1/conversations', method: 'GET' })

export const history = (convId: string, cursor?: string, limit = 50) =>
  request<OutgoingMessage[]>({
    url: `/api/v1/conversations/${encodeURIComponent(convId)}/messages`,
    method: 'GET',
    params: { cursor, limit },
  })
