<template>
  <div class="mail-client">
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
            :replyTo="replyTarget"
            @cancel="cancelCompose"
            @success="handleComposeSuccess"
            @reply="handleReply"
          />
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Sidebar from './Sidebar.vue'
import MailList from './MailList.vue'
import MailDetail from './MailDetail.vue'
import ComposeView from './ComposeView.vue'
import { getInbox, getSent, getMail, deleteMail, getUserSummary } from '../services/api'
import { userStore } from '../store/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()

const mails = ref([])
const selectedMail = ref(null)
const loading = ref(false)
const searchQuery = ref('')
const replyTarget = ref(null)

// 从路由推导当前视图
const currentView = computed(() => {
  if (route.path.includes('/inbox')) return 'inbox'
  if (route.path.includes('/sent')) return 'sent'
  if (route.path.includes('/compose')) return 'compose'
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

const fetchSummary = async () => {
  if (!userStore.id) return
  try {
    const res = await getUserSummary()
    userStore.applySummary(res.data)
  } catch (err) {
    console.error('Failed to fetch user summary:', err)
  }
}

// 监听视图切换，重新抓取列表
watch([currentView, () => userStore.id, () => userStore.sessionId], () => {
  if (currentView.value !== 'compose') {
    fetchMails()
  }
  fetchSummary()
})

const setView = (view) => {
  router.push(`/${view}`)
}

const openCompose = () => {
  replyTarget.value = null
  router.push('/compose')
}

const cancelCompose = () => {
  replyTarget.value = null
  router.back()
}

const handleReply = (mode) => {
  if (!selectedMail.value) return
  replyTarget.value = {
    mode: mode,
    mail: selectedMail.value
  }
  router.push('/compose')
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
  router.push(`/${currentView.value}/${mail.id}`)
  
  if (currentView.value === 'inbox') {
    const r = mail.recipients?.find(rp => rp.recipient_id === userStore.id)
    if (r && r.status === 'unread') {
      userStore.setUnreadCount(Math.max(0, userStore.unreadCount - 1))
      r.status = 'read'
    }
  }
}

const handleComposeSuccess = () => {
  replyTarget.value = null
  ElMessage.success('发送成功')
  router.push('/sent')
}

onMounted(() => {
  fetchMails()
  fetchSummary()
})
</script>

<style scoped>
.mail-client {
  display: flex;
  height: 100%;
  width: 100%;
  overflow: hidden;
  background: #f5f7fa;
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
