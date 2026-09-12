<script setup lang="ts">
defineProps<{ title: string }>()
</script>

<template>
  <div class="the-header">
    <span class="title">{{ title }}</span>
    <div class="actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
.the-header {
  height: 56px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}
</style>
