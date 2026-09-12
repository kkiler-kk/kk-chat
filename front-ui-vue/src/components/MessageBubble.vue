<script setup lang="ts">
import { computed } from 'vue'
import { formatPast } from '@/utils/format'
import type { OutgoingMessage } from '@/types/api'

const props = defineProps<{
  message: OutgoingMessage
  self: boolean
}>()

const isImage = computed(() => props.message.content_type === 'image')

function openImage() {
  if (isImage.value) window.open(props.message.content, '_blank')
}
</script>

<template>
  <div class="bubble-row" :class="{ self }">
    <a-avatar shape="square" :size="40" :src="message.sender_avatar || undefined">
      {{ message.sender_name?.charAt(0) || '?' }}
    </a-avatar>
    <div class="bubble-body">
      <div class="meta">
        <span class="name">{{ message.sender_name }}</span>
        <span class="time">{{ formatPast(message.created_at) }}</span>
      </div>
      <div class="bubble" :class="{ 'bubble-self': self }">
        <img
          v-if="isImage"
          :src="message.content"
          class="msg-image"
          alt="图片消息"
          @click="openImage"
        />
        <template v-else>{{ message.content }}</template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.bubble-row {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  padding: 0 16px;
}

.bubble-row.self {
  flex-direction: row-reverse;
}

.bubble-body {
  max-width: 60%;
}

.meta {
  display: flex;
  gap: 8px;
  align-items: baseline;
  margin-bottom: 4px;
  font-size: 12px;
  color: #999;
}

.self .meta {
  flex-direction: row-reverse;
}

.name {
  color: #666;
}

.bubble {
  position: relative;
  padding: 8px 12px;
  border-radius: 4px 12px 12px 12px;
  background: #fff;
  color: rgb(70, 74, 77);
  font-size: 14px;
  line-height: 22px;
  word-break: break-word;
  white-space: pre-wrap;
  box-shadow: 0 1px 2px rgb(0 0 0 / 8%);
}

.bubble-self {
  border-radius: 12px 4px 12px 12px;
  background: #95ec69;
}

.msg-image {
  max-width: 240px;
  max-height: 320px;
  border-radius: 8px;
  cursor: pointer;
  display: block;
}
</style>
