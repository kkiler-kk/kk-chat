<script setup lang="ts">
import { ref } from 'vue'
import { UserAddOutlined } from '@ant-design/icons-vue'
import { useContactsStore } from '@/stores/contacts'
import { usePresenceStore } from '@/stores/presence'
import { privateConvId } from '@/utils/conversation'
import { useAuthStore } from '@/stores/auth'
import SearchUserModal from '@/components/SearchUserModal.vue'
import UserCard from '@/components/UserCard.vue'
import type { FriendItem, UserDetail } from '@/types/api'

const emit = defineEmits<{ 'go-chat': [convId: string] }>()

const contacts = useContactsStore()
const presence = usePresenceStore()
const auth = useAuthStore()

const searchOpen = ref(false)
const cardUserId = ref<number | null>(null)

function openChat(f: FriendItem) {
  if (!auth.user) return
  emit('go-chat', privateConvId(auth.user.id, f.id))
}

function onCardGoChat(user: UserDetail) {
  cardUserId.value = null
  if (!auth.user) return
  emit('go-chat', privateConvId(auth.user.id, user.id))
}
</script>

<template>
  <div class="contacts-panel">
    <div class="ops">
      <a-button type="primary" ghost size="small" @click="searchOpen = true">
        <template #icon><UserAddOutlined /></template>
        添加好友
      </a-button>
    </div>
    <a-empty v-if="contacts.friends.length === 0" description="暂无好友" class="empty" />
    <div v-for="f in contacts.friends" :key="f.id" class="friend-row" @click="openChat(f)">
      <a-badge :dot="presence.isOnline(f.id)" :offset="[-6, 30]" color="#52c41a">
        <a-avatar :size="36" shape="square" :src="f.avatar || undefined">
          {{ f.name?.charAt(0) || '?' }}
        </a-avatar>
      </a-badge>
      <span class="fname">{{ f.name }}</span>
      <a-button type="text" size="small" class="detail-btn" @click.stop="cardUserId = f.id">
        资料
      </a-button>
    </div>

    <SearchUserModal v-model:open="searchOpen" @go-chat="onCardGoChat" />
    <a-modal :open="cardUserId !== null" :footer="null" :width="360" @cancel="cardUserId = null">
      <UserCard v-if="cardUserId !== null" :user-id="cardUserId" @go-chat="onCardGoChat" />
    </a-modal>
  </div>
</template>

<style scoped>
.contacts-panel {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.ops {
  padding: 4px 12px 10px;
}

.empty {
  margin-top: 48px;
}

.friend-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  cursor: pointer;
}

.friend-row:hover {
  background: #f5f5f5;
}

.fname {
  flex: 1;
  font-size: 14px;
  color: #333;
}

.detail-btn {
  color: #999;
}
</style>
