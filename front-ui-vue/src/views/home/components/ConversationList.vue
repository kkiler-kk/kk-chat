<script setup lang="ts">
import { computed, ref } from 'vue'
import { TeamOutlined } from '@ant-design/icons-vue'
import { useChatStore } from '@/stores/chat'
import { usePresenceStore } from '@/stores/presence'
import { formatPast } from '@/utils/format'
import type { RecentConversation } from '@/types/api'

const chat = useChatStore()
const presence = usePresenceStore()
const keyword = ref('')

const filtered = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return chat.conversations
  return chat.conversations.filter((c) => c.peer_name.toLowerCase().includes(kw))
})

function open(c: RecentConversation) {
  void chat.openConversation(c.conversation_id)
}

function isOnline(c: RecentConversation): boolean {
  return c.type === 'private' && presence.isOnline(c.peer_id)
}

function lastText(c: RecentConversation): string {
  const prefix = c.last_sender_name ? `${c.last_sender_name}: ` : ''
  return prefix + c.last_message_content
}
</script>

<template>
  <div class="conv-list">
    <div class="search">
      <a-input v-model:value="keyword" placeholder="搜索会话" allow-clear />
    </div>
    <a-empty v-if="filtered.length === 0" description="暂无会话" class="empty" />
    <div
      v-for="c in filtered"
      :key="c.conversation_id"
      class="conv-item"
      :class="{ active: chat.activeConvId === c.conversation_id }"
      @click="open(c)"
    >
      <a-badge :dot="isOnline(c)" :offset="[-6, 36]" color="#52c41a">
        <a-avatar :size="42" shape="square" :src="c.peer_avatar || undefined">
          <template v-if="c.type === 'group' && !c.peer_avatar"><TeamOutlined /></template>
          <template v-else-if="!c.peer_avatar">{{ c.peer_name?.charAt(0) || '?' }}</template>
        </a-avatar>
      </a-badge>
      <div class="conv-main">
        <div class="conv-top">
          <span class="conv-name">{{ c.peer_name }}</span>
          <span class="conv-time">{{ formatPast(c.last_time) }}</span>
        </div>
        <div class="conv-last">{{ lastText(c) }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.conv-list {
  flex: 1;
  overflow-y: auto;
}

.search {
  padding: 10px 12px;
}

.empty {
  margin-top: 60px;
}

.conv-item {
  display: flex;
  gap: 10px;
  padding: 10px 12px;
  cursor: pointer;
  transition: background 0.15s;
}

.conv-item:hover {
  background: #f5f5f5;
}

.conv-item.active {
  background: #e6f4ff;
}

.conv-main {
  flex: 1;
  min-width: 0;
}

.conv-top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.conv-name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.conv-time {
  font-size: 12px;
  color: #bbb;
  flex-shrink: 0;
}

.conv-last {
  margin-top: 4px;
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
