import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import './style.css'
import createRouterInstance from './router'

import { renderWithQiankun, qiankunWindow } from 'vite-plugin-qiankun/dist/helper'

import { userStore } from './store/user'

let app = null

function render(props = {}) {
  const { container, modules } = props
  app = createApp(App, { modules })
  
  app.use(ElementPlus)
  app.use(createRouterInstance()) // Create new router with current window.__RAVEN_ROUTE_BASE__
  
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
    console.log('[raven-app] bootstrap')
  },
  mount(props) {
    console.log('[raven-app] mount', props)
    
    // Sync initial state
    if (props.routeBase) {
      window.__RAVEN_ROUTE_BASE__ = props.routeBase
    }
    if (props.user) {
      userStore.setUser(props.user.id || '', props.user.name || '')
    }
    if (props.sessionId) {
      userStore.setSession(props.sessionId)
    }
    if (props.modules) {
      userStore.setModules(props.modules)
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
        console.log('[raven-app] host state change', state)
        if (state.user) {
          userStore.setUser(state.user.id, state.user.name)
        }
        if (state.sessionId) {
          userStore.setSession(state.sessionId)
        }
        if (state.ravenConfig) {
          userStore.applyConfig(state.ravenConfig)
        }
        if (state.modules) {
          userStore.setModules(state.modules)
        }
      }, true) // fireImmediately: true
    }
    
    render(props)
  },
  unmount() {
    console.log('[raven-app] unmount')
    if (app) {
      app.unmount()
      app = null
    }
  },
  update(props) {
    console.log('[raven-app] update')
  }
})

// Standalone initialization
if (!qiankunWindow.__POWERED_BY_QIANKUN__) {
  render({
    fetchUsers: async (query) => {
        return [
            { id: 'user-999', name: 'Test User', dept: 'Testing' },
            { id: 'user-007', name: 'Bond', dept: 'MI6' }
        ].filter(u => !query || u.name.includes(query) || u.id.includes(query))
    }
  })
}
