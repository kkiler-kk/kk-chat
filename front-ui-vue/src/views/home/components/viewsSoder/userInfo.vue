<template>
  <div style="height: 100%">
    <div class="div1">
      <span class="title">信 息</span>
      <a-button class="botton" @click="handlLogout">退 出</a-button>
    </div>
    <a-spin :spinning="spinning">
      <div class="userInfo">
        <a-card hoverable style="width: 100%; margin: 0 auto">
          <template #cover>
            <img alt="example" :src="cover" />
          </template>
          <template #actions>
            <setting-outlined key="setting" />
            <edit-outlined key="edit" @click="handleUpdateInfo"/>
            <ellipsis-outlined key="ellipsis" />
          </template>
          <div>
            <div class="div-image">
              <a-avatar
                class="image"
                :size="100"
                :src="userInfo?.avatar"
              />
              <h4>{{ userInfo.name }}</h4>
            </div>
            <div>
              <div>
                <span class="left"><UserOutlined /> </span
                ><span class="right">{{ userInfo.identity }}</span>
              </div>
              <div>
                <span class="left"><MailFilled /></span
                ><span class="right">{{ userInfo.email }}</span>
              </div>
              <div>
                <span class="left"><PhoneFilled /></span
                ><span class="right">{{ userInfo.phone }}</span>
              </div>
            </div>
          </div>
        </a-card>
      </div>
    </a-spin>
    <div>
  </div>
  </div>
  <updateUserInfo ref="updateUserInfoRef"/>
</template>
<script setup lang="ts">
import cover from "@/assets/images/login.jpg";
import { onMounted } from "vue";
import { Session } from "@/utils/storage";
import { logout } from "@/api/userBasic/userBasic";
import { ref } from "vue";
import { message } from "ant-design-vue";
import router from "@/router";
import {
  UserOutlined,
  MailFilled,
  SettingOutlined,
  EditOutlined,
  EllipsisOutlined,
  PhoneFilled,
} from "@ant-design/icons-vue";
import updateUserInfo from "@/views/home/components/viewsSoder/components/updateUserInfo.vue";
const userInfo = ref<any>({});
const spinning = ref<boolean>(true);
const updateUserInfoRef = ref()
onMounted(() => {
  let user = Session.get("userInfo");
  console.log('user: ', user)
  spinning.value = true;
  if (user == null) {
    spinning.value = false;
    return;
  }
  userInfo.value = user
  spinning.value = false
});
// 退出登录
const handlLogout = () => {
  logout().then((res) => {
    message.success("退出成功!");
    Session.set("token", "");
  Session.set("userInfo", "");
  });
  setTimeout(() => {
    // 重新登录 跳转登录路由
    router.push({
      name: "login",
    });
  }, 200);
};

const handleUpdateInfo =() => {
  updateUserInfoRef.value.showModal()
}

</script>
<style scoped>
.div1 {
  margin-top: 10px;
  width: 100%;
  height: 5%;
}
.title {
  float: left;
  padding-left: 10px;
  font-size: 26px;
  color: rgb(24, 144, 255);
}
.botton {
  float: right;
  margin-right: 10px;
  background-color: rgb(24, 144, 255);
  color: white;
  border-radius: 8%;
}
.userInfo {
  width: 95%;
  height: 90%;
  background-color: white;
  margin: 0 auto;
}
.div-image {
  margin-top: -50px;
  margin-bottom: 20px;
}
.image {
  border-radius: 50%;
}
.left {
  float: left;
}
.right {
  font-weight: normal;
}
</style>
