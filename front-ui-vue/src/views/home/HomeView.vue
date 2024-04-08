<template>
  <div :style="parentDiv">
    <a-layout style="width: 100%; height: 100%">
      <a-layout-sider :style="siderStyle" width="450px"
        ><Sider @getValue="getClickChat"
      /></a-layout-sider>
      <a-layout>
        <a-layout-header :style="headerStyle"
          ><Header :clickChat="clickChat"
        /></a-layout-header>
        <!-- 相继传入 从 sider点击聊天对象 传入组件 -->
        <a-layout-content :style="contentStyle"
          ><Content :clickChat="clickChat"
        /></a-layout-content>
        <a-layout-footer :style="footerStyle"
          ><Footer :clickChat="clickChat"
        /></a-layout-footer>
      </a-layout>
    </a-layout>
  </div>
</template>
<script setup lang="ts">
import type { CSSProperties } from "vue";
import Sider from "@/views/home/components/sider.vue";
import Header from "@/views/home/components/header.vue";
import Content from "@/views/home/components/content.vue";
import Footer from "@/views/home/components/footer.vue";
import { createSocket, sendMsg } from "@/utils/websocket";
import { onMounted, ref } from "vue";
import { message } from "ant-design-vue";
import { Session } from "@/utils/storage";
import router from "@/router";
const headerStyle: CSSProperties = {
  textAlign: "center",
  color: "#fff",
  //   paddingInline: 50,
  //   lineHeight: '64px',
  height: "80px",
  backgroundColor: "rgb(243, 242, 239)",
};
const parentDiv: CSSProperties = {
  width: "100%",
  height: "100%",
  overflow: "hidden",
};
const contentStyle: CSSProperties = {
  width: "100%",
  height: "0",
  overflow: "auto",
  backgroundColor: "rgb(245, 246, 247)",
};

const siderStyle: CSSProperties = {
  textAlign: "center",
  color: "#fff",
  height: "100%",
  width: "600px",
  backgroundColor: "rgb(243, 242, 239)",
};

const footerStyle: CSSProperties = {
  textAlign: "center",
  color: "#fff",
  height: "100px",
  backgroundColor: "rgb(243, 242, 239)",
};
onMounted(() => {
  console.log("Download the sub-item address");
  console.log("https://github.com/KKiller-op/java-chat-room");
  console.log('Session.get("token")', Session.get("token"));
  if (Session.get("token") == undefined) {
    // 检查是否登录
    message.error("登录状态已过期，请重新登录");
    setTimeout(() => {
      // 重新登录 跳转登录路由
      router.push({
        name: "login",
      });
    }, 200);
    return;
  }
  createSocket(); // 连接websocket
});
const clickChat = ref<any>();
// 获取点击要聊天的对象
const getClickChat = (value: any) => {
  clickChat.value = value;
};
</script>
<style scoped></style>
