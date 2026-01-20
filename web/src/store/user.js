import { reactive } from 'vue'

export const userStore = reactive({
  id: '',
  name: '',
  sessionId: 'default',
  unreadCount: 0,
  eventSource: null,
  fetchUsers: null, // Global user directory helper
  qiankunActions: null, // Qiankun state actions
  config: {
    showReset: true,      // 是否显示场次重置按钮
    showSidebar: true,    // 是否显示侧边栏
    primaryColor: '#409EFF' // 主题色
  },

  setUser(id, name) {
    this.id = id
    this.name = name
  },

  setSession(sid) {
    this.sessionId = sid
  },

  setFetchUsers(fn) {
    this.fetchUsers = fn
  },

  setQiankunActions(actions) {
    this.qiankunActions = actions
  },

  applyConfig(newConfig) {
    if (!newConfig) return
    this.config = { ...this.config, ...newConfig }
    
    // Apply primary color to CSS variables
    if (this.config.primaryColor) {
      document.documentElement.style.setProperty('--raven-primary-color', this.config.primaryColor);
      // Generate light/dark versions (simplified for demo)
      document.documentElement.style.setProperty('--raven-primary-light', `${this.config.primaryColor}1a`); // 10% opacity
      document.documentElement.style.setProperty('--raven-primary-dark', this.config.primaryColor); 
    }
  },

  setUnreadCount(count) {
    this.unreadCount = count
    this.notifyHost()
  },

  // SSE Notifications
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
    console.log(`[raven-mail] Notifying host, unreadCount: ${this.unreadCount}`)
    
    // 1. Standalone / Legacy way: CustomEvent
    const event = new CustomEvent('raven-new-mail', { 
        detail: { unreadCount: this.unreadCount, userId: this.id } 
    });
    window.dispatchEvent(event);

    // 2. Standard way: Qiankun state
    if (this.qiankunActions && this.qiankunActions.setGlobalState) {
        this.qiankunActions.setGlobalState({
            unreadCount: this.unreadCount,
            lastUser: this.id
        });
    }
  }
})
