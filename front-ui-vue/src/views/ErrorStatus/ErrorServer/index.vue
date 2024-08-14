<template>
    <a-result status="500" title="500" sub-title="对不起，服务器错误。">
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