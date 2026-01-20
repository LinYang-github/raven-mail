import { reactive } from 'vue'

export const userStore = reactive({
  id: 'user-123', // Default fallback
  name: 'Default User',
  
  setUser(id, name) {
    this.id = id
    this.name = name || id
    console.log('[raven-mail] User switched to:', id)
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
    this.eventSource = new EventSource('http://localhost:8080/api/v1/mails/events')
    
    this.eventSource.onmessage = (event) => {
      const msg = event.data
      console.log('[raven-mail] Notification received:', msg)
      
      if (msg.startsWith('NEW_MAIL:')) {
        const targetIds = msg.split(':')[1].split(',')
        if (targetIds.includes(this.id)) {
          this.unreadCount++
          // Broadcast to host if needed
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
