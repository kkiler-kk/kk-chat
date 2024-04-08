<template>
  <div style="height: 100%; width: 93%; margin: 0 auto">
    <div class="div1">
      <span class="title">联系人</span>
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
      <div v-for="(item, o) in data" :key="index">
        <div style="float: left; "> <h6>{{ item.type }}</h6></div>
        <div v-for="(l, i) in item.list" :key="l.id" class="card"  :class="{click: index == l.id}" @click="handleSelectChat(l,l.id)">
          <div class="media">
            <div class="media-avatar">
              <a-badge dot :color="l.isLogout ? 'green' : 'rgb(204, 204, 204)'">
                <a-avatar shape="square" class="avatar" size="large" :src="l.avatar"/>
              </a-badge>
            </div>
            <div style="margin-top: 10px">
              <div clas="d-flex align-items-center mb-1">
                <h6 class="text-truncate mb-0 me-auto">{{ l.name }}</h6>
                <!-- <p class="small text-muted text-nowrap ms-4 mb-0">11:08 am</p> -->
              </div>
              <div class="text-truncate">
                {{ l.isLogout ? '' : formatPast(l.loginOutTime) }} 在线
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <plusFirendGroup ref="openPlusRef" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from "vue";
import { listFriend } from "@/api/userFriend/userFriend";
import pinyin from "js-pinyin";
import {formatPast} from '@/utils/formatTime'
import plusFirendGroup from "@/views/home/components/viewsSoder/components/plusFirendGroup.vue";

const clickChat = ref<any>()
const visible = ref<boolean>(false);
const value = ref<string>("");
const openPlusRef = ref();
const index = ref<number>()
const data = ref<any>();

const emit = defineEmits(["getValue"])

onMounted(() => {
  pinyin.setOptions({ charCase: 1 });
  handleListFriend();
});
const onSearch = (searchValue: string) => {
  console.log("use value", searchValue);
  console.log("or use this.value", value.value);
};
const handleListFriend = () => {
  listFriend().then((res: any) => {
    // console.log('res', res)
    // data.value = res
    let scenicData = [];
    for (let index = 0; index < res.length; index++) {
      let obj = res[index];
      obj.pin = pinyin.getFullChars(obj.name);
      scenicData.push(obj);
    }
    filterName(scenicData);
    data.value = filterName(scenicData);
  });
};
const handleOpenPlus = () => {
  openPlusRef.value.showModal();
};

// handlSelectChat 选择聊天对象
const handleSelectChat = (item :any, i: number) => {  // 从子组件传值到父组件值
  clickChat.value = item
  index.value = i
  emit("getValue", item)
}


const filterName = (list) => {
  let letterArray = [];
  for (let i = 65; i < 91; i++) {
    letterArray.push(String.fromCharCode(i));
  }
  let newNames = [];
  letterArray.forEach((item) => {
    newNames.push({
      type: `${item}`,
      list: list.filter((val) => {
        return val.pin.slice(0, 1).toUpperCase() === `${item}`;
      }),
    });
  });
  newNames = newNames.filter((item) => item.list.length);
  return newNames;
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
  background-color: rgb(243, 242, 239);
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
  background-color: white;
}
.card:hover {
  border: 1px solid rgb(24, 144, 255);
  box-shadow: 0.5px 0.5px rgb(24, 144, 255);
  color: rgb(24, 144, 255);
}
.click {
  border: 1px solid rgb(24, 144, 255);
  box-shadow: 0.5px 0.5px rgb(24, 144, 255);
}
:deep(.ant-card-body) {
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

h6 {
  font-weight: normal;
  font-size: 1rem;
}


</style>
