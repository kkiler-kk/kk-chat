<template>
  <div class="box">
    <a-form
      :model="formState"
      name="basic"
      autocomplete="off"
      @finish="onFinish"
      @finishFailed="onFinishFailed"
    >
      <a-form-item
        label="用户名"
        name="username"
        :rules="[{ required: true, message: '请输入用户名!' }]"
      >
        <a-input
          v-model:value="formState.username"
          class="input"
          placeholder="请输入用户名"
        />
      </a-form-item>
      <a-form-item
        label="ID"
        name="identity"
        :rules="[{ required: true, message: '请输入ID!' }]"
      >
        <a-input
          v-model:value="formState.identity"
          class="input"
          placeholder="请输入ID"
        />
      </a-form-item>
      <a-form-item
        label="邮箱"
        name="email"
        :rules="[{ required: true, message: '请输入邮箱!' }]"
      >
        <a-input
          v-model:value="formState.email"
          class="input"
          placeholder="请输入邮箱地址"
        />
      </a-form-item>

      <a-form-item
        label="验证码"
        name="verificationCode"
        :rules="[{ required: true, message: '请输入验证码!' }]"
      >
        <a-input
          v-model:value="formState.verificationCode"
          class="verificationCode"
          placeholder="请输入验证码"
        />
        <a-button style="margin-left: 10px" @click="sendEmail">发送</a-button>
      </a-form-item>
      <a-form-item
        label="密码"
        name="oldPassword"
        :rules="[{ required: true, message: '请输入密码!' }]"
      >
        <a-input-password
          v-model:value="formState.oldPassword"
          class="input"
          placeholder="请输入密码"
        />
      </a-form-item>
      <a-form-item
        label="密码"
        name="newPassword"
        :rules="[{ required: true, message: '请输入密码!' }]"
      >
        <a-input-password
          v-model:value="formState.newPassword"
          class="input"
          placeholder="请再次输入密码"
        />
      </a-form-item>

      <a-form-item :wrapper-col="{ offset: 8, span: 16 }">
        <a-button type="primary" html-type="submit" style="width: 50%">注 册</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>
<script lang="ts" setup>
import { reactive, ref } from "vue";
import { message } from 'ant-design-vue';
import { sendEmailCode } from "@/api/auth/auth.ts";
import {createUserBasic} from '@/api/userBasic/userBasic.ts'
interface FormState {
  username: string;
  oldPassword: string;
  newPassword: string;
  identity: string;
  email: string;
  verificationCode: string; // 邮箱验证码
}

const formState = reactive<FormState>({
  username: "",
  oldPassword: "",
  newPassword: "",
  identity: "",
  email: "",
  verificationCode: "",
});
const onFinish = (values: any) => {
  console.log("Success:", values);
  createUserBasic(formState).then((res) => {
    message.success("注册成功！")
  }).catch((err) => {
    console.log('error', err)
    message.error("添加失败", err)
  })
};

const onFinishFailed = (errorInfo: any) => {
};
const sendEmail = () => {
  sendEmailCode(formState.email)
    .then((res) => {
      message.success("发送成功！请留意邮箱")
    })
    .catch((error) => {
      console.log('err', error);
    });
};
</script>
<style scoped>
.text {
  color: rgb(0, 89, 128);
}
.ant-form {
  width: 100%;
  height: 100%;
}
.box {
  width: 80%;
  position: relative;
  left: 10%;
}
.input {
  width: 80%;
}
.verificationCode {
  width: 50%;
}
</style>
