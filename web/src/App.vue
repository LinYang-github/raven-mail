<template>
  <div class="app-container">
    <Sidebar :currentView="currentView" @update:view="setView" @compose="openCompose" />
    
    <div class="main-content">
      <MailList 
        :title="listTitle"
        :mails="mails"
        :loading="loading"
        :selectedId="selectedMail?.id"
        @select="selectMail"
        @refresh="fetchMails"
        @search="handleSearch"
        @delete="handleDelete"
      />
      
      <div class="right-pane">
        <template v-if="isComposing">
          <ComposeView @cancel="cancelCompose" @success="handleComposeSuccess" />
        </template>
        <template v-else>
          <MailDetail :mail="selectedMail" />
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import Sidebar from './components/Sidebar.vue'
import MailList from './components/MailList.vue'
import MailDetail from './components/MailDetail.vue'
import ComposeView from './components/ComposeView.vue'
import { getInbox, getSent, getMail, deleteMail } from './services/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const currentView = ref('inbox')
const mails = ref([])
const selectedMail = ref(null)
const isComposing = ref(false)
const loading = ref(false)
const searchQuery = ref('')

const listTitle = computed(() => currentView.value === 'inbox' ? '收件箱' : '已发送')

const setView = (view) => {
  currentView.value = view
  searchQuery.value = '' // Reset search
  isComposing.value = false
  selectedMail.value = null
}

const openCompose = () => {
  isComposing.value = true
  selectedMail.value = null // Deselect mail logic
}

const cancelCompose = () => {
  isComposing.value = false
  // Optionally select the first mail if available
}

const fetchMails = async () => {
  loading.value = true
  try {
    const res = currentView.value === 'inbox' 
      ? await getInbox(1, searchQuery.value) 
      : await getSent(1, searchQuery.value)
    mails.value = res.data.data || []
    
    // Auto-select logic if needed, but current logic is fine
    if (!isComposing.value && selectedMail.value && !mails.value.find(m => m.id === selectedMail.value.id)) {
      selectedMail.value = null
    }
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const handleSearch = (q) => {
  searchQuery.value = q
  fetchMails()
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这封文电吗？', '提示', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消'
    })
    
    await deleteMail(id)
    ElMessage.success('删除成功')
    
    // Refresh list
    fetchMails()
    
    if (selectedMail.value?.id === id) {
      selectedMail.value = null
    }
  } catch (err) {
    if (err !== 'cancel') {
      console.error(err)
      ElMessage.error('删除失败')
    }
  }
}

const selectMail = async (mail) => {
  isComposing.value = false
  selectedMail.value = mail // Optimistic UI
  // Mark as read or fetch full details if needed
  try {
    const res = await getMail(mail.id)
    selectedMail.value = res.data
  } catch (err) {
    console.error(err)
  }
}

const handleComposeSuccess = () => {
  isComposing.value = false
  if (currentView.value === 'sent') {
    fetchMails()
  } else {
    // Maybe switch to sent view? Or just stay in inbox
    ElMessage.success('已发送')
  }
}

watch(currentView, () => {
  fetchMails()
})

onMounted(() => {
  fetchMails()
})
</script>

<style>
/* Global reset */
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

.app-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  background: #f5f7fa;
  overflow: hidden;
}

.main-content {
  flex: 1;
  display: flex;
  background: white;
  overflow: hidden;
}

.right-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}
</style>
