<template>
  <div class="rich-content" v-html="content"></div>
</template>

<script setup>
defineProps(['content'])
</script>

<style scoped>
.rich-content {
  font-size: 15px;
  line-height: 1.8;
  color: #303133;
}
</style>
