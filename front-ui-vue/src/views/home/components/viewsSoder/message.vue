<template>
  <div style="height: 100%; width: 93%; margin: 0 auto">
    <div class="div1">
      <span class="title">聊天</span>
      <a-popover v-model:open="visible" class="botton" trigger="click">
        <template #content>
          <p @click="handleOpenPlus" style="cursor: pointer">加好友/加群</p>
          <p style="cursor: pointer">创建群聊</p>
        </template>
        <a-button type="primary">找朋友</a-button>
      </a-popover>
    </div>
    <div>
      <a-input-search
        size="large"
        class="search"
        v-model:value="value"
        placeholder="搜索..."
        enter-button
        @search="onSearch"
      />
    </div>
    <div class="content">
      <a-card v-for="(item, index) in data" class="card" @click="handleSelectChat(item)">
        <div class="media">
          <div class="media-avatar">
            <a-badge dot color="green">
              <a-avatar shape="square" class="avatar" size="large">
                <template #icon><UserOutlined /></template>
              </a-avatar>
            </a-badge>
          </div>
          <div style="margin-top: 10px;">
            <div  clas="d-flex align-items-center mb-1">
              <h6 class="text-truncate mb-0 me-auto">KK</h6>
              <p class="small text-muted text-nowrap ms-4 mb-0">11:08 am</p>
            </div>
            <div class="text-truncate">It is a long established fact that a reader w...</div>
          </div>
        </div>
      </a-card>
    </div>
    <plusFirendGroup ref="openPlusRef" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from "vue";
import { UserOutlined } from "@ant-design/icons-vue";
import { sendMsg } from "@/utils/websocket";
import plusFirendGroup from "@/views/home/components/viewsSoder/components/plusFirendGroup.vue";
import {UserInfo} from '@/store/userInfo'
const visible = ref<boolean>(false);
const value = ref<string>("");
const openPlusRef = ref();
const data: UserInfo[] = [
  
];
const clickChat = ref<UserInfo>()
onMounted(() => {
  sendMsg("recentChatList", null); // 发送给客户端消息 recentChatList 获取最近聊天信息列表 （用户列表）
});
const onSearch = (searchValue: string) => {
  console.log("use value", searchValue);
  console.log("or use this.value", value.value);
};
const handleOpenPlus = () => {
  openPlusRef.value.showModal();
};

const handleSelectChat = (item: UserInfo) => {
    clickChat.value = item
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
  /* margin-right: 10px; */
  background-color: rgb(24, 144, 255);
  color: white;
  border-radius: 8%;
}
.search {
  background-color: white;
  border-radius: 15%;
}
.content {
  background-color: white;
  height: 87%;
  margin-top: 10px;
}
.content {
  overflow: auto;
}
.small {
  float: right;
}
.card {
  height: 70px;
  padding: 0;
  border-radius: 3px;
  margin-bottom: 3px;
}
.card:hover {
  border: 1px solid rgb(24, 144, 255);
  box-shadow: 0.5px 0.5px rgb(24,144,255);
}
:deep(.card > .ant-card-body){
  padding: 0px;
}
.avatar {
  width: 40px;
  height: 40px;
}
.media-avatar{
  margin-right: 1rem !important;
  margin-top: 1rem !important;
}

</style>
