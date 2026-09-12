<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import SideNav from './components/SideNav.vue'
import TheHeader from './components/TheHeader.vue'
import ConversationList from './components/ConversationList.vue'
import ChatWindow from './components/ChatWindow.vue'
import type { PanelKey } from './types'
import { useChatStore } from '@/stores/chat'
import { useContactsStore } from '@/stores/contacts'

const activePanel = ref<PanelKey>('chats')
const chat = useChatStore()
const contacts = useContactsStore()

const panelTitle = computed(() => {
  switch (activePanel.value) {
    case 'chats':
      return '消息'
    case 'contacts':
      return '联系人'
    case 'groups':
      return '群组'
    case 'settings':
      return '设置'
    default:
      return ''
  }
})

onMounted(() => {
  void Promise.all([chat.loadConversations(), contacts.loadFriends(), contacts.loadGroups()])
})
</script>

<template>
  <div class="home">
    <SideNav v-model:active-panel="activePanel" />
    <div class="middle">
      <TheHeader :title="panelTitle" />
      <ConversationList v-if="activePanel === 'chats'" />
      <div v-else class="panel-stub">「{{ panelTitle }}」面板见 Task F9</div>
    </div>
    <ChatWindow />
  </div>
</template>

<style scoped>
.home {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.middle {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-right: 1px solid #f0f0f0;
}

.panel-stub {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
}
</style>
