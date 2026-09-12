import type { OutgoingMessage, RecentConversation } from './api'

export interface WsEnvelope<T = unknown> {
  event: string
  data: T
  ts: number
}

/** 服务端 → 客户端事件与 payload 的映射（spec §3.4） */
export interface WsServerEvents {
  pong: undefined
  'chat.message': OutgoingMessage
  'chat.recent_updated': RecentConversation
  'presence.changed': { user_id: number; online: boolean }
  'group.updated': { group_id: number }
  'system.kick': { reason: string }
}

export type WsServerEvent = keyof WsServerEvents
