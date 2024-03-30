<template>
  <div class="formBox">
    <div class="image"></div>
    <div class="form">
      <a-tabs v-model:activeKey="activeKey">
        <a-tab-pane key="login" tab="登录">
        <Login />
        </a-tab-pane>
        <a-tab-pane key="register" tab="注册" force-render><Register /></a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import Register from "@/views/login/component/register.vue";
import Login from "@/views/login/component/login.vue";
const activeKey = ref("login");
</script>
<style scoped>
.form {
  width: 45%;
}
.image{
  width: 55%;
  height: 100%;
  background-image: url('@/assets/images/login.jpg');
  background-size: cover;
  float: left;
  margin-right: 10px;
}
.formBox {
  width: 50%;
  height: 500px;
  box-shadow: 2px 2px 2px 1px rgba(0, 0, 0, 0.2);
  margin: 0 auto;
  display: flex;
  position: absolute;
  top: 25%;
  left: 25%;
}
.ant-tabs {
  width: 100%;
  height: 100%;
}

</style>
