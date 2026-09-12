import { wsClient } from './WSClient'
import { useChatStore } from '@/stores/chat'
import { useContactsStore } from '@/stores/contacts'

/**
 * WS 事件 → store action 的唯一接线处。
 * 在登录成功/页面恢复后调用一次；返回统一解绑函数（登出时调用）。
 */
export function setupWsBindings(): () => void {
  const chat = useChatStore()
  const contacts = useContactsStore()
  const offs = [
    wsClient.on('chat.message', (msg) => chat.handleIncoming(msg)),
    wsClient.on('chat.recent_updated', (conv) => chat.handleRecentUpdated(conv)),
    wsClient.on('presence.changed', (d) => contacts.handlePresenceChanged(d.user_id, d.online)),
    wsClient.on('group.updated', () => {
      void contacts.loadGroups()
    }),
  ]
  return () => offs.forEach((off) => off())
}
