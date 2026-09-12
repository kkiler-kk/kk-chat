<script setup lang="ts">
import { computed, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useContactsStore } from '@/stores/contacts'
import type { GroupItem } from '@/types/api'

const open = defineModel<boolean>('open', { required: true })
const emit = defineEmits<{ created: [group: GroupItem] }>()

const contacts = useContactsStore()
const name = ref('')
const checkedIds = ref<number[]>([])
const submitting = ref(false)

const options = computed(() => contacts.friends.map((f) => ({ label: f.name, value: f.id })))

async function submit() {
  if (!name.value.trim()) {
    message.warning('请输入群名称')
    return
  }
  submitting.value = true
  try {
    const g = await contacts.createGroup(name.value.trim(), checkedIds.value)
    message.success('群聊已创建')
    open.value = false
    name.value = ''
    checkedIds.value = []
    emit('created', g)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <a-modal
    v-model:open="open"
    title="发起群聊"
    :confirm-loading="submitting"
    ok-text="创建"
    cancel-text="取消"
    @ok="submit"
  >
    <a-form layout="vertical">
      <a-form-item label="群名称" required>
        <a-input v-model:value="name" placeholder="群名称" :maxlength="64" />
      </a-form-item>
      <a-form-item label="邀请好友">
        <a-checkbox-group v-model:value="checkedIds" :options="options" class="members" />
        <a-empty
          v-if="contacts.friends.length === 0"
          description="还没有好友，先去添加吧"
          :image-style="{ height: '48px' }"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<style scoped>
.members {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 240px;
  overflow-y: auto;
}
</style>
