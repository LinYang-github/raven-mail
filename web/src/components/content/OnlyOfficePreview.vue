<template>
  <div class="onlyoffice-preview">
    <div class="placeholder">
      <el-icon size="40"><Document /></el-icon>
      <div class="text">文电正文 (ONLYOFFICE 文档)</div>
      <el-button type="primary" plain @click="handleOpen">在线阅读</el-button>
    </div>
  </div>
</template>

<script setup>
import { Document } from '@element-plus/icons-vue'

const props = defineProps(['content'])

const handleOpen = () => {
  // In a real project, this would open a dialog or a new tab with the OnlyOffice viewer
  window.alert('正在从 ONLYOFFICE 服务器加载文档视图...')
}
</script>

<style scoped>
.onlyoffice-preview {
  padding: 40px;
  background: #f8f9fb;
  border-radius: 8px;
  border: 1px solid #ebedf0;
  text-align: center;
}
.placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  color: #606266;
}
.text {
  font-weight: 600;
}
</style>
