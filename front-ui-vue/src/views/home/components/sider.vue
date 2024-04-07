<template>
  <div style="height: 100%;">
    <div style="width: 15%; float:left;">
      <p>
        <a-image :src="logo" :width="44" />
      </p>
      <ul>
        <li>
          <a-avatar :src="userInfo?.avatar" size="large" @click="switComponents(UserInfo)" />
        </li>
        <li>
          <a-avatar size="large" class="li click" @click="switComponents(Message)">
            <template #icon><MessageFilled /></template>
          </a-avatar>
        </li>
        <li>
          <a-avatar size="large" class="li" @click="switComponents(Contact)">
            <template #icon><ContactsFilled /></template>
          </a-avatar>
        </li>

        <li>
          <a-avatar size="large" class="li" @click="switComponents(Like)">
            <template #icon><fire-filled /></template>
          </a-avatar>
        </li>
        <li>
          <a-avatar size="large" class="li" @click="switComponents(Setting)">
            <template #icon><SettingFilled /></template>
          </a-avatar>
        </li>
      </ul>
    </div>
    <div class="componet">
      <keep-alive include="Message" exclude="Setting">
        <component :is="whichComponent" @getValue="getValue"></component>
      </keep-alive>
    </div>
  </div>
</template>
<script setup lang="ts">
import logo from "@/assets/images/chat.png";
import {
  SettingFilled,
  LikeFilled,
  UserOutlined,
  FireFilled,
  ContactsFilled,
  MessageFilled,
} from "@ant-design/icons-vue";
import UserInfo from "@/views/home/components/viewsSoder/userInfo.vue";
import Setting from "@/views/home/components/viewsSoder/setting.vue";
import Like from "@/views/home/components/viewsSoder/like.vue";
import Message from "@/views/home/components/viewsSoder/message.vue";
import Contact from "@/views/home/components/viewsSoder/contact.vue";
import { ref, onMounted,watch } from "vue";
import { Session } from "@/utils/storage";
const whichComponent = ref(UserInfo);
const emit = defineEmits(['getValue'])
const userInfo = ref<any>({
  avatar: 'http://127.0.0.1:9345/resource/public/file/20240403/douwei.jpg'
})
const clickChat = ref<any>({name: 'hello'})
onMounted(() => {
  userInfo.value = Session.get('userInfo')
  if (userInfo.value == null){
    return 
  }
  clickChat.value = new Object
})
function switComponents(componet) {
  whichComponent.value = componet;
}
const getValue = (value: any) => {  // 从子组件传出 当前点击的值
  clickChat.value = value
  emit('getValue', clickChat.value) // 传入父级组件 HomeView
}
// let stopWatch = watch(clickChat, (newVal, oldVal) => {
//   console.log('newVal', newVal)
//   console.log('oldVal', oldVal)
// })
</script>
<style scoped>
div {
  color: black;
}
div > * {
  padding: 0;
  margin: 0;
}
ul {
  padding-bottom: 20px;
}
p {
  width: 90px;
}
ul li {
  list-style-type: none;
  margin-top: 20px;
  width: 90px;
}
.li:hover {
  color: rgb(24, 144, 255);
  background-color: "rgb(24, 144, 255)";
}
.click {
  color: rgb(24, 144, 255);
  background-color: "rgb(24, 144, 255)";
}
.componet{
    width: 85%;
    height: 100%;
    border: 1px solid white;
    float: right;
}
</style>
