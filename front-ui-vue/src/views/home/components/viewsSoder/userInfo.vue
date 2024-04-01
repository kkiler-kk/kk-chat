<template>
  <div style="height: 100%;">
    <div class="div1">
      <span class="title">信 息</span>
      <a-button class="botton">退 出</a-button>
    </div>
    <div class="userInfo">
        <a-avatar :src="logo" :size="78" />
        <p>name: KK</p>
        <p>id:id</p>
        <p>email:kk996icu@qq.com</p>
        <p>phone: 18857392027</p>
    </div>
  </div>
</template>
<script setup lang="ts">
import logo from "@/assets/images/chat.png";
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
  margin-right: 10px;
  background-color: rgb(24, 144, 255);
  color: white;
  border-radius: 8%;
}
.userInfo {
    width: 89%;
    height: 90%;
    background-color: white;
    margin: 0 auto;
}
</style>
