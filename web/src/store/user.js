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
  }
})
