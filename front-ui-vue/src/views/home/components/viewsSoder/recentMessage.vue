<template>
  <div style="height: 100%; width: 93%; margin: 0 auto">
    <div class="div1">
      <span class="title">聊天</span>
      <a-popover v-model:open="visible" class="botton" trigger="click">
        <template #content>
          <p @click="handleOpenPlus" style="cursor: pointer">加好友/加群</p>
          <p style="cursor: pointer" @click="handleCreateGroup">创建群聊</p>
        </template>
        <a-button type="primary">找朋友</a-button>
      </a-popover>
    </div>
    <div>
      <a-input-search
        size="large"
        class="search"
        v-model:value="searchValue"
        placeholder="搜索..."
        enter-button
        @search="onSearch"
      />
    </div>
    <div class="content">
      <a-card
        v-for="(item, i) in searchList"
        class="card"
        :class="{click: index == i}"
        @click="handleSelectChat(item, i)"
      >
        <div class="media">
          <div class="media-avatar">
            <a-badge dot :color="item.isLogout ? 'green' : 'rgb(204, 204, 204)'" v-if="item.type == 1">
              <a-avatar shape="square" class="avatar" size="large" :src="item.avatar" />
            </a-badge>
            <a-avatar v-else shape="square" class="avatar" size="large" :src="item.avatar" />
          </div>
          <div style="margin-top: 10px; position: relative;">
            <div clas="d-flex align-items-center mb-1" style="width: 100%;height: 100%;">
              <h6 class="text-truncate mb-0 me-auto">
                {{ item.name }}
              </h6>
              <p class="small text-muted text-nowrap ms-4 mb-0" style="position: absolute; left: 12rem;">
                {{ formatPast(new Date(item.createdAt), "YYYY-mm-dd HH:MM") }}
              </p>
            </div>
            <div class="text-truncate" style="margin-bottom: 0rem; position: absolute;top: 1.5rem;">
              {{ item.content }}
            </div>
          </div>
        </div>
      </a-card>
    </div>
    <plusFirendGroup ref="openPlusRef" />
    <createGroup  ref="openGroupRef"/>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, computed} from "vue";
import { UserOutlined } from "@ant-design/icons-vue";
import { sendMsg } from "@/utils/websocket";
import plusFirendGroup from "@/views/home/components/viewsSoder/components/plusFirendGroup.vue";
import createGroup from "@/views/home/components/viewsSoder/components/createGroup.vue";
import { message } from "ant-design-vue";
import { formatPast } from "@/utils/formatTime";
import { isJsonString } from "@/utils/is";
const visible = ref<boolean>(false);
const searchValue = ref<string>("");
const openPlusRef = ref();
const openGroupRef = ref();
const dataList = ref<any>([]); // 最近聊天消息列表
const index = ref<any>()
const clickChat = ref<any>();
onMounted(() => {
  sendMsg("recentChatList", null); // 发送给客户端消息 recentChatList 获取最近聊天信息列表 （用户列表）
  window.addEventListener("onmessageWS", getSocketData);
});
const emit = defineEmits(["getValue"]);
const onSearch = (searchValue: string) => {
  // dataList.value =  dataList.value.filter(function(item) {
  //   console.log("item: ", item)
  //   return item.name.indexOf(searchValue) > 0
  // })
};
const searchList = computed(() => {
  return dataList.value.filter((item) => {
    return item.name.toLowerCase().match(searchValue.value.toLowerCase())
  })
})
const handleOpenPlus = () => {
  openPlusRef.value.showModal();
};
function getSocketData(data) {
  const message = data.detail.event;
  if (message.event == "recentChatList") {
    dataList.value = message.data;
  }
}

// handlSelectChat 选择聊天对象
const handleSelectChat = (item: any, i: number) => {
  // 从子组件传值到父组件值
  index.value = i
  clickChat.value = item;
  clickChat.value.type = item.type
  emit("getValue", item);
};
const handleCreateGroup = () =>{
  openGroupRef.value.showModal()
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
  box-shadow: 0.5px 0.5px rgb(24, 144, 255);
}
.click {
  border: 1px solid rgb(24, 144, 255);
  box-shadow: 0.5px 0.5px rgb(24, 144, 255);
}
:deep(.card > .ant-card-body) {
  padding: 0px;
}
.avatar {
  width: 40px;
  height: 40px;
}
.media-avatar {
  margin-right: 1rem !important;
  margin-top: 1rem !important;
}
</style>
