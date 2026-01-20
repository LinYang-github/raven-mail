<template>
  <div class="app-container">
    <Sidebar v-if="userStore.config.showSidebar" :currentView="currentView" @update:view="setView" @compose="openCompose" />
    
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
        <router-view v-slot="{ Component }">
          <component 
            :is="Component" 
            :mail="selectedMail"
            @cancel="cancelCompose"
            @success="handleComposeSuccess"
          />
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import MailList from './components/MailList.vue'
import MailDetail from './components/MailDetail.vue'
import ComposeView from './components/ComposeView.vue'
import { getInbox, getSent, getMail, deleteMail } from './services/api'
import { userStore } from './store/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()

const mails = ref([])
const selectedMail = ref(null)
const loading = ref(false)
const searchQuery = ref('')

// 从路由推导当前视图
const currentView = computed(() => {
  if (route.path.startsWith('/inbox')) return 'inbox'
  if (route.path.startsWith('/sent')) return 'sent'
  if (route.path.startsWith('/compose')) return 'compose'
  return 'inbox'
})

const isComposing = computed(() => currentView.value === 'compose')
const listTitle = computed(() => currentView.value === 'inbox' ? '收件箱' : '已发送')

// 监听路由参数中的 ID 变化，加载邮件详情
watch(() => route.params.id, async (newId) => {
  if (newId) {
    try {
      const res = await getMail(newId)
      selectedMail.value = res.data
    } catch (err) {
      console.error('Failed to fetch mail detail:', err)
      selectedMail.value = null
    }
  } else {
    selectedMail.value = null
  }
}, { immediate: true })

// 监听视图切换，重新抓取列表
watch([currentView, () => userStore.id, () => userStore.sessionId], () => {
  if (currentView.value !== 'compose') {
    fetchMails()
  }
})

const setView = (view) => {
  router.push(`/${view}`)
}

const openCompose = () => {
  router.push('/compose')
}

const cancelCompose = () => {
  router.back()
}

const fetchMails = async () => {
  loading.value = true
  try {
    const res = currentView.value === 'sent' 
      ? await getSent(1, searchQuery.value)
      : await getInbox(1, searchQuery.value)
    
    mails.value = res.data.data || []
    
    // 如果是收件箱，统计未读数并回传主应用
    if (currentView.value === 'inbox') {
      const unread = mails.value.filter(m => {
        const r = m.recipients?.find(rp => rp.recipient_id === userStore.id)
        return r && r.status === 'unread'
      }).length
      userStore.setUnreadCount(unread)
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
    
    fetchMails()
    
    if (route.params.id === id) {
      router.push(`/${currentView.value}`)
    }
  } catch (err) {
    if (err !== 'cancel') {
      console.error(err)
      ElMessage.error('删除失败')
    }
  }
}

const selectMail = (mail) => {
  // 仅跳转路由，详情加载由上面的 watch 处理
  router.push(`/${currentView.value}/${mail.id}`)
  
  // 体验优化：如果是未读，前端先行扣减
  if (currentView.value === 'inbox') {
    const r = mail.recipients?.find(rp => rp.recipient_id === userStore.id)
    if (r && r.status === 'unread') {
      userStore.setUnreadCount(Math.max(0, userStore.unreadCount - 1))
      r.status = 'read'
    }
  }
}

const handleComposeSuccess = () => {
  ElMessage.success('发送成功')
  router.push('/sent')
}

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
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
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
