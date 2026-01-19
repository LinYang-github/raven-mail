import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import './style.css'

import { renderWithQiankun, qiankunWindow } from 'vite-plugin-qiankun/dist/helper'

let app = null

function render(props = {}) {
  const { container } = props
  app = createApp(App)
  
  app.use(ElementPlus)
  
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
    console.log('[raven-mail] mount')
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
