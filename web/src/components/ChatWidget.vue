<template>
  <div class="chat-widget-container">
    <!-- 浮动气泡 -->
    <div class="chat-bubble" @click="toggleChat" :class="{ 'has-unread': totalUnread > 0 }">
      <el-badge :value="totalUnread" :hidden="totalUnread === 0" class="chat-badge">
        <el-icon :size="24"><ChatDotRound /></el-icon>
      </el-badge>
    </div>

    <!-- 聊天窗口 -->
    <transition name="slide-up">
      <div v-if="isOpen" class="chat-window">
        <div class="chat-header">
          <div class="header-info">
            <span class="title">即时通讯</span>
            <span v-if="activePartner" class="partner-name"> - {{ activePartnerName }}</span>
          </div>
          <div class="header-actions">
            <el-button link :icon="Close" @click="isOpen = false"></el-button>
          </div>
        </div>

        <div class="chat-body">
          <!-- 用户列表 (左侧/折叠) -->
          <div v-if="!activePartner" class="user-selector">
            <div class="search-box">
              <el-input 
                v-model="userQuery" 
                placeholder="搜索联系人..." 
                prefix-icon="Search"
                size="small"
                @input="handleSearch"
              />
            </div>
            <div class="user-list">
              <div 
                v-for="user in userOptions" 
                :key="user.id" 
                class="user-item"
                @click="selectPartner(user)"
              >
                <el-badge :value="userStore.chatUnreads[user.id]" :hidden="!userStore.chatUnreads[user.id]" class="user-badge">
                  <el-avatar :size="32">{{ user.name.charAt(0) }}</el-avatar>
                </el-badge>
                <div class="user-info">
                  <div class="name">{{ user.name }}</div>
                  <div class="dept">{{ user.dept }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 对话区域 (展示时) -->
          <div v-else class="conversation">
            <div class="conversation-back" @click="activePartner = null">
              <el-icon><ArrowLeft /></el-icon> 返回联系人
            </div>
            <div class="message-list" ref="messageListRef">
              <div 
                v-for="msg in currentMessages" 
                :key="msg.id" 
                class="message-row"
                :class="{ 'is-me': msg.sender_id === userStore.id }"
              >
                <div class="message-content">
                  <div class="text">{{ msg.content }}</div>
                  
                  <!-- Attachments -->
                  <div v-if="msg.attachments && msg.attachments.length > 0" class="chat-attachments">
                    <div v-for="att in msg.attachments" :key="att.id" class="chat-attachment-item">
                        <template v-if="isImage(att)">
                            <el-image 
                                :src="getAttachmentUrl(att)" 
                                :preview-src-list="[getAttachmentUrl(att)]"
                                class="chat-image"
                                fit="cover"
                            />
                        </template>
                        <template v-else>
                            <a :href="getAttachmentUrl(att)" target="_blank" class="chat-file">
                                <el-icon><Document /></el-icon>
                                <span>{{ att.file_name }}</span>
                            </a>
                        </template>
                    </div>
                  </div>

                  <div class="time">{{ formatTime(msg.created_at) }}</div>
                </div>
              </div>
              <div v-if="currentMessages.length === 0" class="empty-chat">
                暂无聊天记录，开始打个招呼吧
              </div>
            </div>
            
            <div class="input-area">
              <!-- Pending Files -->
              <div v-if="pendingFiles.length > 0" class="pending-files">
                <div v-for="(file, index) in pendingFiles" :key="index" class="pending-file">
                    <el-icon><Document /></el-icon>
                    <span class="filename">{{ file.name }}</span>
                    <el-icon class="remove-btn" @click="removePendingFile(index)"><Close /></el-icon>
                </div>
              </div>

              <el-input
                v-model="inputText"
                type="textarea"
                :rows="2"
                placeholder="输入消息..."
                resize="none"
                @keydown.enter.prevent="handleSend"
              ></el-input>
              <div class="input-actions">
                <input type="file" ref="fileInput" multiple style="display: none" @change="handleFileSelect">
                <el-button link :icon="Paperclip" @click="triggerFileUpload" class="action-btn"></el-button>
                <el-button type="primary" size="small" :disabled="!canSend" @click="handleSend">
                  发送
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { ChatDotRound, Close, Search, ArrowLeft, Paperclip, Document } from '@element-plus/icons-vue'
import { userStore } from '../store/user'
import { sendChatMessage, getChatHistory, markChatAsRead, getPreviewUrl } from '../services/api'
import { ElMessage, ElNotification } from 'element-plus'

const isOpen = ref(false)
const activePartner = ref(null)
const userQuery = ref('')
const userOptions = ref([])
const inputText = ref('')
const messageListRef = ref(null)
const fileInput = ref(null)
const pendingFiles = ref([])

const activePartnerName = computed(() => {
  return userOptions.value.find(u => u.id === activePartner.value)?.name || activePartner.value
})

const currentMessages = computed(() => {
  return userStore.chats[activePartner.value] || []
})

const totalUnread = computed(() => {
  return userStore.totalIMUnread
})

const canSend = computed(() => {
    return (inputText.value.trim() || pendingFiles.value.length > 0) && activePartner.value
})

const toggleChat = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    handleSearch('')
  }
}

const handleSearch = async (query) => {
  if (userStore.fetchUsers) {
    const results = await userStore.fetchUsers(query)
    // 过滤掉自己
    userOptions.value = results.filter(u => u.id !== userStore.id)
  }
}

const selectPartner = async (user) => {
  activePartner.value = user.id
  // 清除本地未读并同步到后端
  if (userStore.chatUnreads[user.id]) {
    userStore.markIMRead(user.id)
    try {
      await markChatAsRead(user.id)
    } catch (e) {
      console.warn('Failed to mark IM as read', e)
    }
  }

  // 加载历史记录
  try {
    const res = await getChatHistory(user.id)
    userStore.chats[user.id] = res.data
    scrollToBottom()
  } catch (err) {
    console.error('Failed to load chat history', err)
  }
}

const triggerFileUpload = () => {
    fileInput.value.click()
}

const handleFileSelect = (event) => {
    const files = Array.from(event.target.files)
    pendingFiles.value.push(...files)
    // Reset input so same file can be selected again
    event.target.value = ''
}

const removePendingFile = (index) => {
    pendingFiles.value.splice(index, 1)
}

const handleSend = async () => {
  if (!canSend.value) return
  
  const content = inputText.value.trim()
  const files = [...pendingFiles.value]
  
  inputText.value = ''
  pendingFiles.value = []

  try {
    const res = await sendChatMessage(activePartner.value, content, files)
    // 推送到本地状态
    if (!userStore.chats[activePartner.value]) {
      userStore.chats[activePartner.value] = []
    }
    userStore.chats[activePartner.value].push(res.data)
    scrollToBottom()
  } catch (err) {
    ElMessage.error('发送失败')
    // Restore text if failed (optional, but good UX)
    if (content) inputText.value = content
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  })
}

const formatTime = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const isImage = (att) => {
    return att.mime_type && att.mime_type.startsWith('image/')
}

const getAttachmentUrl = (att) => {
    return getPreviewUrl(att)
}

// 监听新消息事件
onMounted(() => {
  window.addEventListener('raven-im-received', async (e) => {
    const msg = e.detail
    // 如果窗口是开着的，且对话框正是该发信人，自动标记为已读
    if (isOpen.value && activePartner.value === msg.sender_id) {
        userStore.markIMRead(msg.sender_id)
        try {
            await markChatAsRead(msg.sender_id)
        } catch (e) {}
        scrollToBottom()
    } else if (msg.sender_id === activePartner.value || msg.receiver_id === activePartner.value) {
      scrollToBottom()
    }
    
    // 如果是别人发来的消息，且窗体未打开或当前不是与该人的对话，显示系统通知
    if (msg.sender_id !== userStore.id && (!isOpen.value || activePartner.value !== msg.sender_id)) {
        ElNotification({
            title: '新即时消息',
            message: msg.content || '[附件]',
            type: 'info',
            position: 'bottom-right',
            offset: 80, // 避免挡住气泡
            duration: 3000
        })
    }
  })
})

watch(activePartner, () => {
    scrollToBottom()
    pendingFiles.value = [] // clear pending files on switch
})
</script>

<style scoped>
.chat-widget-container {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 2000;
}

.chat-bubble {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--raven-primary-color);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  transition: all 0.3s;
}

.chat-bubble:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(0,0,0,0.2);
}

.chat-badge :deep(.el-badge__content) {
  box-shadow: 0 0 0 2px white;
}

.user-badge :deep(.el-badge__content) {
  right: 5px;
  top: 5px;
}

.chat-window {
  position: absolute;
  right: 0;
  bottom: 70px;
  width: 350px;
  height: 500px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.15);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid #ebeef5;
}

.chat-header {
  padding: 12px 16px;
  background: var(--raven-primary-color);
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-info .title {
  font-weight: 600;
}

.header-info .partner-name {
  font-size: 13px;
  opacity: 0.9;
}

.chat-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #f5f7fa;
}

.user-selector {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.search-box {
  padding: 12px;
  background: white;
  border-bottom: 1px solid #ebeef5;
}

.user-list {
  flex: 1;
  overflow-y: auto;
}

.user-item {
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  background: white;
  cursor: pointer;
  transition: background 0.2s;
}

.user-item:hover {
  background: #f0f2f5;
}

.user-item .user-info .name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.user-item .user-info .dept {
  font-size: 12px;
  color: #909399;
}

.conversation {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.conversation-back {
  padding: 8px 12px;
  background: white;
  font-size: 13px;
  color: var(--raven-primary-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  border-bottom: 1px solid #ebeef5;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-row {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.message-row.is-me {
  align-items: flex-end;
}

.message-content {
  max-width: 80%;
  padding: 8px 12px;
  border-radius: 8px;
  background: white;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  position: relative;
}

.is-me .message-content {
  background: #e1f3ff;
}

.message-content .text {
  font-size: 14px;
  line-height: 1.5;
  color: #303133;
  word-break: break-word;
}

.message-content .time {
  font-size: 10px;
  color: #999;
  margin-top: 4px;
  text-align: right;
}

.empty-chat {
  text-align: center;
  color: #999;
  font-size: 12px;
  margin-top: 40px;
}

.input-area {
  padding: 12px;
  background: white;
  border-top: 1px solid #ebeef5;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.action-btn {
    font-size: 18px;
    color: #606266;
}
.action-btn:hover {
    color: var(--raven-primary-color);
}

.pending-files {
    padding: 0 0 8px 0;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
}

.pending-file {
    font-size: 12px;
    background: #f0f2f5;
    padding: 4px 8px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    gap: 4px;
}

.pending-file .filename {
    max-width: 150px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.pending-file .remove-btn {
    cursor: pointer;
    color: #909399;
}
.pending-file .remove-btn:hover {
    color: #f56c6c;
}

.chat-attachments {
    margin-top: 8px;
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.chat-image {
    max-width: 100%;
    border-radius: 4px;
    cursor: pointer;
}

.chat-file {
    display: flex;
    align-items: center;
    gap: 6px;
    text-decoration: none;
    color: #409EFF;
    font-size: 13px;
    background: #f4f4f5;
    padding: 4px 8px;
    border-radius: 4px;
}
.chat-file:hover {
    background: #ecf5ff;
}

/* Animations */
.slide-up-enter-active, .slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from, .slide-up-leave-to {
  transform: translateY(20px);
  opacity: 0;
}
</style>
