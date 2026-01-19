<template>
  <div class="mail-list">
    <div class="header">
      <span class="title">{{ title }}</span>
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
        :class="{ active: selectedId === mail.id }"
        @click="$emit('select', mail)"
      >
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
        
        <div class="preview">{{ mail.content }}</div>

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

const props = defineProps(['mails', 'loading', 'selectedId', 'title'])
const emit = defineEmits(['select', 'refresh', 'search', 'delete'])

const mailList = computed(() => props.mails || [])
const searchQuery = ref('')

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
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
  width: 300px; /* Reduced width slightly */
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

.search-bar {
  padding: 8px 12px;
  border-bottom: 1px solid #f0f0f0;
}

.list-container {
  flex: 1;
  overflow-y: auto;
}

.mail-item {
  padding: 10px 16px;
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
  background: #eef5fe;
  border-color: transparent;
}

.mail-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: #409EFF;
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
  max-width: 160px;
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
  color: #409EFF;
}

.preview {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  height: 18px; /* Force single line height */
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
