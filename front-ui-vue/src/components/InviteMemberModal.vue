<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useContactsStore } from '@/stores/contacts'
import * as groupApi from '@/api/group'

const props = defineProps<{ groupId: number | null }>()
const open = defineModel<boolean>('open', { required: true })

const contacts = useContactsStore()
const checkedIds = ref<number[]>([])
const submitting = ref(false)

const options = computed(() => contacts.friends.map((f) => ({ label: f.name, value: f.id })))

watch(open, (v) => {
  if (v) checkedIds.value = []
})

async function submit() {
  if (props.groupId === null || checkedIds.value.length === 0) {
    message.warning('请选择要邀请的好友')
    return
  }
  submitting.value = true
  try {
    await groupApi.invite(props.groupId, checkedIds.value)
    message.success('已邀请入群')
    open.value = false
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <a-modal
    v-model:open="open"
    title="邀请好友入群"
    :confirm-loading="submitting"
    ok-text="邀请"
    cancel-text="取消"
    @ok="submit"
  >
    <a-checkbox-group v-model:value="checkedIds" :options="options" class="members" />
    <a-empty
      v-if="contacts.friends.length === 0"
      description="还没有好友，先去添加吧"
      :image-style="{ height: '48px' }"
    />
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
