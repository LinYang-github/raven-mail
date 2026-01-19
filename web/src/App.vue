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
      />
      
      <MailDetail :mail="selectedMail" />
    </div>

    <ComposeModal v-model="showCompose" @success="handleComposeSuccess" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import Sidebar from './components/Sidebar.vue'
import MailList from './components/MailList.vue'
import MailDetail from './components/MailDetail.vue'
import ComposeModal from './components/ComposeModal.vue'
import { getInbox, getSent, getMail } from './services/api'

const currentView = ref('inbox')
const mails = ref([])
const selectedMail = ref(null)
const loading = ref(false)
const showCompose = ref(false)

const listTitle = computed(() => currentView.value === 'inbox' ? '收件箱' : '已发送')

const setView = (view) => {
  currentView.value = view
}

const openCompose = () => {
  showCompose.value = true
}

const fetchMails = async () => {
  loading.value = true
  try {
    const res = currentView.value === 'inbox' ? await getInbox() : await getSent()
    mails.value = res.data.data || []
    selectedMail.value = null
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const selectMail = async (mail) => {
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
  if (currentView.value === 'sent') {
    fetchMails()
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
</style>
