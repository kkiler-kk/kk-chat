<script setup lang="ts">
import { ref } from 'vue'
import { useContactsStore } from '@/stores/contacts'
import UserCard from '@/components/UserCard.vue'
import type { UserDetail } from '@/types/api'

const open = defineModel<boolean>('open', { required: true })
const emit = defineEmits<{ 'go-chat': [user: UserDetail] }>()

const contacts = useContactsStore()
const keyword = ref('')
const cardUserId = ref<number | null>(null)

let debounceTimer: ReturnType<typeof setTimeout> | null = null

function onSearchInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    void contacts.searchUsers(keyword.value)
  }, 300)
}

function onGoChat(user: UserDetail) {
  open.value = false
  emit('go-chat', user)
}
</script>

<template>
  <a-modal v-model:open="open" title="搜索用户" :footer="null" width="480px">
    <a-input-search
      v-model:value="keyword"
      placeholder="按账号 / 昵称 / 邮箱搜索"
      allow-clear
      @input="onSearchInput"
      @search="() => contacts.searchUsers(keyword)"
    />
    <div class="results">
      <div
        v-for="item in contacts.userSearch"
        :key="item.id"
        class="result-row"
        @click="cardUserId = item.id"
      >
        <a-avatar :size="36" shape="square" :src="item.avatar || undefined">
          {{ item.name?.charAt(0) || '?' }}
        </a-avatar>
        <div class="meta">
          <div class="name">{{ item.name }}</div>
          <div class="identity">{{ item.identity }}</div>
        </div>
        <a-tag v-if="item.is_friend" color="green">好友</a-tag>
      </div>
      <a-empty
        v-if="keyword && contacts.userSearch.length === 0"
        description="未找到用户"
        :image-style="{ height: '48px' }"
      />
    </div>

    <a-modal :open="cardUserId !== null" :footer="null" :width="360" @cancel="cardUserId = null">
      <UserCard v-if="cardUserId !== null" :user-id="cardUserId" @go-chat="onGoChat" />
    </a-modal>
  </a-modal>
</template>

<style scoped>
.results {
  margin-top: 12px;
  max-height: 360px;
  overflow-y: auto;
}

.result-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 4px;
  border-radius: 6px;
  cursor: pointer;
}

.result-row:hover {
  background: #f5f5f5;
}

.meta {
  flex: 1;
  min-width: 0;
}

.name {
  font-size: 14px;
  color: #333;
}

.identity {
  font-size: 12px;
  color: #999;
}
</style>
