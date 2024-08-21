<template>
  <div style="height: 100%">
    <div class="div1">
      <span class="title">信 息</span>
      <a-button class="botton" @click="handlLogout">退 出</a-button>
    </div>
    <a-spin :spinning="spinning">
      <div class="userInfo" v-if="userInfo?.id"> 
        <CarUserInfo :is-slfe="true" :id="userInfo?.id" :flag="true" />
      </div>
    </a-spin>
    <div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { Local } from "@/utils/storage";
import { logout } from "@/api/userBasic/userBasic";
import { message } from "ant-design-vue";
import CarUserInfo from "@/components/card/userInfo.vue"
import router from "@/router";
const userInfo = ref<any>({
});
const spinning = ref<boolean>(true);

onMounted(() => {
  let user = Local.get("userInfo");
  spinning.value = true;
  if (user == null) {
    spinning.value = false;
    return;
  }
  userInfo.value = user
  userInfo.value.isFriend = true
  spinning.value = false
});
// 退出登录
const handlLogout = () => {
  logout().then((res) => {
    message.success("退出成功!");
    Local.clear();
  });
  setTimeout(() => {
    // 重新登录 跳转登录路由
    router.push({
      name: "login",
    });
  }, 200);
};

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

</style>
