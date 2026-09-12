<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { PictureOutlined, SendOutlined } from '@ant-design/icons-vue'
import { useChatStore } from '@/stores/chat'
import { useAuthStore } from '@/stores/auth'
import { usePresenceStore } from '@/stores/presence'
import { ApiError } from '@/api/client'
import { uploadImage } from '@/api/file'
import { ErrCode } from '@/types/errorcode'
import MessageBubble from '@/components/MessageBubble.vue'
import EmojiPicker from '@/components/EmojiPicker.vue'

const chat = useChatStore()
const auth = useAuthStore()
const presence = usePresenceStore()

const draft = ref('')
const sending = ref(false)
const limitAlert = ref(false)
const listRef = ref<HTMLElement | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

const conv = computed(() => chat.activeConversation)
const online = computed(() => (conv.value ? presence.isOnline(conv.value.peer_id) : false))

function scrollToBottom() {
  void nextTick(() => {
    if (listRef.value) listRef.value.scrollTop = listRef.value.scrollHeight
  })
}

watch(
  () => chat.activeMessages.length,
  (len, oldLen) => {
    // 新消息到底部；加载历史（头部插入）保持位置不动
    if (len > (oldLen ?? 0)) scrollToBottom()
  },
)
watch(
  () => chat.activeConvId,
  () => {
    limitAlert.value = false
    scrollToBottom()
  },
)

async function onScroll() {
  const el = listRef.value
  if (!el || !chat.activeConvId) return
  if (el.scrollTop < 40) {
    const prevHeight = el.scrollHeight
    await chat.loadMoreHistory(chat.activeConvId)
    void nextTick(() => {
      el.scrollTop = el.scrollHeight - prevHeight
    })
  }
}

async function send() {
  const text = draft.value.trim()
  if (!text || sending.value) return
  sending.value = true
  try {
    await chat.send(text)
    draft.value = ''
    limitAlert.value = false
    scrollToBottom()
  } catch (e) {
    if (e instanceof ApiError && e.code === ErrCode.NotFriendLimit) limitAlert.value = true
  } finally {
    sending.value = false
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    void send()
  }
}

function insertEmoji(emoji: string) {
  draft.value += emoji
}

async function onPickImage(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = '' // 允许重复选同一文件
  if (!file) return
  if (file.size > 5 * 1024 * 1024) {
    void import('ant-design-vue').then(({ message }) => message.error('图片不能超过 5MB'))
    return
  }
  const res = await uploadImage(file)
  await chat.send(res.url, 'image')
  scrollToBottom()
}
</script>

<template>
  <div class="chat-window">
    <template v-if="conv">
      <div class="chat-header">
        <span class="peer">{{ conv.peer_name }}</span>
        <span v-if="conv.type === 'private'" class="status" :class="{ on: online }">
          {{ online ? '在线' : '离线' }}
        </span>
        <span v-else class="status">群聊</span>
      </div>

      <div ref="listRef" class="msg-list" @scroll="onScroll">
        <MessageBubble
          v-for="m in chat.activeMessages"
          :key="m.id"
          :message="m"
          :self="m.sender_id === auth.user?.id"
        />
      </div>

      <a-alert
        v-if="limitAlert"
        type="warning"
        show-icon
        banner
        message="非好友关系发送消息超过 3 条，请先添加对方为好友"
      />

      <div class="input-area">
        <div class="toolbar">
          <EmojiPicker @select="insertEmoji" />
          <a-button type="text" title="发送图片" @click="fileInput?.click()">
            <template #icon><PictureOutlined /></template>
          </a-button>
          <input
            ref="fileInput"
            type="file"
            accept="image/jpeg,image/png,image/gif,image/webp"
            hidden
            @change="onPickImage"
          />
        </div>
        <a-textarea
          v-model:value="draft"
          :rows="3"
          placeholder="输入消息，Enter 发送 / Shift+Enter 换行"
          :maxlength="4096"
          @keydown="onKeydown"
        />
        <div class="send-row">
          <a-button type="primary" :loading="sending" :disabled="!draft.trim()" @click="send">
            <template #icon><SendOutlined /></template>
            发送
          </a-button>
        </div>
      </div>
    </template>
    <div v-else class="placeholder">
      <a-empty description="选择左侧会话开始聊天" />
    </div>
  </div>
</template>

<style scoped>
.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  background: #f5f5f5;
}

.chat-header {
  height: 56px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.peer {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.status {
  font-size: 12px;
  color: #bbb;
}

.status.on {
  color: #52c41a;
}

.msg-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.input-area {
  background: #fff;
  border-top: 1px solid #f0f0f0;
  padding: 8px 12px 12px;
  flex-shrink: 0;
}

.toolbar {
  display: flex;
  gap: 4px;
  margin-bottom: 4px;
}

.send-row {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}

.placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
