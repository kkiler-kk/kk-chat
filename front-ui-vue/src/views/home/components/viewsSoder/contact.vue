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
      <a-input-search size="large" class="search" v-model:value="value" placeholder="搜索..." enter-button
        @search="onSearch" />
    </div>
    <div class="content">
      <div class="manage">
        <div>
          <span>好友通知</span>
          <span>
            <RightOutlined />
          </span>
        </div>
        <div>
          <span>群通知</span>
          <span>
            <RightOutlined />
          </span>
        </div>
      </div>
      <a-divider
        style="height: 0.5px; background-color: rgb(23, 133, 255); padding: 0; margin-top: 0px; margin-bottom: 12px;" />
      <div>
        <a-collapse v-model:activeKey="activeKey" :bordered="false">
          <a-collapse-panel key="1" header="群聊">
            <div v-for="(item, i) in searhcGroupList" class="card" :class="{ click: index == item.id }"
              @click="handleSelectChat(item, item.id, 2)">
              <div class="media">
                <div class="media-avatar">
                  <a-avatar shape="square" class="avatar" size="large" :src="item.avatar" />
                </div>
                <div style="margin-top: 10px">
                  <div clas="d-flex align-items-center mb-1">
                    <h6 class="text-truncate mb-0 me-auto">{{ item.name }}({{ item.identity }})</h6>
                    <!-- <p class="small text-muted text-nowrap ms-4 mb-0">11:08 am</p> -->
                  </div>
                  <div class="text-truncate">{{ item.countMember }} 成员</div>
                </div>
              </div>
            </div>
          </a-collapse-panel>
          <a-collapse-panel key="2" header="我的好友">
            <div v-for="(item, o) in searhcFriendList" :key="index">
              <div style="float: left">
                <h6>{{ item.type }}</h6>
              </div>
              <div v-for="(l, i) in item.list" :key="l.id" class="card" :class="{ click: index == l.id }"
                @click="handleSelectChat(l, l.id, 1)">
                <div class="media">
                  <div class="media-avatar">
                    <a-badge dot :color="l.isLogout ? 'green' : 'rgb(204, 204, 204)'">
                      <a-avatar shape="square" class="avatar" size="large" :src="l.avatar" />
                    </a-badge>
                  </div>
                  <div style="margin-top: 10px">
                    <div clas="d-flex align-items-center mb-1">
                      <h6 class="text-truncate mb-0 me-auto">{{ l.name }}</h6>
                      <!-- <p class="small text-muted text-nowrap ms-4 mb-0">11:08 am</p> -->
                    </div>
                    <div class="text-truncate">
                      {{ l.isLogout ? "" : formatPast(l.loginOutTime) }} 在线
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </a-collapse-panel>
        </a-collapse>
      </div>
    </div>
    <plusFirendGroup ref="openPlusRef" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { listFriend } from "@/api/userFriend/userFriend";
import pinyin from "js-pinyin";
import { listGroup } from "@/api/group/index";
import { formatPast } from "@/utils/formatTime";
import plusFirendGroup from "@/views/home/components/viewsSoder/components/plusFirendGroup.vue";
import {
  RightOutlined
} from "@ant-design/icons-vue";
const activeKey = ref(['2'])
const clickChat = ref<any>();
const visible = ref<boolean>(false);
const searchValue = ref<string>("");
const openPlusRef = ref();
const index = ref<number>();
const dataList = ref<any>([]);
const groupList = ref<any>();
const emit = defineEmits(["getValue"]);

onMounted(() => {
  pinyin.setOptions({ charCase: 1 });
  handleListFriend();
  handleListGroup();
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
    dataList.value = filterName(scenicData);
  });
};
const handleListGroup = () => {
  listGroup().then((res: any) => {
    groupList.value = res;
    console.log("groupList", groupList.value);
  });
};
const handleOpenPlus = () => {
  openPlusRef.value.showModal();
};

// handlSelectChat 选择聊天对象
const handleSelectChat = (item: any, i: number, type: number) => {
  // 从子组件传值到父组件值
  clickChat.value = item;
  clickChat.value.type = type
  index.value = i;
  emit("getValue", clickChat.value);
};

const filterName = (list: any) => {  // 按字母排序
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
const searhcGroupList = computed(() => { // 过滤群聊
  return groupList.value.filter((item) => {
    return item.name.toLowerCase().match(searchValue.value.toLowerCase())
  })
})

const searhcFriendList = computed(() => { // 过滤群聊
  console.log("dataList", dataList.value)
  return dataList.value.filter((item) => {
    return item.list.filter((n) => {
      return n.name.toLowerCase().match(searchValue.value.toLowerCase())
    })
  })
})
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

.manage {
  width: 100%;
  height: 9%;
  font-size: 15px;
  margin: 0px;
  padding: 0;
}

.manage div {
  margin: 0 auto;
  width: 90%;
  height: 45%;
}

.manage>div span:nth-child(1) {
  float: left;
}

.manage>div span:nth-child(2) {
  float: right;
}
</style>
