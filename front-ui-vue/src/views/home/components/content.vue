<template>
  <div>
    <div class="container">
      <ul class="list-unstyled py-4" id="scrollIV" ref="scrollIV">
        <li v-for="(message, index) in messageList" class="d-flex message">
          <div class="content" :class="{ self: userInfo.id == message?.user?.id }">
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
                <!--  v-if="message.type != 1" -->
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
import { ref, onMounted, watch, nextTick } from "vue";
import { sendMsg, getSocket } from "@/utils/websocket";
import logo from "@/assets/images/chat.png";
import { isJsonString } from "@/utils/is";
import { formatPast } from "@/utils/formatTime";
import { Session } from "@/utils/storage";
import { message } from "ant-design-vue";

interface Props {
  clickChat?: any;
}

const props = defineProps<Props>();
const messageList = ref<any>([]);
const userInfo = ref<any>();
const scrollIV = ref<HTMLElement | null>(null);
onMounted(() => {
  userInfo.value = Session.get("userInfo");
  if (userInfo == undefined) {
    message.error("请先登录!");
    return;
  }
  console.log("scrollIV", scrollIV.value.scrollTop);
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

async function getSocketData(data) {
  const messageModel = data.detail.event;
  if (messageModel.event == "messageList") {
    messageList.value = messageModel.data;
  } else if (messageModel.event == "sendMessage") {
    if (
      props.clickChat.id == messageModel.data.targetId ||
      messageModel.data.userId == props.clickChat.id
    ) {
      // 如果刚好和他聊天 就添加进当前消息聊天记录
      messageList.value.push(messageModel.data);
    }
  } else if (messageModel.event == "error") {
    message.error(messageModel.errorMsg);
  }
}

// onMounted(() => {
//   nextTick(() => {
//     console.log("scrollIV", scrollIV.value.scrollHeight);
//     scrollIV.value.scrollTo({ top: scrollIV.value.scrollHeight, behavior: "smooth" });
//   });
// })
</script>
<style scoped>
.container {
  width: 100%;
  overflow-y: auto;
  height: 100%;
}
.content {
  display: flex;
  width: 100%;
  margin-bottom: 15px;
  /* flex-wrap: wrap; */
}
.message {
  width: 100%;
  display: flex;
  padding-left: 0;
  justify-content: space-between;
}
.message-body {
}
.name {
  width: 100%;
}

.list-unstyled {
  margin-top: 1rem !important;
  list-style: none;
  overflow: auto;
}
.message-row {
  /* margin-bottom: 16px; */
  margin-left: 5px;
  border-radius: 0 32px 32px 32px;
  padding-left: 8px;
  padding-right: 8px;
  height: auto;
  position: relative;
  line-height: 30px;
  font-size: 14px;
  min-height: 32px;
  top: 0px;
  background-color: white;
  color: rgb(70, 74, 77);
}
.message-row::after {
  content: "";
  display: block;
  position: absolute;
  top: 0;
  left: -5px;
  border-left: 3px solid transparent;
  border-top: 3px solid white;
  border-right: 3px solid white;
  border-bottom: 3px solid transparent;
}
.avatar {
  width: 40px;
  height: 40px;
}
.self {
  /* float: right; */
  flex-direction: row-reverse;
  margin-left: 20px;
}
.self  .message-row {
  border-radius: 0px 32px 32px 32px;
}
.self::after {
  content: "";
  display: block;
  position: absolute;
  top: 0;
  left: 0px;
  border-left: 0 solid transparent;
  border-top: 0 solid rgb(204, 204, 204);
  border-right: 0 solid rgb(204, 204, 204);
  border-bottom: 0 solid transparent;
}
</style>
