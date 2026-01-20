import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import './style.css'
import router from './router'

import { renderWithQiankun, qiankunWindow } from 'vite-plugin-qiankun/dist/helper'

import { userStore } from './store/user'

let app = null

function render(props = {}) {
  const { container, modules } = props
  app = createApp(App, { modules })
  
  app.use(ElementPlus)
  app.use(router)
  
  // Register all icons globally
  for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
  }

  // Allow mounting to a specific container if provided (qiankun context)
  const mountPoint = container ? container.querySelector('#app') : '#app'
  app.mount(mountPoint)
}

renderWithQiankun({
  bootstrap() {
    console.log('[raven-mail] bootstrap')
  },
  mount(props) {
    console.log('[raven-mail] mount', props)
    
    // Sync initial state
    if (props.user) {
      userStore.setUser(props.user.id || '', props.user.name || '')
    }
    if (props.sessionId) {
      userStore.setSession(props.sessionId)
    }
    if (props.ravenConfig) {
      userStore.applyConfig(props.ravenConfig)
    }

    // Capture qiankun actions/props
    userStore.setQiankunActions(props)
    
    // Capture fetchUsers function
    if (props.fetchUsers) {
      userStore.setFetchUsers(props.fetchUsers)
    }

    // Start Real-time Notifications
    userStore.initNotifications()

    // Sync state changes from host
    if (props.onGlobalStateChange) {
      props.onGlobalStateChange((state, prev) => {
        console.log('[raven-mail] host state change', state)
        if (state.user) {
          userStore.setUser(state.user.id, state.user.name)
        }
        if (state.sessionId) {
          userStore.setSession(state.sessionId)
        }
        if (state.ravenConfig) {
          userStore.applyConfig(state.ravenConfig)
        }
      }, true) // fireImmediately: true
    }
    
    render(props)
  },
  unmount() {
    console.log('[raven-mail] unmount')
    if (app) {
      app.unmount()
      app = null
    }
  },
  update(props) {
    console.log('[raven-mail] update')
  }
})

// Standalone initialization
if (!qiankunWindow.__POWERED_BY_QIANKUN__) {
  render()
}
