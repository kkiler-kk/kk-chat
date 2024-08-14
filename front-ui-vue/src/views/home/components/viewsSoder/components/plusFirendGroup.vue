<template>
  <div>
    <a-modal v-model:visible="open" @ok="handleOk">
      <a-tabs v-model:activeKey="form.type">
        <a-tab-pane :key="1" tab="找好友"
          ><div>
            <a-input-search
              size="large"
              class="search"
              v-model:value="form.search"
              placeholder="id/手机号/邮箱搜索..."
              enter-button
              @search="onSearch"
            />
          </div>
          <div>
            <a-empty v-if="dataList.length == 0" />
            <div v-else>
              <div class="friendBox">
                <div v-for="item in dataList" :key="item.id" class="friend">
                  <div class="left">
                    <a-avatar :src="item.avatar" shape="square" :size="64" />
                  </div>
                  <div class="right">
                    <p>
                      <span
                        >{{ item.name }} (<span style="color: red">{{
                          item.identity
                        }}</span
                        >)</span
                      >
                    </p>
                    <p v-if="!item.isFriend">
                      <a-button type="primary" @click="handleAddFriend(item)"
                        >加好友</a-button
                      >
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </a-tab-pane>
        <a-tab-pane :key="2" tab="找群" force-render
          ><div>
            <a-input-search
              size="large"
              class="search"
              v-model:value="form.search"
              placeholder="群ID/群昵称"
              enter-button
              @search="onSearch"
            />
          </div>
          <div>
            <a-empty v-if="dataList.length == 0" />
            <div v-else>
              <div class="friendBox">
                <div v-for="item in dataList" :key="item.id" class="friend">
                  <div class="left">
                    <a-avatar :src="item.avatar" shape="square" :size="64" />
                  </div>
                  <div class="right">
                    <p>
                      <span
                        >{{ item.name }} (<span style="color: red">{{
                          item.identity
                        }}</span
                        >)</span
                      >
                    </p>
                    <p v-if="!item.isFriend">
                      <a-button type="primary" @click="handleAddFriend(item)"
                        >加群</a-button
                      >
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
  </div>
</template>
<script setup lang="ts">
import { Modal, message } from "ant-design-vue";
import { ref, createVNode } from "vue";
import { addFriend } from "@/api/userFriend/userFriend";
import { findByName } from "@/api/group";
import { ExclamationCircleOutlined } from "@ant-design/icons-vue";
import { userSearch } from "@/api/userBasic/userBasic";
import { Local } from "@/utils/storage";
const open = ref<boolean>(false);
const form = ref({
  type: 1, // 1找好友 2 找群
  search: "",
});
const dataList = ref<any>([]);
const showModal = () => {
  open.value = true;
};
const hideModal = () => {
  open.value = false;
};
const handleOk = (e: MouseEvent) => {
  console.log(e);
  open.value = false;
};
const handleAddFriend = (item) => {
  let userInfo = Local.get("userInfo");
  let data = {
    userId: userInfo.id,
    friendId: item.id,
  };
  addFriend(data)
    .then((res) => {
      message.success("添加成功");
    })
    .catch((err) => {
      message.error("添加失败", err);
    });
};
const onSearch = (searchValue: string) => {
  console.log("use value", searchValue);
  if (form.value.type == 1) {
    // 1 找好友
    userSearch(form.value)
      .then((res) => {
        dataList.value = res;
      })
      .catch((err) => {});
  } else {
    // 2 找群聊
    findByName(form.value.search).then((res) => {
      dataList.value = res;
    });
  }
};
//暴露state和play方法
defineExpose({
  showModal,
  hideModal,
});
</script>
<style scoped>
.friendBox {
  width: 100%;
  display: flex;
  margin-top: 5px;
  flex-wrap: wrap;
}
.friend {
  /* background-color: rgb(243, 242, 239); */
  width: 50%;
  /* border: 1px solid red; */
}
.left {
  float: left;
  display: inline;
}
.right {
  /* margin-right: 10px; */
}
</style>
