<template>
  <div class="mail-list">
    <div class="header">
      <span class="title">{{ title }}</span>
      <el-button circle size="small" @click="$emit('refresh')">
        <el-icon><Refresh /></el-icon>
      </el-button>
    </div>

    <div class="list-container" v-loading="loading">
      <el-empty v-if="mailList.length === 0 && !loading" description="暂无消息" />
      
      <div 
        v-for="mail in mailList" 
        :key="mail.id"
        class="mail-item"
        :class="{ active: selectedId === mail.id }"
        @click="$emit('select', mail)"
      >
        <div class="item-header">
          <span class="sender">{{ mail.sender_id }}</span>
          <span class="time">{{ formatDate(mail.created_at) }}</span>
        </div>
        <div class="subject" :class="{ 'active-text': selectedId === mail.id }">
          {{ mail.subject || '(无主题)' }}
        </div>
        <div class="preview">{{ mail.content }}</div>
        <div v-if="mail.attachments?.length" class="attachment-tag">
          <el-icon><Paperclip /></el-icon> {{ mail.attachments.length }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Refresh, Paperclip } from '@element-plus/icons-vue'

const props = defineProps(['mails', 'loading', 'selectedId', 'title'])
defineEmits(['select', 'refresh'])

const mailList = computed(() => props.mails || [])

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleDateString()
}
</script>

<style scoped>
.mail-list {
  width: 320px;
  background: white;
  border-right: 1px solid #ebedf0;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.header {
  padding: 20px 24px;
  border-bottom: 1px solid #ebedf0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header .title {
  font-size: 16px;
  font-weight: 700;
  color: #1a1a1a;
}

.list-container {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.mail-item {
  padding: 16px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.mail-item:hover {
  background: #f8f9fb;
}

.mail-item.active {
  background: #f0f7ff;
  border-color: #c6e2ff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 6px;
}

.sender {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
}

.time {
  font-size: 11px;
  color: #909399;
}

.subject {
  font-size: 13px;
  font-weight: 500;
  color: #606266;
  margin-bottom: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.subject.active-text {
  color: #409EFF;
}

.preview {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.attachment-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #f0f2f5;
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 11px;
  color: #909399;
  margin-top: 8px;
}
</style>
