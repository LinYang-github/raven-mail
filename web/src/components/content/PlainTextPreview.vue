<template>
  <pre class="body-text">{{ content }}</pre>
</template>

<script setup>
defineProps(['content'])
</script>

<style scoped>
.body-text {
  white-space: pre-wrap;
  font-family: inherit;
  font-size: 15px;
  color: #303133;
  line-height: 1.8;
  margin: 0;
}
</style>
