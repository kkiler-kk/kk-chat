<template>
  <div>
    <div class="container">
      <ul class="list-unstyled py-4" id="scrollIV">
        <li v-for="(message, index) in messageList" class="d-flex message">
          <div class="content">
            <div class="mr-lg-3 me-2">
              <a-avatar shape="square" class="avatar" :src="logo" size="large" />
            </div>
            <div class="message-body">
              <div class="name">KK 7:19PM</div>
              <div class="message-row">
                Hello WroldHello WroldHello WroldHello WroldHello WroldHello WroldHello
               
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from "vue";
import logo from "@/assets/images/chat.png";
import { UserInfo } from "@/store/userInfo";
interface Props {
  clickChat?: UserInfo;
}
interface Message {
  // 消息结构体
  id: string;
  content: string; // 发送内容
  sendId: number; // 发送人id
  sendUser: UserInfo; // 发送人信息
  type: number; // 消息类型 群聊 or 私聊
  createdAt: string; // 消息创建时间
}
defineProps<Props>();
const messageList = ref<Message[]>([{}, {}]);
</script>
<style scoped>
.container {
  width: 100%;
  overflow: auto;
  color: black;
  height: 100%;
}
.content {
  display: flex;
  width: 80%;
  overflow: auto;
  /* flex-wrap: wrap; */
}
.message {
  /* float: right; */
}
.message-body {
  height: auto;
}
.name {
  float: left;
  display: block;
  width: 100%;
}

.list-unstyled {
  padding-left: 10px;
  margin-top: 1.5rem !important;
  list-style: none;
}
.message-row {
  /* display: -ms-grid;
  display: grid;
  -ms-grid-columns: 1fr auto;
  grid-template-columns: 1fr auto;
  margin: 16px;
  height: 64px;
  border-radius: 16px;
  overflow: auto;
  background-color: black;
  color: black; */
  /* border: solid #f9465c 1px; */
  /* margin: 16px; */
  display: grid;
  margin-bottom: 16px;
  margin-top: 26px;
  border-radius: 12px;
  padding-left: 16px;
  overflow: auto;
  height: auto;
  font-size: 14px;
  background-color: rgb(204, 204, 204);
  color: rgb(70, 74, 77);
}
.avatar {
  width: 40px;
  height: 40px;
}
</style>
