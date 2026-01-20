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
  chats: {}, // { userId: [messages] }
  chatUnreads: {}, // { userId: count }
  totalIMUnread: 0,

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

  markIMRead(userId) {
    this.chatUnreads[userId] = 0
    this.recalculateIMUnread()
  },

  recalculateIMUnread() {
    this.totalIMUnread = Object.values(this.chatUnreads).reduce((sum, count) => sum + count, 0)
  },

  applySummary(summary) {
    if (!summary) return
    if (summary.unread_mail_count !== undefined) {
      this.unreadCount = summary.unread_mail_count
    }
    if (summary.im_unread_counts) {
      this.chatUnreads = summary.im_unread_counts
      this.recalculateIMUnread()
    }
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
      try {
        const payload = JSON.parse(event.data)
        console.log('[raven-mail] Notification received:', payload)

        // 场次隔离校验
        if (payload.session_id !== this.sessionId) return

        // 是否发送给我的
        if (payload.targets && !payload.targets.includes(this.id)) return

        if (payload.type === 'MAIL') {
          this.unreadCount++
          this.notifyHost()
          // 触发邮件更新提示
          window.dispatchEvent(new CustomEvent('raven-mail-updated', { detail: payload.data }))
        } else if (payload.type === 'CHAT') {
          const msg = payload.data
          const chatPartner = msg.sender_id === this.id ? msg.receiver_id : msg.sender_id
          if (!this.chats[chatPartner]) this.chats[chatPartner] = []
          this.chats[chatPartner].push(msg)
          
          // 如果消息是别人发给我的，且不在当前对话中（由 UI 层决定是否标记），累加未读
          if (msg.sender_id !== this.id) {
            this.chatUnreads[msg.sender_id] = (this.chatUnreads[msg.sender_id] || 0) + 1
            this.recalculateIMUnread()
          }

          // 触发 IM 事件
          window.dispatchEvent(new CustomEvent('raven-im-received', { detail: msg }))
        }
      } catch (e) {
        // 兼容旧版或心跳及原始字符串
        const msg = event.data
        if (msg.startsWith('NEW_MAIL:')) {
            const parts = msg.split(':')
            const sessionId = parts[1]
            const targetIds = parts[2].split(',')
            if (sessionId === this.sessionId && targetIds.includes(this.id)) {
              this.unreadCount++
              this.notifyHost()
            }
        }
      }
    }

    this.eventSource.onerror = (err) => {
      console.error('[raven-mail] SSE error:', err)
      if (this.eventSource) {
        this.eventSource.close()
        this.eventSource = null
      }
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
