<script setup lang="ts">
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import * as userApi from '@/api/user'
import { useContactsStore } from '@/stores/contacts'
import { formatTime } from '@/utils/format'
import type { UserDetail } from '@/types/api'

const props = defineProps<{ userId: number }>()
const emit = defineEmits<{ 'go-chat': [user: UserDetail] }>()

const contacts = useContactsStore()
const detail = ref<UserDetail | null>(null)
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    detail.value = await userApi.getDetail(props.userId)
  } finally {
    loading.value = false
  }
}

watch(() => props.userId, load, { immediate: true })

async function addFriend() {
  if (!detail.value) return
  try {
    await contacts.addFriend(detail.value.id)
    message.success('已添加为好友')
    await load()
  } catch {
    /* 拦截器已 toast */
  }
}
</script>

<template>
  <a-spin :spinning="loading">
    <div v-if="detail" class="user-card">
      <a-avatar :size="64" shape="square" :src="detail.avatar || undefined">
        {{ detail.name?.charAt(0) || '?' }}
      </a-avatar>
      <div class="info">
        <div class="name">
          {{ detail.name }}
          <a-tag v-if="detail.is_self" color="blue">本人</a-tag>
          <a-tag v-else-if="detail.is_friend" color="green">好友</a-tag>
        </div>
        <div class="row">账号：{{ detail.identity }}</div>
        <div v-if="detail.signature" class="row">签名：{{ detail.signature }}</div>
        <div v-if="detail.is_friend || detail.is_self" class="row">邮箱：{{ detail.email }}</div>
        <div v-if="(detail.is_friend || detail.is_self) && detail.phone" class="row">
          手机：{{ detail.phone }}
        </div>
        <div class="row">注册于 {{ formatTime(detail.created_at, 'YYYY-MM-DD') }}</div>
      </div>
      <div class="ops">
        <a-button v-if="!detail.is_friend && !detail.is_self" type="primary" @click="addFriend">
          添加好友
        </a-button>
        <a-button v-else-if="!detail.is_self" type="primary" @click="emit('go-chat', detail)">
          发消息
        </a-button>
      </div>
    </div>
  </a-spin>
</template>

<style scoped>
.user-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 260px;
}

.info .name {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 6px;
}

.info .row {
  font-size: 13px;
  color: #666;
  margin-top: 2px;
}

.ops {
  display: flex;
  justify-content: flex-end;
}
</style>
