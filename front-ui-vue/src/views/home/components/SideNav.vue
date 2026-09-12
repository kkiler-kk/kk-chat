<script setup lang="ts">
import {
  MessageOutlined,
  TeamOutlined,
  UsergroupAddOutlined,
  SettingOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import type { PanelKey } from '../types'

const props = defineProps<{ activePanel: PanelKey }>()
const emit = defineEmits<{ 'update:activePanel': [key: PanelKey] }>()
const auth = useAuthStore()

const items: Array<{ key: PanelKey; icon: typeof MessageOutlined; title: string }> = [
  { key: 'chats', icon: MessageOutlined, title: '消息' },
  { key: 'contacts', icon: TeamOutlined, title: '联系人' },
  { key: 'groups', icon: UsergroupAddOutlined, title: '群组' },
  { key: 'settings', icon: SettingOutlined, title: '设置' },
]
</script>

<template>
  <div class="side-nav">
    <a-avatar
      class="me"
      :size="40"
      :src="auth.user?.avatar || undefined"
      title="我的资料"
      @click="emit('update:activePanel', 'settings')"
    >
      {{ auth.user?.name?.charAt(0) || '?' }}
    </a-avatar>
    <div class="nav-items">
      <div
        v-for="item in items"
        :key="item.key"
        class="nav-item"
        :class="{ active: props.activePanel === item.key }"
        :title="item.title"
        @click="emit('update:activePanel', item.key)"
      >
        <component :is="item.icon" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.side-nav {
  width: 64px;
  height: 100%;
  background: #2e2e32;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 0;
  gap: 20px;
  flex-shrink: 0;
}

.me {
  cursor: pointer;
}

.nav-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.nav-item {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: #97979a;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-item:hover {
  color: #fff;
}

.nav-item.active {
  color: #fff;
  background: #44444a;
}
</style>
