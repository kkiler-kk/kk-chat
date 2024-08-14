<template>
    <a-result status="403" title="403" sub-title="对不起，您没有权限访问此页面。">
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