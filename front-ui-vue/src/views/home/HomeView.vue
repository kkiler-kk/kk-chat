<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import SideNav from './components/SideNav.vue'
import TheHeader from './components/TheHeader.vue'
import ConversationList from './components/ConversationList.vue'
import ChatWindow from './components/ChatWindow.vue'
import ContactsPanel from './components/ContactsPanel.vue'
import GroupsPanel from './components/GroupsPanel.vue'
import SettingsPanel from './components/SettingsPanel.vue'
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

/** 联系人/群组面板点击聊天 → 打开会话并切回消息面板 */
async function goChat(convId: string) {
  await chat.openConversation(convId)
  activePanel.value = 'chats'
}
</script>

<template>
  <div class="home">
    <SideNav v-model:active-panel="activePanel" />
    <div class="middle">
      <TheHeader :title="panelTitle" />
      <ConversationList v-if="activePanel === 'chats'" />
      <ContactsPanel v-else-if="activePanel === 'contacts'" @go-chat="goChat" />
      <GroupsPanel v-else-if="activePanel === 'groups'" @go-chat="goChat" />
      <SettingsPanel v-else-if="activePanel === 'settings'" />
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
</style>
