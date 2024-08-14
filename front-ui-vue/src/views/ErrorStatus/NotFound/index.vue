<template>
    <a-result status="404" title="404" sub-title="对不起，您访问的页面不存在。">
        <template #extra>
            <a-button type="primary" @click="backHome">返回首页</a-button>
        </template>
    </a-result>
</template>
<script setup lang="ts">
import router from "@/router";
function backHome(){
    router.push({
        name: "Home",
        params: {}
    });
}
</script>
<style></style>