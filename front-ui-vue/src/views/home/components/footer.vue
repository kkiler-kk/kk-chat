<template>
  <div style="display: flex">
    <div style="width: 80%; margin-right: 20px;">
      <a-input
        v-model:value="formState.content"
        :bordered="false"
        placeholder="请输入消息..."
        @keyup.enter="handleSendMsg"
      />
    </div>
    <div style="margin-right: 20px;">
      <a-button type="primary" shape="circle">
        <template #icon><meh-outlined /></template>
      </a-button>
    </div>
    <div style="margin-right: 20px;">
      <a-button type="primary" @click="handleSendMsg">
        发送
        <template #icon><arrow-right-outlined /></template>
      </a-button>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { ref } from "vue";
import { MehOutlined,ArrowRightOutlined } from "@ant-design/icons-vue";
import { message } from "ant-design-vue";
import { sendMsg,getSocket } from "@/utils/websocket";
import { Session } from "@/utils/storage";
interface DataItem {
  content: string; // 发送内容
  userId: number; // 发送者id
  user: UserInfo,
  targetUser: UserInfo,
  targetId: number; //发送目标
  type: 1; // 类型 私聊或者群聊
}
interface UserInfo {
  id: number,
  name: string,
  avatar: string,
}
interface Props {
  clickChat?: any
}
const props =  defineProps<Props>()
 
const formState = ref<DataItem>({
  content: "",  // 发送内容
  userId: 0, // 发送人id
  targetId: 0, // 消息接收者id
  type: 1, // 类型 群聊 or 私聊 默认1 是私聊 因为群聊晚点写
});

const handleSendMsg = () => { // 发送消息
  if (props.clickChat == undefined) {
    message.error('请选择聊天对象!')
    return 
  }
  let userInfo = Session.get("userInfo")
  if (userInfo == undefined) {
    message.error("请先登录!")
    return
  }
  formState.value.user = {
    id: userInfo.id,
    name: userInfo.name,
    avatar: userInfo.avatar,
  }
  formState.value.targetUser = {
    id: props.clickChat.id,
    name: props.clickChat.name,
    avatar: props.clickChat.avatar,
  }
  formState.value.userId = userInfo.id
  formState.value.targetId = props.clickChat.id
  console.log('props.clickChat handleSendMsg', formState.value)
  sendMsg("sendMessage", formState.value)
}
// getSocket().onmessage = function (event) {
//   console.log("getSocket", event)
// }
</script>
<style scoped></style>
