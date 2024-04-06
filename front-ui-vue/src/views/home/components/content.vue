<template>
  <div>
    <div class="container">
      <ul class="list-unstyled py-4" id="scrollIV">
        <li
          v-for="(message, index) in messageList"
          class="d-flex message"
         
        >
          <div class="content"  :class="{ self: userInfo.id == message?.user?.id }">
            <div class="mr-lg-3 me-2">
              <a-avatar
                shape="square"
                class="avatar"
                :src="message?.user?.avatar"
                size="large"
              />
            </div>
            <div class="message-body">
              <div class="name">
                {{ message?.user?.name }} {{ formatPast(new Date(message?.createTime)) }}
              </div>
              <div class="message-row">
                {{ message.content }}
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { sendMsg, getSocket } from "@/utils/websocket";
import logo from "@/assets/images/chat.png";
import { isJsonString } from "@/utils/is";
import { formatPast } from "@/utils/formatTime";
import { Session } from "@/utils/storage";
import { message } from "ant-design-vue";
import { debug } from "console";
interface Props {
  clickChat?: any;
}

const props = defineProps<Props>();
const messageList = ref<any>([]);
const userInfo = ref<any>();
onMounted(() => {
  userInfo.value = Session.get("userInfo");
  if (userInfo == undefined) {
    message.error("请先登录!");
    return;
  }
  window.addEventListener("onmessageWS", getSocketData);
});
// 监听 点击对象值的变换 请求获取历史消息
watch(
  () => props.clickChat,
  (newVal, oldVal) => {
    let data = {
      targetId: newVal.id,
      userId: userInfo.value.id,
    };
    sendMsg("messageList", data); // 发送给客户端消息 recentChatList 获取最近聊天信息列表 （用户列表）
  }
);
function getSocketData(data) {
  const message = data.detail.event;
  if (message.event == "messageList") {
    messageList.value = message.data;
    console.log("messageList", messageList.value);
  }
}
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
  width: 100%;
  overflow: auto;
  margin-bottom: 10px;
  /* flex-wrap: wrap; */
}
.message {
  width: 100%;
  display: flex;
  justify-content: space-between;
}
.message-body {
  height: auto;
}
.name {
  /* float: left; */
  /* display: block; */
  width: 100%;
}

.list-unstyled {
  margin-top: 1rem !important;
  list-style: none;
}
.message-row {
  /* margin-bottom: 16px; */
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
.self {
  align-self: flex-end;
}
</style>
