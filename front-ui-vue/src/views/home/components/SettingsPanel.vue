<script setup lang="ts">
import { onUnmounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { LogoutOutlined } from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { sendEmailCode } from '@/api/auth'
import AvatarUpload from '@/components/AvatarUpload.vue'

const auth = useAuthStore()
const router = useRouter()

const form = reactive({
  name: '',
  phone: '',
  signature: '',
  birth_date: '',
  email: '',
  email_code: '',
})
const avatarUrl = ref('')
const saving = ref(false)
const countdown = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

// 用户信息变化时回填表单
watch(
  () => auth.user,
  (u) => {
    if (!u) return
    form.name = u.name ?? ''
    form.phone = u.phone ?? ''
    form.signature = u.signature ?? ''
    form.birth_date = u.birth_date ? u.birth_date.slice(0, 10) : ''
    form.email = u.email ?? ''
    avatarUrl.value = u.avatar ?? ''
  },
  { immediate: true },
)

async function sendCode() {
  if (!form.email) {
    message.warning('请先填写新邮箱')
    return
  }
  await sendEmailCode(form.email)
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

async function save() {
  saving.value = true
  try {
    const emailChanged = auth.user && form.email !== auth.user.email
    await auth.updateProfile({
      name: form.name,
      phone: form.phone,
      signature: form.signature,
      birth_date: form.birth_date || undefined,
      avatar: avatarUrl.value,
      email: emailChanged ? form.email : undefined,
      email_code: emailChanged ? form.email_code : undefined,
    })
    message.success('资料已更新')
  } finally {
    saving.value = false
  }
}

async function logout() {
  await auth.logout()
  await router.push('/login')
}

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="settings-panel">
    <div class="avatar-row">
      <AvatarUpload :src="avatarUrl" @uploaded="(url) => (avatarUrl = url)" />
      <div class="identity">
        <div class="id-name">{{ auth.user?.name }}</div>
        <div class="id-account">账号：{{ auth.user?.identity }}</div>
      </div>
    </div>

    <a-form layout="vertical" class="form">
      <a-form-item label="昵称">
        <a-input v-model:value="form.name" :maxlength="64" />
      </a-form-item>
      <a-form-item label="手机号">
        <a-input v-model:value="form.phone" :maxlength="32" />
      </a-form-item>
      <a-form-item label="个性签名">
        <a-input v-model:value="form.signature" :maxlength="255" />
      </a-form-item>
      <a-form-item label="生日">
        <input v-model="form.birth_date" type="date" class="date-input" />
      </a-form-item>
      <a-form-item label="邮箱（修改需验证码）">
        <a-input v-model:value="form.email" />
      </a-form-item>
      <a-form-item v-if="auth.user && form.email !== auth.user.email" label="邮箱验证码">
        <div class="code-row">
          <a-input v-model:value="form.email_code" :maxlength="6" placeholder="6 位验证码" />
          <a-button :disabled="countdown > 0" @click="sendCode">
            {{ countdown > 0 ? `${countdown}s 后重发` : '发送验证码' }}
          </a-button>
        </div>
      </a-form-item>
      <a-button type="primary" block :loading="saving" @click="save"> 保存 </a-button>
      <a-button block danger class="logout" @click="logout">
        <template #icon><LogoutOutlined /></template>
        退出登录
      </a-button>
    </a-form>
  </div>
</template>

<style scoped>
.settings-panel {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.avatar-row {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 20px;
}

.id-name {
  font-size: 16px;
  font-weight: 600;
}

.id-account {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}

.date-input {
  width: 100%;
  padding: 6px 10px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  font-size: 14px;
}

.code-row {
  display: flex;
  gap: 8px;
}

.logout {
  margin-top: 12px;
}
</style>
