<script setup lang="ts">
import { onUnmounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import * as authApi from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import type { FormInstance } from 'ant-design-vue'

const emit = defineEmits<{ registered: [] }>()

const auth = useAuthStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const countdown = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const form = reactive({
  identity: '',
  name: '',
  password: '',
  confirmPassword: '',
  email: '',
  email_code: '',
})

const rules = {
  identity: [
    { required: true, message: '请输入唯一账号' },
    { min: 3, max: 32, message: '账号长度 3-32 位' },
  ],
  name: [{ required: true, message: '请输入昵称' }],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, max: 64, message: '密码长度 6-64 位' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码' },
    {
      validator: (_rule: unknown, value: string) =>
        value === form.password
          ? Promise.resolve()
          : Promise.reject(new Error('两次输入的密码不一致')),
    },
  ],
  email: [
    { required: true, message: '请输入邮箱' },
    { type: 'email' as const, message: '邮箱格式不正确' },
  ],
  email_code: [
    { required: true, message: '请输入邮箱验证码' },
    { len: 6, message: '验证码为 6 位数字' },
  ],
}

async function sendCode() {
  if (!form.email) {
    message.warning('请先填写邮箱')
    return
  }
  await authApi.sendEmailCode(form.email)
  message.success('验证码已发送')
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0 && timer) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

async function onSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    await auth.register({
      identity: form.identity,
      name: form.name,
      password: form.password,
      email: form.email,
      email_code: form.email_code,
    })
    message.success('注册成功，请登录')
    emit('registered')
  } finally {
    loading.value = false
  }
}

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <a-form ref="formRef" :model="form" :rules="rules" layout="vertical" @finish="onSubmit">
    <a-form-item label="唯一账号" name="identity">
      <a-input v-model:value="form.identity" placeholder="登录用账号（3-32位）" />
    </a-form-item>
    <a-form-item label="昵称" name="name">
      <a-input v-model:value="form.name" placeholder="显示昵称" />
    </a-form-item>
    <a-form-item label="密码" name="password">
      <a-input-password v-model:value="form.password" placeholder="密码（6-64位）" />
    </a-form-item>
    <a-form-item label="确认密码" name="confirmPassword">
      <a-input-password v-model:value="form.confirmPassword" placeholder="再次输入密码" />
    </a-form-item>
    <a-form-item label="邮箱" name="email">
      <a-input v-model:value="form.email" placeholder="邮箱" />
    </a-form-item>
    <a-form-item label="邮箱验证码" name="email_code">
      <div class="code-row">
        <a-input v-model:value="form.email_code" placeholder="6 位验证码" :maxlength="6" />
        <a-button :disabled="countdown > 0" @click="sendCode">
          {{ countdown > 0 ? `${countdown}s 后重发` : '发送验证码' }}
        </a-button>
      </div>
    </a-form-item>
    <a-button type="primary" html-type="submit" block :loading="loading"> 注册 </a-button>
  </a-form>
</template>

<style scoped>
.code-row {
  display: flex;
  gap: 8px;
}
</style>
