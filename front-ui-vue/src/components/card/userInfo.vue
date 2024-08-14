<template>
    <div v-if="flag">
        <a-card hoverable style="width: 100%; margin: 0 auto">
          <template #cover>
            <img alt="example" :src="cover" />
          </template>
          <template #actions>
          </template>
          <div>
            <div class="div-image">
              <a-avatar
                class="image"
                :size="100"
                :src="userInfo?.avatar"
              />
              <h4>{{ userInfo.name }}</h4>
            </div>
            <div>
              <div>
                <span class="left">ID </span
                ><span class="right">{{ userInfo.identity }}</span>
              </div>
              <div>
                <span class="left">备注 </span
                ><span class="right">{{ userInfo.identity }}</span>
              </div>
              <div>
                <span class="left">邮箱</span
                ><span class="right">{{ userInfo.email }}</span>
              </div>
              <div>
                <span class="left">手机号</span
                ><span class="right">{{ userInfo.phone }}</span>
              </div>
            </div>
          </div>
        </a-card>
    </div>
</template>
<script setup lang="ts">
import cover from "@/assets/images/login.jpg";
import {ref} from 'vue'
const props = defineProps({
    flag: Boolean,
    id: Number,
})
const userInfo = ref({})

</script>
<style lang="less">
</style>