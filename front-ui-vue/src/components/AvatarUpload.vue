<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { uploadImage } from '@/api/file'

const props = defineProps<{ src?: string }>()
const emit = defineEmits<{ uploaded: [url: string] }>()

const uploading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

async function onChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (file.size > 5 * 1024 * 1024) {
    message.error('图片不能超过 5MB')
    return
  }
  uploading.value = true
  try {
    const res = await uploadImage(file)
    emit('uploaded', res.url)
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <a-spin :spinning="uploading">
    <div class="avatar-upload" title="点击更换头像" @click="fileInput?.click()">
      <a-avatar :size="72" shape="square" :src="props.src || undefined">
        {{ props.src ? '' : '点击上传' }}
      </a-avatar>
      <input
        ref="fileInput"
        type="file"
        accept="image/jpeg,image/png,image/gif,image/webp"
        hidden
        @change="onChange"
      />
    </div>
  </a-spin>
</template>

<style scoped>
.avatar-upload {
  cursor: pointer;
  display: inline-block;
}
</style>
