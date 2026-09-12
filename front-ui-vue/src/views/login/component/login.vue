<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as authApi from '@/api/auth'
import { ApiError } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import { ErrCode } from '@/types/errorcode'
import type { FormInstance } from 'ant-design-vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const captchaId = ref('')
const captchaImage = ref('')

const form = reactive({
  account: '',
  password: '',
  captcha_answer: '',
})

const rules = {
  account: [{ required: true, message: '请输入账号或邮箱' }],
  password: [{ required: true, message: '请输入密码' }],
  captcha_answer: [{ required: true, message: '请输入图形验证码' }],
}

async function refreshCaptcha() {
  const res = await authApi.getCaptcha()
  captchaId.value = res.captcha_id
  captchaImage.value = res.image
  form.captcha_answer = ''
}

async function onSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    await auth.login(form.account, form.password, captchaId.value, form.captcha_answer)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/home'
    await router.push(redirect)
  } catch (e) {
    // 验证码错误/密码错误 → 刷新图形验证码
    if (
      e instanceof ApiError &&
      (e.code === ErrCode.CaptchaInvalid || e.code === ErrCode.BadCredentials)
    ) {
      await refreshCaptcha()
    }
  } finally {
    loading.value = false
  }
}

onMounted(refreshCaptcha)
</script>

<template>
  <a-form ref="formRef" :model="form" :rules="rules" layout="vertical" @finish="onSubmit">
    <a-form-item label="账号 / 邮箱" name="account">
      <a-input v-model:value="form.account" placeholder="identity 或邮箱" size="large" />
    </a-form-item>
    <a-form-item label="密码" name="password">
      <a-input-password v-model:value="form.password" placeholder="密码" size="large" />
    </a-form-item>
    <a-form-item label="图形验证码" name="captcha_answer">
      <div class="captcha-row">
        <a-input
          v-model:value="form.captcha_answer"
          placeholder="验证码"
          size="large"
          :maxlength="4"
        />
        <img
          v-if="captchaImage"
          :src="captchaImage"
          class="captcha-img"
          title="点击刷新"
          alt="captcha"
          @click="refreshCaptcha"
        />
      </div>
    </a-form-item>
    <a-button type="primary" html-type="submit" size="large" block :loading="loading">
      登录
    </a-button>
  </a-form>
</template>

<style scoped>
.captcha-row {
  display: flex;
  gap: 8px;
}

.captcha-img {
  height: 40px;
  width: 120px;
  border-radius: 6px;
  cursor: pointer;
  border: 1px solid #d9d9d9;
}
</style>
