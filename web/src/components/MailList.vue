<template>
  <div class="mail-list">
    <div class="header">
      <div class="header-left">
        <span class="title">{{ title }}</span>
      </div>
      <div class="actions">
        <el-button circle size="small" @click="$emit('refresh')">
          <el-icon><Refresh /></el-icon>
        </el-button>
      </div>
    </div>
    
    <div class="search-bar">
      <el-input 
        v-model="searchQuery" 
        placeholder="搜索邮件..." 
        prefix-icon="Search"
        clearable
        @input="handleSearch"
      />
    </div>

    <div class="list-container" v-loading="loading">
      <el-empty v-if="mailList.length === 0 && !loading" description="暂无消息" />
      
      <div 
        v-for="mail in mailList" 
        :key="mail.id"
        class="mail-item"
        :class="{ active: selectedId === mail.id, unread: isUnread(mail) }"
        @click="$emit('select', mail)"
      >
        <div v-if="isUnread(mail)" class="unread-dot"></div>
        <div class="item-header">
          <span class="sender">{{ mail.sender_id }}</span>
          <div class="header-right">
            <div v-if="mail.attachments?.length" class="attachment-icon">
              <el-icon><Paperclip /></el-icon>
            </div>
            <span class="time">{{ formatDate(mail.created_at) }}</span>
          </div>
        </div>
        
        <div class="subject-row">
          <div class="subject" :class="{ 'active-text': selectedId === mail.id }">
            {{ mail.subject || '(无主题)' }}
          </div>
        </div>
        
        <div class="preview">
          <span v-if="mail.content_type === 'onlyoffice'" class="office-label">
            <el-icon style="vertical-align: middle; margin-right: 4px;"><Document /></el-icon>
            在线正文
          </span>
          <span v-else>{{ mail.content || '(无内容)' }}</span>
        </div>

        <el-button 
          type="danger" 
          link 
          size="small" 
          class="delete-btn"
          @click.stop="$emit('delete', mail.id)"
        >
          <el-icon><Delete /></el-icon>
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { Refresh, Paperclip, Search, Delete } from '@element-plus/icons-vue'
import { debounce } from 'lodash'
import { userStore } from '../store/user'

const props = defineProps(['mails', 'loading', 'selectedId', 'title'])
const emit = defineEmits(['select', 'refresh', 'search', 'delete'])

const isUnread = (mail) => {
  const r = mail.recipients?.find(rp => rp.recipient_id === userStore.id)
  return r && r.status === 'unread'
}

const mailList = computed(() => props.mails || [])
const searchQuery = ref('')

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  const pad = (n) => n.toString().padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

let timeout
const handleSearch = () => {
  clearTimeout(timeout)
  timeout = setTimeout(() => {
    emit('search', searchQuery.value)
  }, 500)
}

</script>

<style scoped>
.mail-list {
  width: 320px; /* Increased to fit date */
  background: white;
  border-right: 1px solid #ebedf0;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.header {
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header .title {
  font-size: 15px;
  font-weight: 700;
  color: #1a1a1a;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-bar {
  padding: 8px 12px;
  border-bottom: 1px solid #f0f0f0;
}

.list-container {
  flex: 1;
  overflow-y: auto;
}

.mail-item {
  padding: 12px 18px; /* Increased padding */
  border-bottom: 1px solid #f5f7fa;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.mail-item:hover {
  background: #f8f9fb;
}

.mail-item:hover .delete-btn {
  display: flex;
}

.mail-item.active {
  background: var(--raven-primary-light);
  border-color: transparent;
}

.mail-item.unread .sender {
  color: #1a1a1a;
  font-weight: 700;
}

.mail-item.unread .subject {
  color: #303133;
  font-weight: 600;
}

.unread-dot {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 8px;
  height: 8px;
  background-color: #f56c6c;
  border-radius: 50%;
  box-shadow: 0 0 4px rgba(245, 108, 108, 0.4);
}

.mail-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--raven-primary-color);
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 6px;
}

.sender {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100px; /* Reduced to fit date */
}

.time {
  font-size: 12px;
  color: #909399;
}

.attachment-icon {
  color: #909399;
  font-size: 12px;
  display: flex;
}

.subject-row {
  margin-bottom: 2px;
}

.subject {
  font-size: 13px;
  font-weight: 500;
  color: #606266;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.subject.active-text {
  color: var(--raven-primary-color);
}

.preview {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  height: 18px; /* Force single line height */
}

.office-label {
  color: var(--raven-primary-color);
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  background: var(--raven-primary-light);
  padding: 0 6px;
  border-radius: 4px;
}

.delete-btn {
  display: none; /* Hidden by default */
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: white; /* Cover text behind it */
  box-shadow: -10px 0 10px white; /* Fade out effect */
  border-radius: 4px;
}
</style>
