import { reactive } from 'vue'

export const userStore = reactive({
  id: '',
  name: '',
  sessionId: localStorage.getItem('raven_session_id') || 'default',
  
  setUser(id, name) {
    this.id = id
    this.name = name || id
    console.log('[raven-mail] User switched to:', id)
  },

  setSession(sid) {
    this.sessionId = sid || 'default'
    localStorage.setItem('raven_session_id', this.sessionId)
    console.log('[raven-mail] Session switched to:', this.sessionId)
    // Refetching might be needed, handled by components watching this
  },

  // 远程搜索函数引用
  fetchUsers: null,
  
  setFetchUsers(fn) {
    this.fetchUsers = fn
  },

  // SSE Notifications
  eventSource: null,
  unreadCount: 0,
  
  initNotifications() {
    if (this.eventSource) return;
    
    console.log('[raven-mail] Initializing notifications...')
    const isDev = import.meta.env.DEV
    const backendUrl = import.meta.env.VITE_BACKEND_URL || (isDev ? 'http://localhost:8080' : '')
    this.eventSource = new EventSource(`${backendUrl}/api/v1/mails/events`)
    
    this.eventSource.onmessage = (event) => {
      const msg = event.data
      console.log('[raven-mail] Notification received:', msg)
      
      if (msg.startsWith('NEW_MAIL:')) {
        const parts = msg.split(':')
        const sessionId = parts[1]
        const targetIds = parts[2].split(',')
        
        // 仅当场次匹配且用户在目标列表中时才提醒
        if (sessionId === this.sessionId && targetIds.includes(this.id)) {
          this.unreadCount++
          this.notifyHost()
        }
      }
    }

    this.eventSource.onerror = (err) => {
      console.error('[raven-mail] SSE error:', err)
      this.eventSource.close()
      this.eventSource = null
      // Retry after 5s
      setTimeout(() => this.initNotifications(), 5000)
    }
  },

  notifyHost() {
    // Notify host via custom event or qiankun state
    const event = new CustomEvent('raven-new-mail', { 
        detail: { unreadCount: this.unreadCount, userId: this.id } 
    });
    window.dispatchEvent(event);
  }
})
