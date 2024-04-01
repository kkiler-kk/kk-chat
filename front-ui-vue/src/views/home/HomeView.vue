<template>
  <div :style="parentDiv">
    <a-layout style="width: 100%; height: 100%">
      <a-layout-sider :style="siderStyle" width="450px"><Sider /></a-layout-sider>
      <a-layout>
        <a-layout-header :style="headerStyle">Header</a-layout-header>
        <a-layout-content :style="contentStyle">Content</a-layout-content>
        <a-layout-footer :style="footerStyle">Footer</a-layout-footer>
      </a-layout>
    </a-layout>
  </div>
</template>
<script setup lang="ts">
import type { CSSProperties } from "vue";
import Sider from "@/views/home/components/sider.vue";
import { createSocket, sendMsg } from "@/utils/websocket";
import { onMounted } from "vue";
import { message, Modal } from "ant-design-vue";
import { Session } from "@/utils/storage";
import router from "@/router";
const headerStyle: CSSProperties = {
  textAlign: "center",
  color: "#fff",
  //   paddingInline: 50,
  //   lineHeight: '64px',
  backgroundColor: "rgb(243, 242, 239)",
};
const parentDiv: CSSProperties = {
  width: "100%",
  height: "100%",
};
const contentStyle: CSSProperties = {
  textAlign: "center",
  lineHeight: "120px",
  color: "#fff",
  width: "100%",
  height: "100%",
  backgroundColor: "white",
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
  backgroundColor: "rgb(243, 242, 239)",
};
onMounted(() => {
  console.log("Session.get(\"token\")", Session.get("token"))
  if (Session.get("token") == undefined) {  // 检查是否登录
    Modal.warning({
      title: () => "提示",
      content: () => "登录状态已过期，请重新登录",
      onOk() {
        setTimeout(() => { // 重新登录 跳转登录路由
        router.push({
          name: "login",
        });
      }, 200);
      }
    });
    return;
  }
  createSocket(); // 连接websocket
  sendMsg("recentChatList", null); // 发送给客户端消息
  sendMsg("sendMessage", { id: 1, type: 1 }); // 发送给客户端消息
});
</script>
<style scoped></style>
