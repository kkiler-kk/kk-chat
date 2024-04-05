<template>
  <div>
    <a-row type="flex">
      <a-col :span="8">
        <div class="media">
          <div class="me-3 show-user-detail">
            <a-avatar :src="clickChat?.avatar" shape="square" class="avatar" size="large" />
          </div>
          <div>
            <div class="d-flex align-items-center mb-1"><h6>{{ clickChat?.name }}</h6></div>
          </div>
        </div>
      </a-col>
      <a-col :span="8" :offset="8">
      <!-- 暂定功能点 -->
      </a-col>
    </a-row>
  </div>
</template>
<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { UserInfo } from "@/store/userInfo";
interface Props {
  clickChat?: UserInfo
}
 defineProps<Props>()
// const clickChat = ref<UserInfo>(props.clickChat)
// onMounted(() => {
//   console.log("header: ", props.clickChat)
//   console.log("headerclickChat", clickChat.value)
// })
</script>
<style scoped>
* {
  color: black;
}

.avatar {
  width: 40px;
  height: 40px;
}
</style>
