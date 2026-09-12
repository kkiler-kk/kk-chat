<script setup lang="ts">
import { ref } from 'vue'
import LoginPanel from './component/login.vue'
import RegisterPanel from './component/register.vue'

const activeTab = ref<'login' | 'register'>('login')

function onRegistered() {
  activeTab.value = 'login'
}
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="brand">kk-chat</h1>
      <p class="slogan">简约 · 即时 · 畅聊</p>
      <a-tabs v-model:active-key="activeTab" centered>
        <a-tab-pane key="login" tab="登录">
          <LoginPanel />
        </a-tab-pane>
        <a-tab-pane key="register" tab="注册">
          <RegisterPanel @registered="onRegistered" />
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  padding: 32px 36px 24px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 12px 40px rgb(0 0 0 / 18%);
}

.brand {
  margin: 0;
  text-align: center;
  font-size: 30px;
  font-weight: 700;
  letter-spacing: 1px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.slogan {
  margin: 6px 0 12px;
  text-align: center;
  color: #999;
  font-size: 13px;
}
</style>
