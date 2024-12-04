<template>
    <div v-if="flag">
      <a-card hoverable style="width: 100%; margin: 0 auto">
          <template #cover>
            <img alt="example" :src="cover" />
          </template>
          <template #actions v-if="isSlfe">
            <setting-outlined key="setting" />
            <edit-outlined key="edit" @click="handleUpdateInfo" />
            <ellipsis-outlined key="ellipsis" />
          </template>
          <div class="info">
            <div class="div-image">
              <a-avatar class="image" :size="100" :src="userInfo?.avatar" />
              <h4>{{ userInfo?.name }}</h4>
            </div>
            <div>
              <div>
                <span class="left">
                  <UserOutlined />
                </span><span class="right">{{ userInfo?.identity }}</span>
              </div>
              <div>
                <span class="left">
                  <icon :style="{ fontSize: '32px', color: 'rgb(38, 38, 38)' }">
                    <template #component>
                      <svg t="1723729736938" class="icon" width="0.5em" height="0.5em" viewBox="0 0 1024 1024"
                        version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4282">
                        <path
                          d="M630.207802 39.808686A39.254933 39.254933 0 0 0 590.989785 0.553753L78.755979 0A78.755979 78.755979 0 0 0 0 78.755979v866.475737a78.755979 78.755979 0 0 0 78.755979 78.755978h866.475737a78.755979 78.755979 0 0 0 78.755978-78.755978V433.957747a39.377989 39.377989 0 0 0-39.377989-39.316461h-0.135362a39.377989 39.377989 0 0 0-39.377989 39.316461v511.273969H78.755979V78.755979l512.11075 0.319946a39.254933 39.254933 0 0 0 39.279544-39.267239z"
                          fill="#070707" p-id="4283"></path>
                        <path
                          d="M877.649242 89.757204l55.707549 55.70755-467.010647 467.010647-80.921768 25.226524 25.214218-80.934073L877.649242 89.757204m0-78.755978a78.755979 78.755979 0 0 0-55.682938 23.07304L354.955657 501.084914a78.755979 78.755979 0 0 0-19.504411 32.253034l-25.226524 80.934074a78.755979 78.755979 0 0 0 98.641863 98.617251l80.921768-25.23883a78.755979 78.755979 0 0 0 32.240729-19.492104L989.002812 201.147691a78.755979 78.755979 0 0 0 0-111.378181l-55.695244-55.695244a78.755979 78.755979 0 0 0-55.695243-23.07304z"
                          fill="#070707" p-id="4284"></path>
                      </svg>
                    </template>
                  </icon>
                </span>
                <span class="right">{{ userInfo?.personalitySignature == "" ? "暂无个人签名": userInfo?.personalitySignature }}</span>
              </div>
              <div>
                <span class="left">
                  <icon :style="{ fontSize: '32px', color: 'rgb(38, 38, 38)' }">
                    <template #component>
                      <svg t="1723730125589" class="icon" viewBox="0 0 1024 1024" version="1.1"
                        xmlns="http://www.w3.org/2000/svg" p-id="5305" width="0.5em" height="0.5em">
                        <path
                          d="M914.285714 1024h-804.571428a109.714286 109.714286 0 0 1 0-219.428571h804.571428a109.714286 109.714286 0 0 1 0 219.428571zM804.571429 731.428571h-7.387429c4.169143-11.556571 7.387429-23.625143 7.387429-36.571428a109.714286 109.714286 0 0 0-219.428572 0c0 12.946286 3.218286 25.014857 7.387429 36.571428H431.469714c4.169143-11.556571 7.387429-23.625143 7.387429-36.571428a109.714286 109.714286 0 0 0-219.428572 0c0 12.946286 3.218286 25.014857 7.387429 36.571428H219.428571a146.285714 146.285714 0 0 1-146.285714-146.285714V512a146.285714 146.285714 0 0 1 146.285714-146.285714h585.142858a146.285714 146.285714 0 0 1 146.285714 146.285714v73.142857a146.285714 146.285714 0 0 1-146.285714 146.285714zM582.802286 273.554286C583.460571 267.629714 585.142857 262.217143 585.142857 256 585.142857 195.437714 552.374857 146.285714 512 146.285714s-73.142857 49.152-73.142857 109.714286c0 6.217143 1.682286 11.629714 2.340571 17.554286A145.554286 145.554286 0 0 1 365.714286 146.285714C365.714286 65.462857 512 0 512 0s146.285714 65.462857 146.285714 146.285714a145.554286 145.554286 0 0 1-75.483428 127.268572z"
                          fill="" p-id="5306"></path>
                      </svg>
                    </template>
                  </icon>
                </span>
                <span class="right">{{ dayjs(userInfo?.birthDate).format("YYYY-MM-DD") }}</span>
              </div>
              <div v-if="userInfo?.isFriend">
                <span class="left">
                  <MailFilled />
                </span><span class="right">{{ userInfo?.email }}</span>
              </div>
              <div v-if="userInfo?.isFriend">
                <span class="left">
                  <PhoneFilled />
                </span><span class="right">{{ userInfo?.phone }}</span>
              </div>
            </div>
          </div>
        </a-card>
        <updateUserInfo ref="updateUserInfoRef" />
        
    </div>
</template>
<script setup lang="ts">
import cover from "@/assets/images/login.jpg";
import { onMounted } from "vue";
import {ref} from 'vue'
import updateUserInfo from "@/views/home/components/viewsSoder/components/updateUserInfo.vue";
import dayjs from "dayjs"
import router from "@/router";
import {detail} from "@/api/userBasic/userBasic"
import Icon, {
  UserOutlined,
  MailFilled,
  SettingOutlined,
  EditOutlined,
  EllipsisOutlined,
  PhoneFilled,
} from "@ant-design/icons-vue";
const props = defineProps({
    flag: Boolean,
    isSlfe: Boolean,
    id: Number,
})
const updateUserInfoRef = ref()
const userInfo = ref<any>()

onMounted(() => {
  getUserDetail()
})
const getUserDetail = () => {
  detail(props.id).then((res) => {
    userInfo.value = res
  }).catch((err) => {
    // 重新登录 跳转登录路由
    router.push({
        name: "login",
      });
  })
}
const handleUpdateInfo = () => {
  updateUserInfoRef.value.showModal()
}
</script>
<style>

.left {
  float: left;
}

.right {
  font-weight: normal;
}
</style>