<template>
  <div class="box">
    <a-modal v-model:visible="open" title="修改信息" @ok="handleOk">
      <a-form
        :model="formState"
        name="basic"
        ref="formRef"
        autocomplete="off"
        @finish="handleOk"
        @finishFailed="onFinishFailed"
      >
        <a-form-item label="头像">
          <a-upload
            v-model:file-list="fileList"
            name="avatar"
            list-type="picture-card"
            class="avatar-uploader"
            :show-upload-list="false"
            :before-upload="beforeUpload"
            @change="handleChange"
            :customRequest="handleUpload"
          >
            <a-avatar v-if="formState.avatar" :src="formState.avatar" alt="avatar" :size="100" />
            <div v-else>
              <loading-outlined v-if="loading"></loading-outlined>
              <plus-outlined v-else></plus-outlined>
              <div class="ant-upload-text">Upload</div>
            </div>
          </a-upload>
        </a-form-item>
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
          :rules="[
            {
              required: true,
              message: '邮箱格式不正确!',
              pattern: /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$/,
            },
          ]"
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
        >
          <a-input
            v-model:value="formState.verificationCode"
            class="verificationCode"
            placeholder="请输入验证码"
          />
          <a-button style="margin-left: 10px" @click="sendEmail">发送</a-button>
        </a-form-item>
        <a-form-item
          label="手机号"
          name="phone"
          :rules="[
            {
              pattern: /^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$/,
              message: '请输入正确手机号',
            },
          ]"
        >
          <a-input
            v-model:value="formState.phone"
            class="input"
            placeholder="请输入手机号"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
<script lang="ts" setup>
import { reactive, ref, onMounted } from "vue";
import { message } from "ant-design-vue";
import { sendEmailCode } from "@/api/auth/auth.ts";
import { updateUserBasic } from "@/api/userBasic/userBasic.ts";
import type { UploadChangeParam, UploadProps } from "ant-design-vue";
import { uploadImages } from "@/api/file/file";
import { Session } from "@/utils/storage";
import { detail } from "@/api/userBasic/userBasic";
const open = ref<boolean>(false);
interface FormState {
  id: number;
  avatar: string;
  username: string;
  identity: string;
  email: string;
  phone: string;
  verificationCode: string; // 邮箱验证码
}

const formState = reactive<FormState>({
  id: 0,
  avatar: "",
  username: "",
  identity: "",
  email: "",
  phone: "",
  verificationCode: "",
});
const formRef = ref(null)
const fileList = ref([]);
const url = ref(import.meta.env.VITE_API_URL + "api/app/file/upload/image");
const loading = ref<boolean>(false);

onMounted(() => {
});


const onFinishFailed = (errorInfo: any) => {};

const handleOk = (e: MouseEvent) => {
  formRef.value.validateFields()
  let userInfo = Session.get("userInfo");
  formState.id = userInfo.id;
  updateUserBasic(formState)
    .then((res) => {
      message.success("修改成功！");
      open.value = false
    })
    .catch((err) => {
      console.log("error", err);
      if (err == undefined) {
        message.success("修改成功") // 修改头像了
        return
      }
      message.error("修改失败", err);
    });
  open.value = false;
};
const beforeUpload = (file: UploadProps["fileList"][number]) => {
  const isJpgOrPng = file.type === "image/jpeg" || file.type === "image/png";
  if (!isJpgOrPng) {
    message.error("You can only upload JPG file!");
  }
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (!isLt2M) {
    message.error("Image must smaller than 2MB!");
  }
  return isJpgOrPng && isLt2M;
};
// 处理上传
const handleUpload = async (e) => {
  console.log("files", e.flle);
  // 设置头像上传状态为 true
  loading.value = true;

  // 发起上传请求
  try {
    const res = await uploadImages(e.file);
    formState.avatar = res;
  } catch (e: any) {
    message.error("上传失败" + e.message);
  }
  // 设置头像上传状态为 false
  loading.value = false;
};
const sendEmail = () => {
  sendEmailCode(formState.email)
    .then((res) => {
      message.success("发送成功！请留意邮箱");
    })
    .catch((error) => {
      console.log("err", error);
    });
};
const handleChange = (info: UploadChangeParam) => {
  if (info.file.status === "uploading") {
    loading.value = true;
    return;
  }
  if (info.file.status === "done") {
    // Get this url from response in real world.
    console.log("info.file", info.file);
  }
  if (info.file.status === "error") {
    loading.value = false;
    message.error("upload error");
  }
};
const showModal = () => {
  open.value = true;
  let userInfo = Session.get("userInfo");
  if (userInfo == null) {
    return
  }
  formState.id = userInfo.id;
  detail(formState.id.toString())
    .then((res) => {
      formState.id = res.id;
      formState.username = res.name;
      formState.phone = res.phone;
      formState.identity = res.identity;
      formState.email = res.email;
      formState.avatar = res.avatar;
    })
    .catch((err) => {
    });
};
const hideModal = () => {
  open.value = false;
};
//暴露state和play方法
defineExpose({
  showModal,
  hideModal,
});
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
