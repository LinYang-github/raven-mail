<template>
  <div class="mail-detail">
    <div v-if="!mail" class="empty-state">
      <el-empty description="请选择一封邮件阅读" />
    </div>
    
    <div v-else class="content-wrapper">
      <div class="detail-header">
        <h1 class="subject">{{ mail.subject }}</h1>
        <div class="meta">
          <el-tag size="small" type="info">{{ new Date(mail.created_at).toLocaleString() }}</el-tag>
        </div>
        
        <div class="sender-info">
          <el-avatar :size="40" class="sender-avatar">{{ mail.sender_id.charAt(0).toUpperCase() }}</el-avatar>
          <div class="info-text">
            <div class="sender-name">{{ mail.sender_id }}</div>
            <div class="recipients">收件人: {{ mail.recipients?.map(r => r.recipient_id).join(', ') }}</div>
          </div>
        </div>
      </div>

      <div class="detail-body">
        <pre class="body-text">{{ mail.content }}</pre>
      </div>

      <div v-if="mail.attachments?.length" class="attachments-section">
        <div class="att-title">
          <el-icon><Paperclip /></el-icon> 附件 ({{ mail.attachments.length }})
        </div>
        <div class="att-list">
          <a v-for="att in mail.attachments" :key="att.id" :href="getDownloadUrl(att)" target="_blank" class="att-item">
            <el-icon class="att-icon"><Document /></el-icon>
            <div class="att-info">
              <div class="att-name">{{ att.file_name }}</div>
              <div class="att-size">{{ formatSize(att.file_size) }}</div>
            </div>
            <el-icon class="dl-icon"><Download /></el-icon>
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getDownloadUrl } from '../services/api'
import { Paperclip, Document, Download } from '@element-plus/icons-vue'

defineProps(['mail'])

const formatSize = (bytes) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}
</script>

<style scoped>
.mail-detail {
  flex: 1;
  background: white;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.content-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.detail-header {
  padding: 40px;
  border-bottom: 1px solid #f0f2f5;
  background: #fff;
}

.subject {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 16px 0;
  color: #1a1a1a;
  line-height: 1.4;
}

.sender-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 24px;
}

.sender-avatar {
  background: #409EFF;
  font-size: 18px;
  font-weight: 600;
}

.info-text {
  flex: 1;
}

.sender-name {
  font-weight: 600;
  color: #303133;
  font-size: 15px;
}

.recipients {
  font-size: 13px;
  color: #909399;
  margin-top: 2px;
}

.detail-body {
  flex: 1;
  padding: 40px;
  overflow-y: auto;
  background: #fff;
}

.body-text {
  white-space: pre-wrap;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
  font-size: 15px;
  color: #303133;
  line-height: 1.8;
  max-width: 800px;
}

.attachments-section {
  padding: 24px 40px;
  background: #f8f9fb;
  border-top: 1px solid #ebedf0;
}

.att-title {
  font-size: 13px;
  font-weight: 600;
  color: #606266;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.att-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border: 1px solid #ebedf0;
  border-radius: 6px;
  background: white;
  text-decoration: none;
  width: 240px;
  transition: all 0.2s;
}

.att-item:hover {
  border-color: #409EFF;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.att-icon {
  font-size: 20px;
  color: #409EFF;
  margin-right: 10px;
}

.att-info {
  flex: 1;
  overflow: hidden;
}

.att-name {
  font-size: 13px;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.att-size {
  font-size: 12px;
  color: #909399;
}

.dl-icon {
  color: #c0c4cc;
}
</style>
