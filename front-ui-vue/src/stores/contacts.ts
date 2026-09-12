import { defineStore } from 'pinia'
import * as friendApi from '@/api/friend'
import * as groupApi from '@/api/group'
import * as userApi from '@/api/user'
import { usePresenceStore } from './presence'
import type { FriendItem, GroupItem, UserSearchItem } from '@/types/api'

export const useContactsStore = defineStore('contacts', {
  state: () => ({
    friends: [] as FriendItem[],
    groups: [] as GroupItem[],
    userSearch: [] as UserSearchItem[],
    groupSearch: [] as GroupItem[],
  }),
  actions: {
    async loadFriends() {
      this.friends = await friendApi.list()
      usePresenceStore().applyFriends(this.friends)
    },
    async loadGroups() {
      this.groups = await groupApi.list()
    },
    async addFriend(uid: number) {
      await friendApi.add(uid)
      await this.loadFriends()
    },
    async createGroup(name: string, memberIds: number[]) {
      const g = await groupApi.create({ name, member_ids: memberIds })
      await this.loadGroups()
      return g
    },
    async joinGroup(gid: number) {
      await groupApi.join(gid)
      await this.loadGroups()
    },
    async searchUsers(kw: string) {
      this.userSearch = kw.trim() ? await userApi.search(kw) : []
    },
    async searchGroups(kw: string) {
      this.groupSearch = kw.trim() ? await groupApi.search(kw) : []
    },
    /** WS presence.changed → 同步在线状态 */
    handlePresenceChanged(uid: number, online: boolean) {
      usePresenceStore().set(uid, online)
      const f = this.friends.find((item) => item.id === uid)
      if (f) f.online = online
    },
    $reset() {
      this.friends = []
      this.groups = []
      this.userSearch = []
      this.groupSearch = []
      usePresenceStore().$reset()
    },
  },
})
