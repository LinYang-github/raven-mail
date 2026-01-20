import { createRouter, createWebHistory } from 'vue-router'
import { qiankunWindow } from 'vite-plugin-qiankun/dist/helper'

const routes = [
  {
    path: '/',
    redirect: '/inbox'
  },
  {
    path: '/inbox',
    name: 'inbox-list',
    component: () => import('../components/EmptyView.vue')
  },
  {
    path: '/inbox/:id',
    name: 'inbox-detail',
    component: () => import('../components/MailDetail.vue')
  },
  {
    path: '/sent',
    name: 'sent-list',
    component: () => import('../components/EmptyView.vue')
  },
  {
    path: '/sent/:id',
    name: 'sent-detail',
    component: () => import('../components/MailDetail.vue')
  },
  {
    path: '/compose',
    name: 'compose',
    component: () => import('../components/ComposeView.vue')
  }
]

const createRouterInstance = () => createRouter({
  // Base will be dynamically set in main.js via window.__RAVEN_ROUTE_BASE__
  history: createWebHistory(qiankunWindow.__POWERED_BY_QIANKUN__ ? (window.__RAVEN_ROUTE_BASE__ || '/') : '/'),
  routes
})

export default createRouterInstance
