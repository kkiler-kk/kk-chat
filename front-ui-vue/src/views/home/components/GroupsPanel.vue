<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { TeamOutlined } from '@ant-design/icons-vue'
import { useContactsStore } from '@/stores/contacts'
import { groupConvId } from '@/utils/conversation'
import CreateGroupModal from '@/components/CreateGroupModal.vue'
import type { GroupItem } from '@/types/api'

const emit = defineEmits<{ 'go-chat': [convId: string] }>()

const contacts = useContactsStore()
const createOpen = ref(false)
const searchKw = ref('')
const searching = ref(false)

function openChat(g: GroupItem) {
  emit('go-chat', groupConvId(g.id))
}

async function onSearch() {
  const kw = searchKw.value.trim()
  searching.value = !!kw
  if (kw) await contacts.searchGroups(kw)
}

async function join(g: GroupItem) {
  await contacts.joinGroup(g.id)
  message.success(`已加入「${g.name}」`)
  openChat(g)
}

function onCreated(g: GroupItem) {
  openChat(g)
}
</script>

<template>
  <div class="groups-panel">
    <div class="ops">
      <a-button type="primary" ghost size="small" @click="createOpen = true">
        <template #icon><PlusOutlined /></template>
        发起群聊
      </a-button>
    </div>
    <div class="search">
      <a-input-search
        v-model:value="searchKw"
        placeholder="搜索群名称"
        size="small"
        allow-clear
        @search="onSearch"
      />
    </div>

    <template v-if="searching">
      <div class="section-title">搜索结果</div>
      <a-empty v-if="contacts.groupSearch.length === 0" description="未找到群聊" />
      <div v-for="g in contacts.groupSearch" :key="g.id" class="group-row">
        <a-avatar :size="36" shape="square" :src="g.avatar || undefined">
          <template v-if="!g.avatar"><TeamOutlined /></template>
        </a-avatar>
        <div class="gmeta">
          <div class="gname">{{ g.name }}</div>
          <div class="gcount">{{ g.member_count }} 名成员</div>
        </div>
        <a-button size="small" type="primary" ghost @click="join(g)">加入</a-button>
      </div>
    </template>

    <template v-else>
      <a-empty v-if="contacts.groups.length === 0" description="暂无群聊" class="empty" />
      <div v-for="g in contacts.groups" :key="g.id" class="group-row" @click="openChat(g)">
        <a-avatar :size="36" shape="square" :src="g.avatar || undefined">
          <template v-if="!g.avatar"><TeamOutlined /></template>
        </a-avatar>
        <div class="gmeta">
          <div class="gname">{{ g.name }}</div>
          <div class="gcount">{{ g.member_count }} 名成员</div>
        </div>
      </div>
    </template>

    <CreateGroupModal v-model:open="createOpen" @created="onCreated" />
  </div>
</template>

<style scoped>
.groups-panel {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.ops {
  padding: 4px 12px 8px;
}

.search {
  padding: 0 12px 8px;
}

.section-title {
  padding: 4px 12px;
  font-size: 12px;
  color: #999;
}

.empty {
  margin-top: 48px;
}

.group-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  cursor: pointer;
}

.group-row:hover {
  background: #f5f5f5;
}

.gmeta {
  flex: 1;
  min-width: 0;
}

.gname {
  font-size: 14px;
  color: #333;
}

.gcount {
  font-size: 12px;
  color: #999;
}
</style>
