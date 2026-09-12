import { defineStore } from 'pinia'
import * as chatApi from '@/api/chat'
import type { OutgoingMessage, RecentConversation } from '@/types/api'

export const useChatStore = defineStore('chat', {
  state: () => ({
    conversations: [] as RecentConversation[],
    messagesByConv: {} as Record<string, OutgoingMessage[]>,
    activeConvId: null as string | null,
    historyCursor: {} as Record<string, string>, // convId → 最旧一条的 created_at
    historyDone: {} as Record<string, boolean>,
    loading: false,
  }),
  getters: {
    activeMessages: (s): OutgoingMessage[] =>
      s.activeConvId ? (s.messagesByConv[s.activeConvId] ?? []) : [],
    activeConversation: (s): RecentConversation | undefined =>
      s.conversations.find((c) => c.conversation_id === s.activeConvId),
  },
  actions: {
    async loadConversations() {
      this.conversations = await chatApi.conversations()
    },
    async openConversation(convId: string) {
      this.activeConvId = convId
      if (this.messagesByConv[convId]) return // 已加载过
      this.loading = true
      try {
        const msgs = await chatApi.history(convId, undefined, 50)
        this.messagesByConv[convId] = msgs
        if (msgs.length > 0) this.historyCursor[convId] = msgs[0].created_at
        this.historyDone[convId] = msgs.length < 50
      } finally {
        this.loading = false
      }
    },
    async loadMoreHistory(convId: string) {
      const cursor = this.historyCursor[convId]
      if (!cursor || this.historyDone[convId]) return
      const msgs = await chatApi.history(convId, cursor, 50)
      this.messagesByConv[convId] = [...msgs, ...(this.messagesByConv[convId] ?? [])]
      if (msgs.length > 0) this.historyCursor[convId] = msgs[0].created_at
      this.historyDone[convId] = msgs.length < 50
    },
    async send(content: string, contentType: 'text' | 'image' = 'text') {
      if (!this.activeConvId) return
      const msg = await chatApi.sendMessage({
        conversation_id: this.activeConvId,
        content,
        content_type: contentType,
      })
      this.handleIncoming(msg) // 自己发的也走统一入口（后端不回推给自己）
    },
    /** WS chat.message → 唯一入口：消息归属判断只在这里做一次 */
    handleIncoming(msg: OutgoingMessage) {
      const list = this.messagesByConv[msg.conversation_id]
      if (list && !list.some((m) => m.id === msg.id)) list.push(msg)
      else if (!list) this.messagesByConv[msg.conversation_id] = [msg]
    },
    /** WS chat.recent_updated → 更新/置顶会话 */
    handleRecentUpdated(conv: RecentConversation) {
      const idx = this.conversations.findIndex((c) => c.conversation_id === conv.conversation_id)
      if (idx >= 0) this.conversations.splice(idx, 1, conv)
      else this.conversations.unshift(conv)
      this.conversations.sort(
        (a, b) => new Date(b.last_time).getTime() - new Date(a.last_time).getTime(),
      )
    },
    /** 退出登录时清空 */
    $resetAll() {
      this.conversations = []
      this.messagesByConv = {}
      this.activeConvId = null
      this.historyCursor = {}
      this.historyDone = {}
      this.loading = false
    },
  },
})
