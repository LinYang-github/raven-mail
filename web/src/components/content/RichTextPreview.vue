<template>
  <div class="rich-content-wrapper editor-content-view" v-html="content"></div>
</template>

<script setup>
import '@wangeditor/editor/dist/css/style.css'
defineProps(['content'])
</script>

<style scoped>
.rich-content-wrapper {
  font-size: 15px;
  line-height: 1.6;
  color: #303133;
  padding: 10px 0;
}

/* 兼容 wangEditor 展示样式 */
:deep(ul) {
  list-style-type: disc;
  margin-left: 20px;
}
:deep(ol) {
  list-style-type: decimal;
  margin-left: 20px;
}
:deep(blockquote) {
  border-left: 4px solid #ccc;
  padding-left: 10px;
  color: #666;
  margin: 10px 0;
}
:deep(table) {
  border-collapse: collapse;
  width: 100%;
}
:deep(table td), :deep(table th) {
  border: 1px solid #ccc;
  padding: 5px;
}
:deep(pre) {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
}
</style>
