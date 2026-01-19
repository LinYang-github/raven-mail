<template>
  <div class="sidebar">
    <div class="logo-container">
      <div class="logo-icon">
        <el-icon :size="20" color="#fff"><Message /></el-icon>
      </div>
      <span class="title">文电管理</span>
    </div>

    <div class="compose-container">
      <el-button type="primary" class="compose-btn" size="large" @click="$emit('compose')">
        <el-icon class="mr-2"><Plus /></el-icon> 新建文电
      </el-button>
    </div>

    <div class="menu-container">
      <div 
        class="menu-item" 
        :class="{ active: currentView === 'inbox' }"
        @click="$emit('update:view', 'inbox')"
      >
        <el-icon><Box /></el-icon>
        <span>收件箱</span>
      </div>
      <div 
        class="menu-item" 
        :class="{ active: currentView === 'sent' }"
        @click="$emit('update:view', 'sent')"
      >
        <el-icon><Position /></el-icon>
        <span>已发送</span>
      </div>
    </div>

    <div class="user-info">
      <el-avatar :size="32" icon="UserFilled" />
      <div class="user-details">
        <div class="name">当前用户</div>
        <div class="id">user-123</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Message, Plus, Box, Position, UserFilled } from '@element-plus/icons-vue'

defineProps(['currentView'])
defineEmits(['update:view', 'compose'])
</script>

<style scoped>
.sidebar {
  width: 260px;
  height: 100%;
  background: #f8f9fb;
  border-right: 1px solid #ebedf0;
  display: flex;
  flex-direction: column;
}

.logo-container {
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #409EFF 0%, #304bd6 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 6px rgba(64, 158, 255, 0.3);
}

.title {
  font-size: 18px;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: 0.5px;
}

.compose-container {
  padding: 0 20px 20px 20px;
}

.compose-btn {
  width: 100%;
  border-radius: 8px;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
  transition: all 0.3s;
}

.compose-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.3);
}

.mr-2 {
  margin-right: 8px;
}

.menu-container {
  flex: 1;
  padding: 0 12px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  margin-bottom: 4px;
  border-radius: 8px;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.menu-item:hover {
  background: #ecedf0;
  color: #303133;
}

.menu-item.active {
  background: #e6f0ff;
  color: #409EFF;
}

.menu-item .el-icon {
  font-size: 18px;
}

.user-info {
  padding: 20px;
  border-top: 1px solid #f0f2f5;
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-details .name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.user-details .id {
  font-size: 12px;
  color: #909399;
}
</style>
