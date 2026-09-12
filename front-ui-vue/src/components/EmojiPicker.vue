<script setup lang="ts">
import Vue3EmojiPicker, { type EmojiExt } from 'vue3-emoji-picker'
import 'vue3-emoji-picker/css'
import { ref } from 'vue'
import { SmileOutlined } from '@ant-design/icons-vue'

const emit = defineEmits<{ select: [emoji: string] }>()
const open = ref(false)
const selected = ref('')

function onSelect(emoji: EmojiExt) {
  emit('select', emoji.i)
  open.value = false
}
</script>

<template>
  <a-popover v-model:open="open" trigger="click" placement="topRight">
    <template #content>
      <Vue3EmojiPicker v-model="selected" :hide-search="false" @select="onSelect" />
    </template>
    <a-button type="text" title="表情">
      <template #icon><SmileOutlined /></template>
    </a-button>
  </a-popover>
</template>
