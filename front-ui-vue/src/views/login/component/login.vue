<template>
  <div class="box">
    <a-form
      :model="formState"
      name="basic"
      @finish="onFinish"
      @finishFailed="onFinishFailed"
    >
      <a-form-item
        label=""
        name="username"
        :rules="[{ required: true, message: '请输入邮箱/手机号/id' }]"
      >
        <a-input v-model:value="formState.username" class="input" placeholder="邮箱/手机号/ID" >
          <template #prefix><UserOutlined style="color: rgba(0, 0, 0, 0.25)" /></template>
        </a-input>
      </a-form-item>

      <a-form-item
        label=""
        name="password"
        :rules="[{ required: true, message: '请输入密码!' }]"
      >
        <a-input-password v-model:value="formState.password" class="input" placeholder="请输入密码" >
          <template #prefix><LockOutlined style="color: rgba(0, 0, 0, 0.25)" /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item
        label=""
        name="code.answer"
      >
        <a-input v-model:value="formState.code.answer" class="verificationCode" placeholder="请输入验证码" />
        <img style="margin-left: 10px" @click="getCaptchaGenerate" :src="formState.code.data" />
        <div>
        </div>
      </a-form-item>
      <a-form-item name="remember" >
        <a-checkbox v-model:checked="formState.remember"  style="float: left;">记住我</a-checkbox>
        <span @click="isLoginModel = !isLoginModel" class="text" style="float: right;">{{
          isLoginModel ? "短信验证登录" : "用户名密码登录"
        }}</span>
      </a-form-item>

      <a-form-item >
        <a-button type="primary" html-type="submit" style="width: 100%; margin: 0 auto;" :loading="loginLoading">登 录</a-button>
      </a-form-item>
      <span class="text" style="text-align: center;">已有账号，忘记密码?</span>
    </a-form>
    <a-divider>其他方式登录</a-divider>

  </div>
</template>
<script lang="ts" setup>
import { reactive, ref ,onMounted } from "vue";
import { UserOutlined, LockOutlined, SettingFilled } from '@ant-design/icons-vue';
import {captchaGenerate} from '@/api/auth/auth'
import {login} from '@/api/userBasic/userBasic'
import { Session } from '@/utils/storage';
import router from "@/router";
interface FormState {
  username: string;
  password: string;
  remember: boolean;
  code: CaptchaGenerateMode;
}
interface CaptchaGenerateMode {
  captchaId: string;
  data: string;
  answer: string;
}
const isLoginModel = ref(true);
const loginLoading = ref(false);
const formState = reactive<FormState>({
  username: "",
  password: "",
  remember: true,
  code: {
    captchaId: '',
    data: '',
    answer: ''
  }
});

const onFinish = (values: any) => {
  loginLoading.value = true
  login(formState).then((res: any) => {
    Session.set('token', res.token)
    Session.set('userInfo', res.user)  // 将信息载入浏览器本地缓存
    loginLoading.value = false
    setTimeout(() => {
        router.push({
          name: "home",
        });
      }, 200);
  }).catch((err) => {
    getCaptchaGenerate() // 刷新验证码
      loginLoading.value = false
  })
};

const onFinishFailed = (errorInfo: any) => {
  console.log("Failed:", errorInfo);
};
onMounted(() => {
  getCaptchaGenerate()
})
const getCaptchaGenerate = () => {
  captchaGenerate().then((res) => {
    formState.code = {
      captchaId: res.captchaId,
      data: res.data
    }
  })
}
</script>
<style  scoped>
.text {
  color: rgb(0, 89, 128);
  cursor: pointer;
}
.ant-form{
  width: 100%;
  height: 100%;
}
.box {
  width: 80%;
  position: relative;
  left: 10%;
}
.input {
  width: 90%;
}
.verificationCode {
  width: 50%;
}
</style>
