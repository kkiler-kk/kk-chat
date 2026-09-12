import { defineStore } from 'pinia'
import type { FriendItem } from '@/types/api'

export const usePresenceStore = defineStore('presence', {
  state: () => ({
    onlineIds: new Set<number>(),
  }),
  actions: {
    set(uid: number, online: boolean) {
      const next = new Set(this.onlineIds)
      if (online) next.add(uid)
      else next.delete(uid)
      this.onlineIds = next // 整体替换保证响应式
    },
    isOnline(uid: number): boolean {
      return this.onlineIds.has(uid)
    },
    applyFriends(items: FriendItem[]) {
      this.onlineIds = new Set(items.filter((f) => f.online).map((f) => f.id))
    },
    $reset() {
      this.onlineIds = new Set()
    },
  },
})
