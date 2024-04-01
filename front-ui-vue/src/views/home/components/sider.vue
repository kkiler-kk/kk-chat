<template>
  <div style="height: 100%;">
    <div style="width: 15%;float:left;">
      <p>
        <a-image :src="logo" :width="44" />
      </p>

      <ul>
        <li>
          <a-avatar :src="logo" size="large" @click="switComponents(UserInfo)" />
        </li>
        <li>
          <a-avatar size="large" class="li" @click="switComponents(Message)">
            <template #icon><MessageFilled /></template>
          </a-avatar>
        </li>
        <li>
          <a-avatar size="large" class="li" @click="switComponents(Contact)">
            <template #icon><UserOutlined /></template>
          </a-avatar>
        </li>

        <li>
          <a-avatar size="large" class="li" @click="switComponents(Like)">
            <template #icon><LikeFilled /></template>
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
        <component :is="whichComponent"></component>
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
  ContactsFilled,
  MessageFilled,
} from "@ant-design/icons-vue";
import UserInfo from "@/views/home/components/viewsSoder/userInfo.vue";
import Setting from "@/views/home/components/viewsSoder/setting.vue";
import Like from "@/views/home/components/viewsSoder/like.vue";
import Message from "@/views/home/components/viewsSoder/message.vue";
import Contact from "@/views/home/components/viewsSoder/contact.vue";
import { ref } from "vue";
const whichComponent = ref(UserInfo);

function switComponents(componet) {
  whichComponent.value = componet;
}
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
.componet{
    width: 85%;
    height: 100%;
    border: 1px solid white;
    float: right;
}

</style>
